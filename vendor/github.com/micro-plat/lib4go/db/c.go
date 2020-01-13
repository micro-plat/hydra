package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var errCodes = []string{
	"1051", "1050", "1061", "1091", "ORA-02289", "ORA-00942",
}

//ParseConnectString 解析数据库连接串
//输入串可以是:"[用户名]:[密码]@[tns名称]" hydra:123456@hydra
//也可以是:"[用户名]:[密码]@[数据库名]/数据库ip"  hydra:123456@hydra/123456
func ParseConnectString(tp string, conn string) (string, error) {
	var uName, pwd, db, ip string
	ips := strings.SplitN(conn, "/", 2)
	if len(ips) > 1 {
		ip = ips[1]
	}
	dbs := strings.SplitN(ips[0], "@", 2)
	if len(dbs) > 1 {
		db = dbs[1]
	}
	up := strings.SplitN(dbs[0], ":", 2)
	if len(up) > 1 {
		pwd = up[1]
	}
	uName = up[0]
	switch tp {
	case "oracle", "ora":
		if uName == "" || pwd == "" || db == "" {
			return "", fmt.Errorf("数据为连接串错误:%s(格式:%s)", conn, `"[用户名]:[密码]@[tns名称]" hydra:123456@hydra`)
		}
		return fmt.Sprintf("%s/%s@%s", uName, pwd, db), nil
	case "mysql":
		if uName == "" || pwd == "" || db == "" || ip == "" {
			return "", fmt.Errorf("数据为连接串错误:%s(格式:%s)", conn, `"[用户名]:[密码]@[数据库名]/数据库ip"  hydra:123456@hydra/123456`)
		}
		return fmt.Sprintf("%s:%s@tcp(%s)/%s", uName, pwd, ip, db), nil
	default:
		return "", fmt.Errorf("不支持的数据库类型:%s", tp)
	}

}

//CreateDB 创建数据库结构
//dir 为相对路径时应为基于$GOPATH的相对路径，否则应使用绝对路径
func CreateDB(db IDBExecuter, dir string) error {
	path, err := getSQLPath(dir)
	if err != nil {
		return err
	}
	sqls, err := getSQL(path)
	if err != nil {
		return err
	}
	for _, sql := range sqls {
		if sql != "" {
			if _, q, _, err := db.Execute(sql, map[string]interface{}{}); err != nil {
				c := err.Error()
				exists := false
				for _, code := range errCodes {
					if strings.Contains(c, code) {
						exists = true
						break
					}
				}
				if !exists {
					return fmt.Errorf("执行SQL失败： %v %s", err, q)
				}
			}
		}
	}
	return nil
}

//getSQLPath 获取getSQLPath
func getSQLPath(dir string) (string, error) {
	if filepath.IsAbs(dir) {
		return dir, nil
	}
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return "", fmt.Errorf("未配置环境变量GOPATH")
	}
	path := strings.Split(gopath, ";")
	if len(path) == 0 {
		return "", fmt.Errorf("环境变量GOPATH配置的路径为空")
	}
	return filepath.Join(path[0], dir), nil
}
func getSQL(dir string) ([]string, error) {
	files, err := filepath.Glob(join(dir, "*.sql"))
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBufferString("")
	for _, f := range files {
		buf, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, err
		}
		_, err = buff.Write(buf)
		if err != nil {
			return nil, err
		}
		buff.WriteString(";")
	}
	tables := make([]string, 0, 8)
	tbs := strings.Split(buff.String(), ";")
	for _, t := range tbs {
		if tb := strings.TrimSpace(t); len(tb) > 0 {
			tables = append(tables, tb)
		}
	}
	return tables, nil
}
func join(elem ...string) string {
	path := filepath.Join(elem...)
	return strings.Replace(path, "\\", "/", -1)

}
