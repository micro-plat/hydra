#!/bin/sh
  
dt=$(date "+%Y%m%d%H%M%S")

#------------------------------------" 


echo "1. 编译项目" 

go build --tags="prod"

echo "2. 复制文件"

scp ./docs root@121.196.168.242:/tmp


ssh -t  root@121.196.168.242 "cd /root/work/docs/bin;./docs stop;mv ./docs ./docs_${dt} ;cp /tmp/docs ./;sleep 3;./docs start;rm -rf /tmp/docs"


rm -rf ./docs

echo "3. 发布成功"