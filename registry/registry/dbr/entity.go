package dbr

import (
	"strings"

	"github.com/micro-plat/lib4go/types"
)

type input map[string]interface{}

func newInput(path string) input {
	return map[string]interface{}{
		FieldPath: path,
	}
}

func newInputByWatch(sec int, path ...string) input {
	return map[string]interface{}{
		FieldPath: `"` + strings.Join(path, `","`) + `"`,
		"sec":     sec,
	}
}

func newInputByUpdate(path string, value string, version int32) input {
	return map[string]interface{}{
		FieldPath:        path,
		FieldValue:       value,
		FieldDataVersion: version,
	}
}

func newInputByInsert(path string, value string, temp bool) input {
	return map[string]interface{}{
		FieldPath:        path,
		FieldIsTemp:      types.DecodeInt(temp, true, 0, 1),
		FieldValue:       value,
		FieldDataVersion: 1,
	}
}

//FieldCreateTime 字段创建时间的数据库名称
const FieldCreateTime = "create_time"

//FieldDataVersion 字段数据版本号的数据库名称
const FieldDataVersion = "data_version"

//FieldIsDelete 字段已删除的数据库名称
const FieldIsDelete = "is_delete"

//FieldIsTemp 字段临时节点的数据库名称
const FieldIsTemp = "is_temp"

//FieldPath 字段路径的数据库名称
const FieldPath = "path"

//FieldUpdateTime 字段更新时间的数据库名称
const FieldUpdateTime = "update_time"

//FieldValue 字段内容的数据库名称
const FieldValue = "value"