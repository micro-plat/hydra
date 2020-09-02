#!/bin/bash

export path=$1
export name=$2
export params=$3

if [ "$path" == '' ]; then
    echo "服务路径不能为空"
    return
fi

if [ "$name" == '' ]; then
    echo "应用名称不能为空"
    return
fi




backupName=$name"`date +%Y%m%d%H%M%S`"



#文件是否存在,备份文件
if [ ! -d "$name"]; then
    mv $name $backupName
fi

exec "./$name stop"
exec "./$name remove"

if [ "$params" != '' ]; then
    exec "./$name install "$params
fi

exec "./$name start" 