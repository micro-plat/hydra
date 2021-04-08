package nfs

import (
	"fmt"
	"hash/crc64"
	"net/http"
	"path/filepath"
	"time"

	"github.com/micro-plat/lib4go/errs"
)

//module 协调本地文件、本地指纹、远程指纹等处理
type module struct {
	path          string
	local         *local
	remoting      *remoting
	remotingHosts []string
	msg           *msg
	isMaster      bool
}

func newModule(path string, hosts []string, masterHost string, isMaster bool) (m *module, err error) {
	l, err := newLocal(path)
	if err != nil {
		return nil, err
	}
	m = &module{
		path:          path,
		remotingHosts: hosts,
		local:         l,
		isMaster:      isMaster,
		remoting:      newRemoting(path, hosts, masterHost, isMaster),
	}

	//查询远程服务器列表
	mp, err := m.remoting.Query()
	if err != nil {
		return nil, err
	}

	//合并到本地
	l.MergeFPS(mp)

	//处理本地新文件上报
	for _, f := range l.NFPS {
		m.msg.Report(f)
	}

	return m, nil
}

//Update 更新环境配置
func (m *module) Update(hosts []string, path string, masterHost string, isMaster bool) {
	m.remotingHosts = hosts
	m.path = path
	m.isMaster = isMaster
	m.remoting.Update(hosts, masterHost, isMaster)
}

//GetFile 获取文件,
//1. 从本地获取,有此文件直接返回结果
//2. 从master获取指纹信息，看哪些服务器有此文件
//3. 向有此文件的服务器拉取文件
//4. 保存到本地
//5. 通知master我也有这个文件了,如果是master则告诉所有人我也有此文件了
func (m *module) GetFile(name string) ([]byte, error) {

	//从本地获取
	buff, err := m.local.GetFile(name)
	if err == nil {
		return buff, err
	}

	//从远程获取
	fp, err := m.remoting.GetFPFormMaster(name)
	if err != nil {
		return nil, err
	}

	//从远程拉取文件
	buff, err = m.remoting.Pull(name, fp.GetAliveHost(m.remotingHosts...))
	if err != nil {
		return nil, err
	}

	//保存到本地
	fp, err = m.local.SaveFile(name, buff, fp.Hosts...)
	if err != nil {
		return nil, err
	}

	//上报给其它服务器
	m.msg.Report(fp)
	return buff, nil
}

//SaveNewFile 保存新文件到本地
//1. 查询本地是否有此文件了,有则报错
//2. 保存到本地，返回指纹信息
//3. 通知master我也有这个文件了,如果是master则告诉所有人我也有此文件了
func (m *module) SaveNewFile(name string, buff []byte) (*eFileFP, error) {
	//检查文件是否存在
	name = getFileName(name)
	if m.local.Has(name) {
		return nil, fmt.Errorf("文件名称重复:%s", name)
	}

	//保存到本地
	fp, err := m.local.SaveFile(name, buff)
	if err != nil {
		return nil, err
	}

	//远程通知
	m.msg.Report(fp)
	return fp, nil
}

//GetFP 获取本地的指纹信息，用于master对外提供服务
//1. 查询本地是否有文件的指纹信息
//2. 如果是master返回不存在
//3. 向master发起查询
func (m *module) GetLocalFP(name string) (*eFileFP, error) {
	//从本地文件获取
	if f, ok := m.local.FPS[name]; ok {
		return f, nil
	}
	return nil, errs.NewError(http.StatusNotFound, "文件不存在")
}

//GetFPList 获取本地包含有服务列表的指纹清单
func (m *module) GetFPList() eFileFPLists {
	return m.local.FPS
}

//RecvNotify 接收远程发过来的新文件通知
//1. 检查本地是否有些文件
//2. 文件不存在则自动下载
//3. 合并服务列表
func (m *module) RecvNotify(f *eFileFP) {
	if !m.local.Has(f.Path) {
		m.msg.Download(f)
	}
	m.local.MergeHost(f)
	return
}

//Close 关闭服务
func (m *module) Close() error {
	m.local.Close()
	m.msg.Close()
	return nil
}

func getCRC64(buff []byte) uint64 {
	return crc64.Checksum(buff, crc64.MakeTable(crc64.ISO))
}

func getFileName(name string) string {
	return filepath.Join(time.Now().Format("20060102"), name)
}
