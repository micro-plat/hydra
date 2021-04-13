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
	path     string
	local    *local
	remoting *remoting
	msg      *msg
}

func newModule(path string) (m *module) {
	m = &module{
		path:     path,
		local:    newLocal(path),
		remoting: newRemoting(),
	}
	m.msg = newMsg(m.local, m.remoting)
	return m
}

//Update 更新环境配置
func (m *module) Update(path string, hosts []string, masterHost string, currentAddr string, isMaster bool) {
	m.path = path
	m.remoting.Update(hosts, masterHost, currentAddr, isMaster)
	m.local.Update(path, currentAddr)
	m.msg.Update(hosts)
	if isMaster {
		m.Report()
	}
}

//从服务器拉取配置，并进行同步
func (m *module) Report() {
	//查询远程服务器列表
	mp, err := m.remoting.Query()
	if err != nil {
		return
	}

	mp.Merge(m.local.GetFPList())
	//结合外部传入，与当前服务器，进行整体合并，并进行通知
	m.msg.Report(GetAllNotify(mp, append(m.remoting.hosts, m.remoting.currentAddr)))
	return
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
	fmt.Println("get.from.master:", name, m.local.GetFPList())
	fp, err := m.remoting.GetFPFormMaster(name)
	if err != nil {
		return nil, err
	}

	//从远程拉取文件
	buff, err = m.remoting.Pull(fp)
	if err != nil {
		return nil, err
	}

	//保存到本地
	fp, err = m.local.SaveFile(name, buff, fp.Hosts...)
	if err != nil {
		return nil, err
	}

	//上报给其它服务器
	m.msg.Report(GetAllNotify(fp.GetMAP(), m.remoting.hosts))
	return buff, nil
}

//GetLocalFile 获取本地文件
func (m *module) GetLocalFile(name string) ([]byte, error) {
	buff, err := m.local.GetFile(name)
	return buff, err
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
	m.msg.Report(GetAllNotify(fp.GetMAP(), m.remoting.getHosts()))
	return fp, nil
}

//GetFP 获取本地的指纹信息，用于master对外提供服务
//1. 查询本地是否有文件的指纹信息
//2. 如果是master返回不存在
//3. 向master发起查询
func (m *module) GetLocalFP(name string) (*eFileFP, error) {
	//从本地文件获取
	if f, ok := m.local.GetFPByName(name); ok {
		return f, nil
	}
	return nil, errs.NewError(http.StatusNotFound, "文件不存在")
}

//GetFPList 获取本地包含有服务列表的指纹清单
func (m *module) GetFPList() eFileFPLists {
	return m.local.GetFPList()
}

//RecvNotify 接收远程发过来的新文件通知
//1. 检查本地是否有些文件
//2. 文件不存在则自动下载
//3. 合并服务列表
func (m *module) RecvNotify(f eFileFPLists) error {
	//处理本地新文件上报
	reports, downloads, err := m.local.MergeLocal(f)
	if err != nil {
		return err
	}
	m.msg.Report(GetAllNotify(reports, m.remoting.getHosts()))
	for _, f := range downloads {
		m.msg.Download(f)
	}
	return m.local.FPWrite(m.local.FPS.Items())
}

//Close 关闭服务
func (m *module) Close() error {
	if m.local != nil {
		m.local.Close()
	}
	if m.msg != nil {
		m.msg.Close()
	}
	return nil
}

func getCRC64(buff []byte) uint64 {
	return crc64.Checksum(buff, crc64.MakeTable(crc64.ISO))
}

func getFileName(name string) string {
	return filepath.Join(time.Now().Format("20060102"), name)
}
