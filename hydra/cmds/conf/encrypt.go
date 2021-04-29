package conf

import (
	"fmt"

	"github.com/micro-plat/hydra/conf/pkgs/security"
	"github.com/urfave/cli"
)

func encrypt(c *cli.Context) (err error) {
	if len(orgData) == 0 {
		return fmt.Errorf("未指定加密的内容")
	}

	cipherData := security.Encrypt([]byte(orgData))
	fmt.Println("原始内容：")
	fmt.Println(orgData)
	fmt.Println("加密结果：")
	fmt.Println(cipherData)
	return nil
}
