# -*- coding: utf-8 -*-  

import os 
import time
import test_all
from test_help import test,execute,runApp
from test_help import LMAddress

testName = os.path.basename(__file__)

@test(testName)
def test_run_Normal():

    args = ["run","-r",LMAddress,"-c","c"]
    response = runApp(args)

    if not u"启动成功" in  response:
        return u"服务正常启动"
    

@test(testName)
def test_run_Withtrace():

    args = ["run","-r",LMAddress,"-c","c","-trace", "web"]
    response = runApp(args)

    if not u"启动成功:pprof.web" in  response:
        return u"trace服务正常启动"
    
    if not u"启动成功(api" in  response:
        return u"服务正常启动"
    
    time.sleep(1)
    #清理trace.out的监控文件
    os.remove("trace.out")


@test(testName)
def test_run_error_registry():

    args = ["run","-r","xx://xxx","-c","c"]
    response = runApp(args)

    if not u"不支持的协议类型" in  response:
        return u"不支持的协议类型测试"
    





def main():
     execute(testName)
 
if __name__ == "__main__":
    main()