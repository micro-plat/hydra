#源镜像
FROM ubuntu:latest
#作者
MAINTAINER yanglei "lib4go@163.com"
#设置工作目录
WORKDIR /usr/local/bin
#将服务器的go工程代码加入到docker容器中
ADD ./dockerserver /usr/local/bin
#go构建可执行文件
# RUN go build .
#暴露端口
EXPOSE 8090
#最终运行docker的命令
ENTRYPOINT  ["./dockerserver", "start"]