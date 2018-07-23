// Example of a daemon with echo service
package hydra

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zkfy/daemon"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{}
var errlog = log.New(os.Stderr, "", 0)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage(name string, run func() (string, error)) (string, error) {

	usage := fmt.Sprintf("用法: %s install | remove | start | stop | status", name)

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	return run()

}

func dStart(run func() (string, error)) {
	name := filepath.Base(os.Args[0])
	desc := "基于hydra的微服务应用"
	srv, err := daemon.New(name, desc, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage(name, run)
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
