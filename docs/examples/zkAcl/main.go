package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/micro-plat/hydra/registry"
	_ "github.com/micro-plat/hydra/registry/zookeeper"
	"github.com/micro-plat/lib4go/logger"
)

func main() {
	logger := logger.New("zkAcl")
	defer time.Sleep(time.Second * 2)

	if len(os.Args) < 2 {
		logger.Error("请输入zookeeper服务器地址")
		return
	}

	registry, err := registry.NewRegistryWithAddress(os.Args[1], logger)
	if err != nil {
		logger.Error(err)
		return
	}

	root := ""
	for {
		time.Sleep(time.Second * 1)
		fmt.Print("zk>")
		var cmd string
		var path string
		_, err := fmt.Scanf("%s %s", &cmd, &path)
		if err != nil {
			logger.Error(err)
			continue
		}
		switch cmd {
		case "list":
			if path == "" {
				path = "/"
			}
			if path == "/" {
				root = ""
			}
			cpath := fmt.Sprintf("%s/%s", strings.TrimRight(root, "/"), strings.Trim(path, "/"))
			paths, _, err := registry.GetChildren(cpath)
			if err != nil {
				logger.Error("list error:", err)
				continue
			}
			root = cpath
			fmt.Println("----list-----")
			fmt.Println(strings.Join(paths, "\r\n"))
			fmt.Println("-------------")
		case "value":
			cpath := fmt.Sprintf("%s/%s", strings.TrimRight(root, "/"), strings.Trim(path, "/"))
			buffer, _, err := registry.GetValue(cpath)
			if err != nil {
				logger.Error("value error", err)
				continue
			}

			fmt.Println(string(buffer))
		case "del":
			err = registry.Delete(path)
			if err != nil {
				logger.Error(err)
				continue
			}
			fmt.Println("删除成功:", path)

		case "pwd":
			fmt.Println(root)

		}

	}
}
