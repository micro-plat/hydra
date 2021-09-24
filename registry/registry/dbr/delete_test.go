/*
 * @Autor: taoshouyin
 * @Date: 2021-09-18 16:06:07
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-18 18:17:35
 */
package dbr

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
)

//逻辑删除节点数据
func TestDBR_Delete(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		isExist  bool
		isTemp   bool
		wantErr  bool
	}{
		{name: "mysql-删除存在的临时节点", provider: "mysql", isTemp: true, isExist: true, wantErr: false},
		{name: "mysql-删除不存在的临时节点", provider: "mysql", isTemp: true, isExist: false, wantErr: true},
		{name: "mysql-删除存在的永久节点", provider: "mysql", isTemp: false, isExist: true, wantErr: false},
		{name: "mysql-删除不存在的永久节点", provider: "mysql", isTemp: false, isExist: false, wantErr: true},
		{name: "mysql-删除父级节点,子节点也要被删除", provider: "mysql", isTemp: false, isExist: false, wantErr: true},

		{name: "oracle-删除存在的临时节点", provider: "oracle", isTemp: true, isExist: true, wantErr: false},
		{name: "oracle-删除不存在的临时节点", provider: "oracle", isTemp: true, isExist: false, wantErr: true},
		{name: "oracle-删除存在的永久节点", provider: "oracle", isTemp: false, isExist: true, wantErr: false},
		{name: "oracle-删除不存在的永久节点", provider: "oracle", isTemp: false, isExist: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), "获取数据库注册中心失败", err)
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/test/TestDBRDelete/%s", data)
			if tt.isExist {
				if tt.isTemp {
					err = rgt.CreateTempNode(path, data)
					assert.Equal(t, false, (err != nil), fmt.Sprintf("删除时创建临时节点异常,%v", err))
				} else {
					err = rgt.CreatePersistentNode(path, data)
					assert.Equal(t, false, (err != nil), fmt.Sprintf("删除时创建永久节点异常,%v", err))
				}
			}

			err = rgt.Delete(path)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("删除节点数据异常,%v", err))
			b, err := rgt.Exists(path)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("检查节点是否存在异常,%v", err))
			assert.Equal(t, false, b, "节点还存在,没有删除成功")
		})
	}
}

//物理删除-更新时间30秒内被逻辑删除的临时节点--私有无法测试
func TestDBR_clear(t *testing.T) {
	tests := []struct {
		name       string
		provider   string
		isTemp     bool
		clearcount int
	}{
		{name: "mysql-清理没有删除的临时节点", provider: "mysql", isTemp: true, clearcount: 1},
		{name: "mysql-清理已经被删除的临时节点", provider: "mysql", isTemp: true, clearcount: 2},
		{name: "mysql-清理已经被删除超过30s的临时节点", provider: "mysql", isTemp: true, clearcount: 2},
		{name: "mysql-清理没有删除的永久节点", provider: "mysql", isTemp: false, clearcount: 3},
		{name: "mysql-清理已经被删除的永久节点", provider: "mysql", isTemp: false, clearcount: 3},
		{name: "mysql-清理已经被删除超过30s的永久节点", provider: "mysql", isTemp: false, clearcount: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// rgt, err := getRegistryForTest(tt.provider)
			// assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			// rgt.clear("path")
		})
	}
}
