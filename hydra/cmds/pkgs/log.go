package pkgs

import (
	"fmt"
	"strings"

	"github.com/micro-plat/hydra/global/compatible"

	"github.com/micro-plat/lib4go/errs"
)

//GetCmdsResult  GetCmdsResult
func GetCmdsResult(serviceName, action string, err error, args ...string) error {
	if err != nil {
		return fmt.Errorf("%s %s %s:%w", action, serviceName, compatible.FAILED, err)
	}
	if len(args) > 0 {
		serviceName = serviceName + " " + strings.Join(args, " ")
	}
	return errs.NewIgnoreError(0, fmt.Sprintf("%s %s %s", action, serviceName, compatible.SUCCESS))
}
