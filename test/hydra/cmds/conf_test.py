# -*- coding: utf-8 -*-  
import os 
import test_all
from test_help import test,execute,runApp
from test_help import ZKAddress

testName = os.path.basename(__file__)

@test(testName)
def test_Conf_Show_Normal():
    """
    测试conf show 
    """
    #执行conf show
    args = ["conf","show","-r","lm://.","-c","c"]
    response = runApp(args,["1","2","3","111","q"])
    #获取main
    if not '19003' in response:
        return u"未查询到端口配置"
    #获取var
    if not '192.168.5.79:1000' in response:
        return u"未查询到redis配置"

    #输入无效数字
    if not u'输入的数字无效' in response:
        return u"未查询到无效的输入"
    
    return 


@test(testName)
def test_Conf_Install_Normal():
    #执行conf install
    args = ["conf","install","-r",ZKAddress,"-c","t"]
    response = runApp(args)

   #获取main
    if not 'OK' in response:
        return u"安装到配置中心失败"

    



@test(testName)
def test_Conf_Install_cover():
    #执行conf install -cover
    args = ["conf","install","-r",ZKAddress,"-c","c","-cover","true"]
    response = runApp(args,exe_name="new")

    if not 'OK' in response:
        return u"覆盖-安装到配置中心失败"

    args = ["conf","show","-r",ZKAddress]
    response = runApp(args,["3"],exe_name="new")
    if not '192.168.5.79:6379' in response:
        return u"覆盖-安装到配置中心失败"
 
    pass

def main():
     execute(testName)
 
if __name__ == "__main__":
    main()