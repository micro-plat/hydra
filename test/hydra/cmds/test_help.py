# -*- coding: utf-8 -*-  
from subprocess import Popen, PIPE
import platform
import time
import sys

defaultencoding = 'utf-8'
if sys.getdefaultencoding() != defaultencoding:
    reload(sys)
    sys.setdefaultencoding(defaultencoding)


ZKAddress="zk://192.168.0.101"
LMAddress="lm://."
 
testcases={}
failcases=[]

def test(name):
    def testInner(callback):
        if testcases.get(name) ==None:
            testcases[name]={}
        caseName = callback.__name__
     
        def inner():
           # print(testcases[name])
            result = callback()
            if result != None and result != "":
                print(u"Test:%s FAIL ,%s" %( caseName, result))
                failcases.append("%s.%s"% (name,caseName))
            else:
                print(u"Test:%s OK" % caseName)

        testcases[name][caseName ]=inner
        return  inner
    return testInner
   

def getcases(name):
    if testcases.get(name)==None:
        return []
    return testcases[name]

def execute(testName):
    cases = getcases(testName)
    for k in cases:
        cases[k].__call__()
    pass


class Server(object):
    def __init__(self, args, server_env = None):
        if server_env:
            self.process = Popen(args, stdin=PIPE, stdout=PIPE, stderr=PIPE, env=server_env)
        else:
            self.process = Popen(args, stdin=PIPE, stdout=PIPE, stderr=PIPE)

    def send(self, data, tail = '\n'):
       # print("send:",data)
        self.process.stdin.write(data + tail)
        self.process.stdin.flush()

    def recv(self, t=.1, stderr=0):
        r = ''
        pr = self.process.stdout
        if stderr:
            pr = self.process.stdout
        r = pr.read()
        return r.rstrip()

    def kill(self):
        self.process.kill()


exenames={
    "org":"./testapporg/testapporg",
    "new":"./testappnew/testappnew",
}


def runApp(args,steps=[],exe_name=""):
    if exe_name == "":
        exe_name='./testapporg/testapporg'

    exe_name = exenames.get(exe_name) or exe_name
    if platform.system()== 'Windows':
        exe_name=exe_name+'.exe'

    srvArgs = [exe_name]
    for i in args:
        srvArgs.append(i)

    server = Server(srvArgs)
    time.sleep(2)

    try:
        for val in steps:
            server.send(val)
            time.sleep(1)

        if len(steps)>0:
            time.sleep(2)
            
    except Exception as ex:
        print(ex)
        return "error:%s"% ex.message       

    finally:   
        server.kill()
        print("server.kill")

    response = server.recv()
    #print(unicode(response,"utf-8"))
    return unicode(response)

