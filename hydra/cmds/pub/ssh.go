package pub

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/micro-plat/lib4go/utility"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var client = &sshClient{}

type sshClient struct {
	ip         string
	userName   string
	pwd        string
	client     *ssh.Client
	tmpPath    string
	localpath  string
	remotePath string
	remoteFile string
}

func (s *sshClient) Bind(host string, localpath string, pwd string) error {
	if host == "" {
		return fmt.Errorf("未指定远程服务器信息")
	}
	if pwd == "" {
		return fmt.Errorf("未指远程服务器登录密码")
	}
	if !strings.Contains(host, ":") {
		return fmt.Errorf("%s 服务器信息应包含远程目录,格式:userName@ip:/path", host)
	}
	if !strings.Contains(host, "@") {
		return fmt.Errorf("%s 服务器信息应包含远程服务器ip地址,格式:userName@ip:/path", host)
	}
	paths := strings.Split(host, ":")
	if len(paths) != 2 || paths[1] == "" || paths[0] == "" {
		return fmt.Errorf("%s 远程路径有误,格式:userName@ip:/path", host)
	}

	hosts := strings.Split(paths[0], "@")
	if len(hosts) != 2 || hosts[1] == "" || hosts[0] == "" {
		return fmt.Errorf("%s 远程服务有误,格式:userName@ip:/path", host)
	}
	s.userName = hosts[0]
	s.ip = hosts[1]
	s.localpath = localpath
	s.remotePath = paths[1]
	s.remoteFile = path.Join(s.remotePath, path.Base(s.localpath))
	s.tmpPath = path.Join(os.TempDir(), utility.GetGUID())
	s.pwd = pwd
	return nil
}

//登录到服务器
func (s *sshClient) Login() (err error) {
	if s.ip == "" || s.userName == "" || s.pwd == "" {
		return fmt.Errorf("服务器ip,用户名，密码不能为空")
	}
	//通过ssh连接到远程服务器
	s.client, err = ssh.Dial("tcp", s.ip, &ssh.ClientConfig{
		User: s.userName,
		Auth: []ssh.AuthMethod{ssh.Password(s.pwd)},
	})
	return err
}

//run 执行命令
func (s *sshClient) run(cmd string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("服务器未登录")
	}
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buff, err := session.CombinedOutput(cmd)
	return string(buff), err
}

//获取当前文件名
func (s *sshClient) GetFileName() string {
	return path.Base(s.localpath)
}

//上传文件
func (s *sshClient) UploadFile() error {

	//1.处理文件名
	srcFile, err := os.Open(s.localpath)
	if err != nil {
		return fmt.Errorf("打开文件失败%w", err)

	}
	defer srcFile.Close()

	//2. 构建sftp客户端
	ftpclient, err := sftp.NewClient(s.client)
	if err != nil {
		return fmt.Errorf("创建sftp客户端失败%w", err)
	}
	defer ftpclient.Close()

	//3. 创建远程文件

	dstFile, e := ftpclient.Create(s.tmpPath)
	if e != nil {
		return fmt.Errorf("创建文件失败%w", e)

	}
	defer dstFile.Close()
	buffer := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return fmt.Errorf("读取文件出错%w", err)
			}
		}
		dstFile.Write(buffer[:n])
	}
	return nil
}

//上传脚本
func (s *sshClient) UploadScript() (string, error) {

	//1. 构建sftp客户端
	ftpclient, err := sftp.NewClient(s.client)
	if err != nil {
		return "", fmt.Errorf("创建ftp客户端失败%w", err)
	}
	defer ftpclient.Close()

	//2. 创建远程文件
	p, script := getScript()
	dstFile, e := ftpclient.Create(path.Join(s.remotePath, p))
	if e != nil {
		return "", fmt.Errorf("创建文件失败%w", e)

	}
	defer dstFile.Close()
	dstFile.Write([]byte(script))

	return p, nil
}

//GoWorkDir 转到工作目录
func (s *sshClient) GoWorkDir() (err error) {
	if _, err = s.run(cmdCD.CMD(s.remotePath)); err == nil {
		return nil
	}
	if _, err = s.run(cmdMkdir.CMD(s.remotePath)); err != nil {
		return err
	}
	if _, err = s.run(cmdCD.CMD(s.remotePath)); err == nil {
		return nil
	}
	return err
}

//ExecScript
func (s *sshClient) ExecScript(p string) error {
	if _, err := s.run(cmdRunScript.CMD(p)); err == nil {
		return nil
	}
	return nil
}
func (s *sshClient) Close() error {
	return nil
}
