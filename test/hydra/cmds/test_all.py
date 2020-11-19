# -*- coding: utf-8 -*-  

import sys 
import getopt

import platform

import select
import os,fnmatch
import time
import importlib
from test_help import failcases


def get_build_cmd():
    cmd=""
    if platform.system()== 'Windows':
        cmd= "where go"
        r=os.popen(cmd)
        text = r.read()
        r.close()
        cmd = text + " build"
    else:
        cmd= "/usr/local/go/bin/go build"

    print("get_cmd:",cmd)
    return cmd


    
def build_test_app():
    time.sleep(0.1)
    print(u"编译备用的服务程序：testapporg")
    os.chdir("./testapporg")
    time.sleep(0.1)
    os.system(get_build_cmd())
    os.chdir("../")

    time.sleep(0.1)
    print(u"编译备用的服务程序：testappnew")
    os.chdir("./testappnew")
    time.sleep(0.1)
    os.system(get_build_cmd())
    time.sleep(0.1)
    os.chdir("../")
    print(u"编译备用的服务程序：编译完成")

if __name__ == "__main__":    
    
    build_test_app()
    time.sleep(1)

    files=fnmatch.filter(os.listdir('.'), '*test.py')
    print(files) 
    for f in files:
        moduleName= f.strip(".py")
        print(u"执行:%s "% moduleName)
        m = importlib.import_module(moduleName)
        m.main()
        print(u"完成---------------")
        print(u"\n"*2)

    print(u"\n"*2)
    for k in failcases:
        print(u"测试失败:%s"% k)



