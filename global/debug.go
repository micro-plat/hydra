package global

//IsDebug 是否是debug模式
//启动后服务器会打印详细的启动、运行日志。也会将详细的错误信息输出到client
//禁用后，当发生异常时client只会看到`Internal Server Error`，不会看到详细的错误消息
var IsDebug = false
