#!/bin/bash

CURDIR=.

SYSOS='linux'

TEST_APP="testapp"



EXE_FILE=$TEST_APP

osname=`uname  -a`

if [[ $osname =~ 'Darwin' ]];then
    SYSOS="mac"
elif [[ $osname =~ 'centos' ]];then
    SYSOS="centos"
elif [[ $osname =~ 'ubuntu' ]];then
    SYSOS="ubuntu"
elif [[ $osname =~ 'WIN' ]];then
    SYSOS="windows"
    EXE_FILE="$TEST_APP.exe"
else
    echo "不支持的操作系统:$osname"
    exit 8
fi

echo "当前操作系统：$SYSOS"

exist_file()
{
    if [ -e "$1" ]
    then
        return 1
    else
        return 2
    fi
}

cd testapp 
echo "编译备用的服务程序：$EXE_FILE"

go build 

echo "编译备用的服务程序：编译完成"
sleep 1

cd ..

for k in $(ls ${CURDIR})
do
    echo "检查:$k" 
    if [ -d $k ] 
    then
        cd $k 
        exist_file test.sh
        value=$?
        if [ $value -eq 1 ]
        then
            echo "-------------测试执行开始：$k"
            sh test.sh  $EXE_FILE
            echo "-------------测试执行完毕：$k"
        fi 
        sleep 5
        cd ..
    fi
    
done 