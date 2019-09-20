
### docker中运行程序
* 编写`dockerfile`文件
* 生成镜像`sudo docker build -t [镜像名称] .`
* 启动容器`sudo docker run --name [容器名称] -p 8090:8090 -d [镜像名称]`


### 镜像
* 查看镜像`sudo docker images`
* 删除镜像`sudo docker rmi [镜像名称]`

### 容器
* 查看容器 `sudo docker ps -a`
* 删除容器 `sudo docker rm -f [容器名称]`


### 日志
* 查看日志 `sudo docker logs hello`