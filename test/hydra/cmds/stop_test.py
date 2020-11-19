# -*- coding: utf-8 -*-  

import os ,platform
import test_all
from test_help import test,execute,runApp
from test_help import ZKAddress

testName = os.path.basename(__file__)

@test(testName)
def test_stop_Normal_running():
    try:
        #1.安装服务
        args = ["install","-r",ZKAddress,"-c","c"]
        response = runApp(args)
        #print("install",response)
        if not "OK" in response:
            return u"安装服务失败"

        #2.启动
        args = ["start"]
        response = runApp(args)
        #print("start",response)
        if not "OK" in response:
            return u"启动服务失败"

        #3.停止
        args = ["stop"]
        response = runApp(args)
        #print("stop",response)
        if not "OK" in response:
            return u"停止服务失败"

    except Exception as err :
        return "error:"+err.message
        
    finally:    
        #4.删除
        args = ["remove"]
        response = runApp(args)
        #print("remove",response)

        if not "OK" in response:
            return u"删除服务失败"




@test(testName)
def test_stop_Not_installed():
    #1. 清理服务，避免其他遗留存在服务
    args = ["remove"]
    runApp(args)

    #2.启动服务
    args = ["stop"]
    response = runApp(args)
    
    if not ("not" in response and "installed" in response):
        return u"未安装服务验证"



@test(testName)
def test_stop_has_stopped():
    try:
        args = ["remove"]
        runApp(args)
        
        #1.安装服务
        args = ["install","-r",ZKAddress,"-c","c"]
        runApp(args)
       
        #3.停止
        args = ["stop"]
        response = runApp(args)  
        if "service has been stopped" not in response:
            return u"停止关闭的服务"

    except Exception as err :
        return "error:"+err.message
        
    finally:    
        #4.删除
        print("finally:remove")
        args = ["remove"]
        runApp(args)


def main():
     execute(testName)
 
if __name__ == "__main__":
    main()