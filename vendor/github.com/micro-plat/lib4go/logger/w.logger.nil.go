package logger

var _ ILogger = &nilLog{}

type nilLog struct {
	session string
	name    string
}

//Nil 构建空的日志组件
func Nil() (logger *nilLog) {
	return &nilLog{
		session: CreateSession(),
	}
}

//SetName 设置日志名称
func (logger *nilLog) SetName(name string) {
	logger.name = name
}

//Close 关闭当前日志组件
func (logger *nilLog) Close() {
}

//Pause 暂停记录
func (logger *nilLog) Pause() {
}

//Resume 恢复记录
func (logger *nilLog) Resume() {
}

//GetSessionID 获取当前日志的session id
func (logger *nilLog) GetSessionID() string {
	return logger.session
}

//Debug 输出debug日志
func (logger *nilLog) Debug(content ...interface{}) {
}

//Debugf 输出debug日志
func (logger *nilLog) Debugf(format string, content ...interface{}) {
}

//Info 输出info日志
func (logger *nilLog) Info(content ...interface{}) {
}

//Infof 输出info日志
func (logger *nilLog) Infof(format string, content ...interface{}) {
}

//Warn 输出info日志
func (logger *nilLog) Warn(content ...interface{}) {
}

//Warnf 输出info日志
func (logger *nilLog) Warnf(format string, content ...interface{}) {
}

//Error 输出Error日志
func (logger *nilLog) Error(content ...interface{}) {
}

//Errorf 输出Errorf日志
func (logger *nilLog) Errorf(format string, content ...interface{}) {
}

//Fatal 输出Fatal日志
func (logger *nilLog) Fatal(content ...interface{}) {
}

//Fatalf 输出Fatalf日志
func (logger *nilLog) Fatalf(format string, content ...interface{}) {
}

//Fatalln 输出Fatal日志
func (logger *nilLog) Fatalln(content ...interface{}) {
}

//Print 输出info日志
func (logger *nilLog) Print(content ...interface{}) {
}

//Printf 输出info日志
func (logger *nilLog) Printf(format string, content ...interface{}) {
}

//Println 输出info日志
func (logger *nilLog) Println(content ...interface{}) {

}
