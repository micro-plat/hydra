package update

import (
	"os"
	"strings"

	"github.com/micro-plat/hydra/compatible"
	"github.com/micro-plat/lib4go/ps"
)

//Restart 重启当前服务
func restart() (err error) {
	var name string
	cpid := os.Getpid()
	pids, err := getProcess(name)
	if err != nil {
		return err
	}
	for _, pid := range pids {
		if pid != cpid {
			killProcess(pid)
		}
	}
	return nil
}

// func close(){

// 	pid,_ := os.FindProcess(0)

// }

func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	process.Signal(compatible.CmdsUpdateProcessSignal)
	_, err = process.Wait()
	if err != nil {
		return err
	}
	return nil
}

//getProcess 根据进程名称获取所有进程编号
func getProcess(name string) ([]int, error) {
	pss, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	pids := make([]int, 0, 1)
	for _, v := range pss {
		if strings.EqualFold(v.Executable(), name) {
			pids = append(pids, v.Pid())
		}
	}
	return pids, nil
}
