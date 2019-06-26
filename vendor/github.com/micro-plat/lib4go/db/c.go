package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro-plat/hydra/registry"
)

var errCodes = []string{
	"1051", "1050", "1061", "1091", "ORA-02289", "ORA-00942",
}

//CreateDB 创建数据库结构
//dir 为相对路径时应为基于$GOPATH的相对路径，否则应使用绝对路径
func CreateDB(db IDB, dir string) error {
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
	files, err := filepath.Glob(registry.Join(dir, "*.sql"))
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
