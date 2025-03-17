package utility

import (
	"strings"

	"github.com/google/uuid"
)

// GetGUID 生成Guid字串
func GetGUID() string {
	uuidV4 := uuid.New()
	return strings.ReplaceAll(uuidV4.String(), "-", "")
}
