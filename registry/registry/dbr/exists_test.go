/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 16:51:02
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-18 17:16:05
 */
package dbr

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
)

//不存在的节点和被逻辑删除的节点都将被任务节点不存在
func TestDBR_Exists(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		nodeType int  //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete bool //被逻辑删除
		wantErr  bool
		want     bool
	}{
		{name: "mysql-不存在的节点", provider: "mysql", nodeType: 4, isDelete: false, wantErr: false, want: false},
		{name: "mysql-没有被逻辑删除的临时节点", provider: "mysql", nodeType: 2, isDelete: false, wantErr: false, want: true},
		{name: "mysql-没有被逻辑删除的永久节点", provider: "mysql", nodeType: 1, isDelete: false, wantErr: false, want: true},
		{name: "mysql-没有被逻辑删除的序列节点", provider: "mysql", nodeType: 3, isDelete: false, wantErr: false, want: true},
		{name: "mysql-被逻辑删除的序列节点", provider: "mysql", nodeType: 2, isDelete: true, wantErr: false, want: false},
		{name: "mysql-被逻辑删除的序列节点", provider: "mysql", nodeType: 1, isDelete: true, wantErr: false, want: false},
		{name: "mysql-被逻辑删除的序列节点", provider: "mysql", nodeType: 3, isDelete: true, wantErr: false, want: false},

		// {name: "oracle-不存在的节点", provider: "oracle", nodeType: 4, isDelete: false, wantErr: false},
		// {name: "oracle-没有被逻辑删除的临时节点", provider: "oracle", nodeType: 2, isDelete: false, wantErr: false},
		// {name: "oracle-没有被逻辑删除的永久节点", provider: "oracle", nodeType: 1, isDelete: false, wantErr: false},
		// {name: "oracle-没有被逻辑删除的序列节点", provider: "oracle", nodeType: 3, isDelete: false, wantErr: false},
		// {name: "oracle-被逻辑删除的序列节点", provider: "oracle", nodeType: 2, isDelete: true, wantErr: false},
		// {name: "oracle-被逻辑删除的序列节点", provider: "oracle", nodeType: 1, isDelete: true, wantErr: false},
		// {name: "oracle-被逻辑删除的序列节点", provider: "oracle", nodeType: 3, isDelete: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/test/TestDBRExists/%s", data)
			if tt.nodeType == 4 {
				got, err := rgt.Exists(path)
				assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("检查节点是否存在异常,%v", err))
				assert.Equal(t, tt.want, got)
				return
			}
			if tt.nodeType == 1 {
				err = rgt.CreatePersistentNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建永久节点异常,%v", err))
			}
			if tt.nodeType == 2 {
				err = rgt.CreateTempNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建临时节点异常,%v", err))
			}
			if tt.nodeType == 3 {
				path, err = rgt.CreateSeqNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建序列节点异常,%v", err))
			}

			if tt.isDelete {
				rgt.Delete(path)
			}
			got, err := rgt.Exists(path)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("检查节点是否存在异常,%v", err))
			assert.Equal(t, tt.want, got)
			rgt.Close()
		})
	}
}
