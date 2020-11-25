# -*- coding: utf-8 -*-  

import os 
import test_all
from test_help import test,execute,runApp
from test_help import ZKAddress

testName = os.path.basename(__file__)

@test(testName)
def test_Install_Normal():
    #执行install
    args = ["remove"]
    runApp(args)
    
    args = ["install","-r",ZKAddress,"-c","c"]
    response = runApp(args)
    if not "OK" in response:
        return u"安装服务失败"

    args = ["remove"]
    response = runApp(args)
    if not "OK" in response:
        return u"删除服务失败"


@test(testName)
def test_Install_Less_param():
    args = ["remove"]
    runApp(args,exe_name="new")
    #缺少集群名称
    args = ["install","-r",ZKAddress]
    response = runApp(args,exe_name="new")
    #print(response)
    if not u"平台名称不能为空" in response:
        return u"缺少参数安装服务用例失败"



@test(testName)
def test_Install_Cover():
    args = ["install","-r",ZKAddress,"-c","c","-cover","true"]
    response = runApp(args)
    if not u"OK" in response:
        return u"覆盖安装服务用例失败"
    
    #清理安装的服务
    args = ["remove"]
    runApp(args)


@test(testName)
def test_remove_NotExists():
    #1.清理已存在的服务
    args = ["remove"]
    runApp(args)

    #2.执行删除不存在的服务
    args = ["remove"]
    response = runApp(args)
    if not (u"not" in response and u"installed" in response):
        return u"删除不存在的服务用例"
    


def main():
     execute(testName)
 
if __name__ == "__main__":
    main()