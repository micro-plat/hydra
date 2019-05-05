package rqs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/context"
	"github.com/micro-plat/hydra/registry"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/osext"
	"github.com/micro-plat/lib4go/security/crc32"
	"github.com/zkfy/archiver"
)

type updater struct {
	targetPath   string
	currentDir   string
	newDir       string
	oldDir       string
	needRollback bool
}

func newUpdater() (u *updater, err error) {
	u = &updater{}
	u.targetPath, err = osext.Executable()
	if err != nil {
		return nil, err
	}
	u.currentDir = filepath.Dir(u.targetPath)
	u.newDir = u.currentDir + ".new"
	u.oldDir = u.currentDir + ".old"
	return
}

//Apply 更新文件
func (u *updater) Apply(update io.Reader, opts updaterOptions) (err error) {
	//读取文件内容
	var newBytes []byte
	if newBytes, err = ioutil.ReadAll(update); err != nil {
		return err
	}
	if opts.CRC32 > 0 {
		if crc32.Encrypt(newBytes) != opts.CRC32 {
			err = fmt.Errorf("文件校验值有误")
			return
		}
	}

	//创建目标目录
	defer os.Chdir(u.currentDir)
	os.RemoveAll(u.newDir)
	err = os.MkdirAll(u.newDir, 0755)
	if err != nil {
		err = fmt.Errorf("权限不足，无法创建文件:%s(err:%v)", u.newDir, err)
		return err
	}
	defer os.RemoveAll(u.newDir)

	//读取归档并解压文件
	archiver := archiver.MatchingFormat(opts.TargetName)
	if archiver == nil {
		err = fmt.Errorf("文件不是有效的归档或压缩文件")
		return
	}
	err = archiver.Read(bytes.NewReader(newBytes), u.newDir)
	if err != nil {
		err = fmt.Errorf("读取归档文件失败:%v", err)
		return
	}
	//备份当前目录
	err = os.RemoveAll(u.oldDir)
	if err != nil {
		err = fmt.Errorf("移除目录失败:%s(err:%v)", u.oldDir, err)
		return
	}
	err = os.Rename(u.currentDir, u.oldDir)
	if err != nil {
		err = fmt.Errorf("无法修改当前工作目录:%s(%s)(err:%v)", u.currentDir, u.oldDir, err)
		return err
	}
	u.needRollback = true
	//将新的目标文件修改为当前目录
	err = os.Rename(u.newDir, u.currentDir)
	if err != nil {
		err = fmt.Errorf("重命名文件夹失败:%v", err)
		return
	}
	err = os.Chdir(u.currentDir)
	if err != nil {
		err = fmt.Errorf("切换工作目录失败:%v", err)
		return
	}
	return
}

//Rollback 回滚当前更新
func (u *updater) Rollback() error {
	if !u.needRollback {
		return nil
	}
	defer os.Chdir(u.currentDir)
	if _, err := os.Stat(u.oldDir); os.IsNotExist(err) {
		return fmt.Errorf("无法回滚，原备份文件(%s)不存在", u.oldDir)
	}

	//当前工作目录不存在，直接将备份目录更改为当前目录
	if _, err := os.Stat(u.currentDir); os.IsNotExist(err) {
		return os.Rename(u.oldDir, u.currentDir)
	}
	os.RemoveAll(u.currentDir)
	err := os.Rename(u.oldDir, u.currentDir)
	if err != nil {
		return err
	}
	return nil
}

//UpdaterOptions 文件更新选项
type updaterOptions struct {
	CRC32      uint32
	TargetName string
}

//NeedUpdate 检查当前服务是否需要更新
func NeedUpdate(r registry.IRegistry, platName string, systemName string, v string) (bool, *conf.Package, error) {
	path := registry.Join("/", platName, "package", systemName, v)
	b, err := r.Exists(path)
	if err != nil {
		return false, nil, context.NewError(406, err)
	}
	if !b {
		return false, nil, context.NewError(406, fmt.Sprintf("版本%s不存在(%s未配置)", v, path))
	}
	buffer, vr, err := r.GetValue(path)
	if err != nil {
		return false, nil, context.NewError(406, err)
	}
	cnf, err := conf.NewJSONConf(buffer, vr)
	if err != nil {
		return false, nil, err
	}
	var pkg conf.Package
	if err = cnf.Unmarshal(&pkg); err != nil {
		return false, nil, context.NewError(406, err)
	}
	if b, err := govalidator.ValidateStruct(&pkg); !b {
		return false, nil, fmt.Errorf("package配置有误:%v", err)
	}
	if pkg.Version > v {
		return true, &pkg, nil
	}
	return false, nil, nil
}

//UpdateNow 更新当前服务
func UpdateNow(pkg *conf.Package, logger *logger.Logger, closeFunc func()) (err error) {
	logger.Info("开始下载更新包:", pkg.URL)
	resp, err := http.Get(pkg.URL)
	if err != nil || resp == nil {
		err = fmt.Errorf("无法下载更新包:%s", pkg.URL)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("无法读取更新包,状态码:%d", resp.StatusCode)
		return
	}
	if resp.ContentLength == 0 {
		err = fmt.Errorf("无法读取更新包长度:%d", resp.ContentLength)
		return
	}
	updater, err := newUpdater()
	if err != nil {
		err = fmt.Errorf("无法创建updater:%v", err)
		return
	}
	err = updater.Apply(resp.Body, updaterOptions{CRC32: pkg.CRC32, TargetName: filepath.Base(pkg.URL)})
	if err != nil {
		if err1 := updater.Rollback(); err1 != nil {
			err = fmt.Errorf("更新失败%+v,回滚失败%v", err, err1)
			return
		}
		err = fmt.Errorf("更新失败,已回滚(err:%v)", err)
		return
	}
	logger.Info("更新成功，停止所有服务，并准备重启")
	closeFunc()
	return restart()
}

//Restart 重启当前服务
func restart() (err error) {
	args := getExecuteParams(os.Args)
	//h.Info("准备重启：", args)
	go func() {
		time.Sleep(time.Second * 5)
		cmd1 := exec.Command("/bin/bash", "-c", args)
		cmd1.Stdout = os.Stdout
		cmd1.Stderr = os.Stderr
		cmd1.Stdin = os.Stdin
		err = cmd1.Start()
		if err != nil {
			return
		}

		os.Exit(20)
	}()
	return nil
}
func getExecuteParams(input []string) string {
	args := make([]string, 0, len(input))
	for i := 0; i < len(input); i++ {
		if strings.HasPrefix(input[i], "-") {
			args = append(args, input[i])
			if i+1 < len(input) && !strings.HasPrefix(input[i+1], "-") {
				args = append(args, fmt.Sprintf(`"%s"`, input[i+1]))
				i++
			}
		} else {
			args = append(args, input[i])
		}
	}
	return strings.Join(args, " ")
}
