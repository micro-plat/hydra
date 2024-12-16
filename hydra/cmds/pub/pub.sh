#!/bin/bash

export TMP_FILE=$1
export BIN_NAME=$2
export PROJECT_PATH=$3
export INSTALL_PARAMS=$4

#1. tmp文件名
#2. 真实文件名（程序名称）
#3. 安装文件路径
#4. 安装命令

source /etc/profile

curstep=0
newname=""
time=$(date "+%Y%m%d%H%M")
newfile="${BIN_NAME}_daemons_${time}"
 
function rollback(){
	
	case $curstep in 
	1)
		echo "创建目标目录失败"
		return 0
	;;
	2)
		echo "停止原有服务失败"
		return 0
	;;
	3)
		echo "复制新文件失败，执行回滚"
		if [ -n ${newname} ] ; then
			echo "rollback 1:恢复文件mv ./${newname} ${BIN_NAME}"
			mv  -f ./${newname} ${BIN_NAME}
			echo "rollback 2:启动原文件./${BIN_NAME} start"
			start_tips=`./${BIN_NAME} start`
			echo "rollback 3:启动结果：${start_tips}"
		else
			echo "rollback 1:原程序未启动，无需处理"
		fi
	;;
	4)
		echo "backup失败，执行回滚"
		if [ -n ${newname} ] ; then
			echo "rollback 1:恢复文件mv ./${newname} ${BIN_NAME}"
			mv  -f ./${newname} ${BIN_NAME}
			echo "rollback 2:启动原文件./${BIN_NAME} start"
			start_tips=`./${BIN_NAME} start`
			echo "rollback 3:启动结果：${start_tips}"
		else
			echo "rollback 1:原程序未启动，无需处理"
		fi 
		
	;;
	5)
		echo "remove失败，执行回滚"
		if [ -n ${newname} ] ; then
			echo "rollback 1:恢复文件mv ./${newname} ${BIN_NAME}"
			mv -f ./${newname} ${BIN_NAME}
			echo "rollback 2:启动原文件./${BIN_NAME} start"
			start_tips=`./${BIN_NAME} start`
			echo "rollback 3:启动结果：${start_tips}"
		else
			echo "rollback 1:原程序未启动，无需处理"
		fi 
		
	
	;;
	6)
		echo "install失败，执行回滚"
		if [ -n ${newname} ] ; then
			echo "rollback 1:恢复文件mv ./${newname} ${BIN_NAME}"
			mv -f ./${newname} ${BIN_NAME}
			echo "rollback 2:启动原文件./${BIN_NAME} start"
			start_tips=`./${BIN_NAME} start`
			echo "rollback 3:启动结果：${start_tips}"
		else
			echo "rollback 1:原程序未启动/未安装，无需处理"
		fi  
	;;
	7)
		echo "start失败，执行回滚"
		if [ -n ${newname} ] ; then
			echo "rollback 1:恢复文件mv ./${newname} ${BIN_NAME}"
			mv -f ./${newname} ${BIN_NAME}
			
			
			if [ -n ${newfile} ] ; then 
				echo "rollback 2:启动原文件:./${BIN_NAME} rollabck -f ${newfile}"
				start_tips=`./${BIN_NAME} rollabck -f ${newfile}`
			fi
			start_tips=`./${BIN_NAME} start`
			echo "rollback 3:启动结果：${start_tips}" 
		else
			echo "rollback 1:原程序未启动/未安装，无需处理"
		fi  		
	
	;;
	8)
		echo "none case 8"
	
	;;
	esac
	
	echo rollback
 
}

function checksucc(){
	msg=$1
	chk1=$2
	chk2=$3
	echo "checksucc,msg:${msg},chk1:${chk1},chk2:${chk2}"
	if echo $msg | grep -e $chk1
	then
		return 0
	else
		if echo $msg | grep -e $chk2
		then
			return 0
		else
			rollback
			return 1
		fi
	fi 	
}

 
output_name=""

BIN_DIR=~/${PROJECT_PATH}/bin

echo "1.检查目标目录${BIN_DIR}"
if [ ! -d ${BIN_DIR} ];then 
	echo "1.1.目标目录不存在，执行创建:${BIN_DIR}"
    mkdir -p ${BIN_DIR}
	if [ $? -gt 0 ]
	then  
		curstep=1
		rollabck
		exit
	fi  
	sleep 1
fi

echo "2.0.打开目标目录cd ${BIN_DIR}"
cd ${BIN_DIR}
sleep 1
pwd


echo "2.检查是否有${BIN_NAME}文件"
if [ -f ${BIN_NAME} ];then
	echo "2.1.文件存在，执行stop指令"
	stop_tips=`./${BIN_NAME} stop`
	#Service is not installed
	#Stopping hbs-api:					[  OK  ]
	echo ${stop_tips}
	curstep=2
	checksucc "${stop_tips}" "OK" "installed"
	if [ $? -gt 0 ]
	then 
		exit
	fi  
	sleep 1
	
	echo "2.2.文件存在，执行stop指令"	
	newname=${BIN_NAME}_${time}
	mv ./${BIN_NAME} ${newname}

fi


echo "3.拷贝生成文件${BIN_NAME}到目标地址(/tmp/${TMP_FILE} ./${BIN_NAME})" 
cp /tmp/${TMP_FILE} ./${BIN_NAME}

if [ $? -gt 0 ]
then  
	curstep=3
	rollabck
	exit
fi 


sleep 1

echo "4.修改${BIN_NAME}的权限为755"
chmod 755 ${BIN_NAME}
sleep 1


echo "5.判定安装参数:-n ${INSTALL_PARAMS}"
if [ -n "${INSTALL_PARAMS}" ] ; then
	echo "5.0.存在安装参数"
	
	echo "5.1 执行backup指令:./${BIN_NAME} backup -f ${newfile}"
	backup_tips=`./${BIN_NAME} backup -f ${newfile}`
	curstep=4
	checksucc "${backup_tips}" "OK" "installed"
	if [ $? -gt 0 ]
	then
		exit
	fi

	echo "5.1 执行remove指令:./${BIN_NAME} remove"
	remove_tips=`./${BIN_NAME} remove`
	echo $remove_tips
	#Removing hbs-api:					[  OK  ]
	#Service is not installed
	curstep=5
	checksucc "${remove_tips}" "OK" "installed"
	if [ $? -gt 0 ]
	then 
		exit
	fi  
	sleep 1
		
	echo "5.2.执行install指令:./${BIN_NAME} install ${INSTALL_PARAMS}"
	install_tips=`./${BIN_NAME} install ${INSTALL_PARAMS}`
	echo $install_tips
	#Install hjjymgrweb:					[  OK  ]	
	curstep=6
	checksucc "${install_tips}" "OK"
	if [ $? -gt 0 ]
	then 
		exit
	fi  
 
	sleep 1
fi 
 
echo "启动服务：${BIN_NAME} start"
start_tips=`./${BIN_NAME} start`
echo ${start_tips}
curstep=7
checksucc "${start_tips}" "OK"
if [ $? -gt 0 ]
then 
	exit
fi  

if [ -f $newfile ] ; then
	echo "删除daemons备份文件：${newfile}"
	rm -f $newfile	
fi

#rm -f /tmp/${TMP_FILE}

sleep 1
 