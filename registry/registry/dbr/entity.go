package dbr

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type input map[string]interface{}

func newInput(path string) input {
	return map[string]interface{}{
		"path": path,
	}
}

func newInputByWatch(sec int, path ...string) input {
	return map[string]interface{}{
		"path": `"` + strings.Join(path, `","`) + `"`,
		"sec":  sec,
	}
}

func newInputByUpdate(path string, value string, version int32) input {
	return map[string]interface{}{
		"path":         path,
		"value":        value,
		"data_version": version,
	}
}

func newInputByInsert(path string, value string, temp bool) input {
	return map[string]interface{}{
		"path":         path,
		"temp":         types.DecodeInt(temp, true, 1, 0),
		"value":        value,
		"data_version": 1,
		"acl_version":  1,
	}
}
