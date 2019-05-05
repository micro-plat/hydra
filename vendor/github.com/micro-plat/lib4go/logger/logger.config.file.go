package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"strings"

	"github.com/micro-plat/lib4go/file"
)

var loggerPath, _ = file.GetAbs("../conf/logger.json")
var configAdapter map[string]func() []*Appender
var defaultConfigAdapter string

//register 设置配置文件
func register(adapterName string, f func() []*Appender) error {
	if configAdapter == nil {
		configAdapter = make(map[string]func() []*Appender)
	}
	if _, ok := configAdapter[adapterName]; ok {
		return fmt.Errorf("adapter(%s) is exist", adapterName)
	}
	configAdapter[adapterName] = f
	defaultConfigAdapter = adapterName
	return nil
}

func readFromFile() (appenders []*Appender) {
	var err error
	appenders, err = read()
	if err == nil {
		return
	}
	appenders = getDefConfig()
	//	sysLoggerError(err)
	err = writeToFile(loggerPath, appenders)
	if err != nil {
		sysLoggerError(err)
	}

	return
}

//NewAppender 构建appender
func NewAppender(conf string) (appenders []*Appender, err error) {
	appenders = make([]*Appender, 0, 2)
	if err = json.Unmarshal([]byte(conf), &appenders); err != nil {
		err = errors.New("配置文件格式有误，无法序列化")
		return
	}
	return
}

// // TimeClear 定时清理loggermanager时间间隔
// var TimeClear = time.Second

// TimeWriteToFile 定时写入文件时间间隔

func read() (appenders []*Appender, err error) {
	currentAppenders := make([]*Appender, 0, 2)
	if !exists(loggerPath) {
		err = errors.New("配置文件不存在:" + loggerPath)
		return
	}
	bytes, err := ioutil.ReadFile(loggerPath)
	if err != nil {
		err = errors.New("无法读取配置文件")
		return
	}
	if err = json.Unmarshal(bytes, &currentAppenders); err != nil {
		err = errors.New("配置文件格式有误，无法序列化")
		return
	}
	if len(currentAppenders) == 0 {
		return
	}
	appenders = make([]*Appender, 0, len(currentAppenders))
	for _, v := range currentAppenders {
		if strings.EqualFold(v.Level, "off") {
			continue
		}
		appenders = append(appenders, v)
	}
	return
}
func writeToFile(loggerPath string, appenders []*Appender) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	fwriter, err := file.CreateFile(loggerPath)
	if err != nil {
		return
	}
	data, err := json.Marshal(appenders)
	if err != nil {
		return
	}
	_, err = fwriter.Write(data)
	if err != nil {
		return
	}
	fwriter.Close()
	//sysLoggerError("已创建日志配置文件:", loggerPath)
	return
}
func getDefConfig() (appenders []*Appender) {
	fileAppender := &Appender{Type: "file", Level: SLevel_ALL}
	fileAppender.Path, _ = file.GetAbs("../logs/%date.log")
	fileAppender.Layout = "[%datetime.%ms][%l][%session] %content%n"
	appenders = append(appenders, fileAppender)

	sdtoutAppender := &Appender{Type: "stdout", Level: SLevel_ALL}
	sdtoutAppender.Layout = "[%datetime.%ms][%l][%session]%content"
	appenders = append(appenders, sdtoutAppender)

	return
}
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil || os.IsExist(err)
}
func getCaller(index int) string {
	defer recover()
	_, file, line, ok := runtime.Caller(index)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return ""
}
