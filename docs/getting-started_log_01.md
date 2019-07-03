# 日志配置与管理

hydra 使用[lib4go](https://github.com/micro-plat/lib4go)提供的日志库进行本地日志输出与存储。基于[rlog](https://github.com/micro-plat/rlog)实现了远程日志存储。

#### 1. 本地日志配置

一般日志配置文件不需要手动创建，初始运行时会检查../conf/logger.json 文件是否存在，不存在则会自动创建。可手工修改../conf/logger.json 改变日志存储目录和日志输出内容:

`../conf/logger.json`内容台下:

```json
[
  {
    "type": "file",
    "level": "All",
    "path": "/home/colin/work/logs/%date.log",
    "layout": "[%datetime.%ms][%l][%session] %content%n"
  },
  {
    "type": "stdout",
    "level": "All",
    "layout": "[%datetime.%ms][%l][%session]%content"
  }
]
```

参数说明:
|参数名|说明|
|:----:|----|
|type|日志输出类型，file:文件，stdout：终端|
|level|日志级别控制开关，值:All < Debug < Info < Warn < Error < Fatal < Off,日志只显示大于当前配置级别的日志，即配置为 Info 则，只有 Info,Warn,Error,Fatal 级别的文件会显示到终端或保存到日志文件|
|path|输出路径，日志输出类型为 file 时必须|
|layout|输出格式|

日志通过`%`开头的变量名进行参数转换:

| 参数名    | 说明                                                                                              |
| :-------- | ------------------------------------------------------------------------------------------------- |
| %session  | 当前日志的 sesssion id,每个日志组件创建时都会有一个 session id,用于从日志内容中区别同一组日志信息 |
| %date     | 日期，格式:20190703                                                                               |
| %datetime | 日期时间格式，格式:2019/07/03 11:18:07                                                            |
| %yy       | 年,格式:2019                                                                                      |
| %mm       | 月，格式：07                                                                                      |
| %dd       | 日期，格式：03                                                                                    |
| %hh       | 小时，24 小时制                                                                                   |
| %mi       | 分钟，格式：09                                                                                    |
| %ss       | 秒数,格式:04                                                                                      |
| %ms       | 毫秒，纳秒                                                                                        |
| %level    | 日志等级全称                                                                                      |
| %l        | 日志等级首字母                                                                                    |
| %name     | 日志名称                                                                                          |
| %pid      | 当前进程编号                                                                                      |
| %caller   | 日志调用函数名称                                                                                  |
| %content  | 日志内容                                                                                          |
| %index    | 日志序号                                                                                          |
| %ip       | 当前主机 ip                                                                                       |
| %n        | 换行符                                                                                            |
