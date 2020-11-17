package update

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/lib4go/security/crc32"

	"github.com/asaskevich/govalidator"
	"github.com/mholt/archiver"
	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/global"
	"github.com/micro-plat/lib4go/errs"
	"github.com/micro-plat/lib4go/logger"
	"github.com/micro-plat/lib4go/sysinfo/pipes"
)

//GetPackage 获取文件package信息
func GetPackage(url string) (*Package, error) {
	resp, err := http.Get(url)
	if err != nil || resp == nil {
		err = fmt.Errorf("无法下载package更新包:%s %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("无法读取package更新包,状态码:%d", resp.StatusCode)
		return nil, err
	}
	if resp.ContentLength == 0 {
		err = fmt.Errorf("无法读取package更新包长度:%d", resp.ContentLength)
		return nil, err
	}

	var buffer []byte
	if buffer, err = ioutil.ReadAll(resp.Body); err != nil {
		return nil, fmt.Errorf("下载package更新包出错:%w", err)
	}

	cnf, err := conf.NewByText(buffer, 0)
	if err != nil {
		return nil, fmt.Errorf("package更新包必须是有效的json格式:%w", err)
	}
	var pkg Package
	if err = cnf.ToStruct(&pkg); err != nil {
		return nil, errs.NewError(406, err)
	}
	if b, err := govalidator.ValidateStruct(&pkg); !b {
		return nil, fmt.Errorf("package配置有误:%v", err)
	}
	return &pkg, nil
}

//Package 更新包
type Package struct {
	URL     string `json:"url,omitempty" valid:"url,required"`
	Version string `json:"version" valid:"ascii,required"`
	CRC32   uint32 `json:"crc32" valid:"required"`
}

//NewPackage 构建CRON任务
func NewPackage(url string, version string, crc32 uint32) *Package {
	return &Package{
		URL:     url,
		Version: version,
		CRC32:   crc32,
	}
}

//Check 是否需要更新
func (p *Package) Check() (bool, error) {
	if p.Version >= global.Version {
		return true, nil
	}
	return false, fmt.Errorf("更新的版本号%s不能低于当前应用的版本号%s", p.Version, global.Version)
}

//Update 更新当前服务
func (p *Package) Update(logger logger.ILogging, closeFunc func()) (err error) {
	logger.Info("开始下载更新包:", p.URL)
	resp, err := http.Get(p.URL)
	if err != nil || resp == nil {
		err = fmt.Errorf("无法下载更新包:%s", p.URL)
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
	updater, err := newUpdater(p.CRC32, filepath.Base(p.URL))
	if err != nil {
		err = fmt.Errorf("无法创建updater:%v", err)
		return
	}
	err = updater.Apply(resp.Body)
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

//Decode 将package信息保存到文件
func (p *Package) Decode(path string) error {
	buff, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("无法创建package json:%w", err)
	}

	f, err := createTruncFile(path)
	if err != nil {
		return fmt.Errorf("无法创建文件:%s %v", path, err)
	}
	defer f.Close()
	_, err = f.Write(buff)
	if err != nil {
		return fmt.Errorf("写入到文件失败:%s %v", path, err)
	}
	return nil
}

//Archive 生成压缩文件
func Archive(source string, destination string, url string) (err error) {

	//创建压缩包
	rpath := filepath.Join(destination, filepath.Base(source)+".zip")
	if err := archiver.Archive([]string{source}, rpath); err != nil {
		return fmt.Errorf("无法构建压缩包:[%s]%s %w", source, rpath, err)
	}
	defer func() {
		if err != nil {
			os.RemoveAll(rpath)
		}
	}()

	//读取文件
	buff, err := ioutil.ReadFile(rpath)
	if err != nil {
		return fmt.Errorf("无法读取文件%s %w", rpath, err)
	}

	//获取版本号
	v, err := getVerion(source)
	if err != nil {
		return err
	}

	//生成pkg文件
	pkg := NewPackage(url, v, crc32.Encrypt(buff))
	if err := pkg.Decode(filepath.Join(destination, "package.json")); err != nil {
		return err
	}
	return nil
}

//CreateTruncFile 根据文件路径(相对或绝对路径)创建文件，如果文件所在的文件夹不存在则自动创建
func createTruncFile(path string) (f *os.File, err error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return
	}
	dir := filepath.Dir(absPath)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return
	}
	return os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
}
func getVerion(p string) (string, error) {

	//获取版本号
	v, err := pipes.RunString(fmt.Sprintf("%s --version", p))
	if err != nil {
		return "", fmt.Errorf("无法获取应用程序版本号%w", err)
	}
	vers := strings.Split(v, " ")
	return vers[len(vers)-1], nil

}
