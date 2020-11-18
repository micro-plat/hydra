package pkgs

import (
	"fmt"
	"strings"

	"github.com/micro-plat/lib4go/errs"
)

const (
	SUCCESS = "\t\t\t\t\t[  \033[32mOK\033[0m  ]" // Show colored "OK"
	FAILED  = "\t\t\t\t\t[\033[31mFAILED\033[0m]" // Show colored "FAILED"
)
const msgtemplate = "Install %s%s%v"

//GetCmdsResult  GetCmdsResult
func GetCmdsResult(serviceName, action string, err error, args ...string) error {
	if err != nil {
		return fmt.Errorf("%s %s %s:%w", action, serviceName, FAILED, err)
	}
	if len(args) > 0 {
		serviceName = serviceName + " " + strings.Join(args, " ")
	}
	return errs.NewIgnoreError(0, fmt.Sprintf("%s %s %s", action, serviceName, SUCCESS))
}
