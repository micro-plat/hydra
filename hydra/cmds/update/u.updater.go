package update

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mholt/archiver"
	"github.com/micro-plat/lib4go/osext"
	"github.com/micro-plat/lib4go/security/crc32"
)

type updater struct {
	targetPath   string
	currentDir   string
	newDir       string
	oldDir       string
	CRC32        uint32
	targetName   string
	tmpPath      string
	needRollback bool
}

// //UpdaterOptions 文件更新选项
// type updaterOptions struct {
// 	CRC32      uint32
// 	TargetName string
// }

// func newUpdaterOptions(crc32 uint32, name string) (*updaterOptions, error) {
// 	u := &updaterOptions{
// 		CRC32: crc32,
// 	}
// 	tmpDir, err := ioutil.TempDir("", "hydra-updater")
// 	if err != nil {
// 		return nil, fmt.Errorf("创建临时文件失败:%v", err)
// 	}

// }

func newUpdater(crc32 uint32, targetName string) (u *updater, err error) {
	u = &updater{CRC32: crc32, targetName: targetName}
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
func (u *updater) Apply(update io.Reader) (err error) {
	//读取文件内容
	var buff []byte
	if buff, err = ioutil.ReadAll(update); err != nil {
		return err
	}
	if u.CRC32 > 0 {
		if v := crc32.Encrypt(buff); v != u.CRC32 {
			err = fmt.Errorf("文件校验值有误当前[%d]%d", v, u.CRC32)
			return
		}
	}
	if err := u.write2Tmp(buff); err != nil {
		return err
	}

	//创建目标目录
	defer os.RemoveAll(u.tmpPath)
	defer os.Chdir(u.currentDir)
	os.RemoveAll(u.newDir)
	err = os.MkdirAll(u.newDir, 0755)
	if err != nil {
		err = fmt.Errorf("权限不足，无法创建文件:%s(err:%v)", u.newDir, err)
		return err
	}
	defer os.RemoveAll(u.newDir)

	//读取归档并解压文件

	err = archiver.Unarchive(u.tmpPath, u.newDir)
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

	//将新的目标文件修改为当前目录
	u.needRollback = true
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
func (u *updater) write2Tmp(buff []byte) error {
	tmpDir, err := ioutil.TempDir("", "hydra-updater")
	if err != nil {
		return fmt.Errorf("创建临时文件失败:%v", err)
	}
	u.tmpPath = filepath.Join(tmpDir, u.targetName)
	f, err := os.Create(u.tmpPath)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("无法创建文件:%s %v", u.tmpPath, err)
	}

	_, err = f.Write(buff)
	if err != nil {
		return fmt.Errorf("写入到临时文件失败:%s %v", u.tmpPath, err)
	}
	return nil
}
