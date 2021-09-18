/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 16:51:16
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-18 18:10:57
 */
package dbr

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/types"
)

func TestDBR_GetValue(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		wantVersion int
		nodeType    int //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete    bool
		wantErr     bool
	}{
		{name: "mysql-获取的节点不存在", provider: "mysql", isDelete: false, nodeType: 4, wantErr: true},
		{name: "mysql-获取没有逻辑删除永久节点数据", provider: "mysql", wantVersion: 1, nodeType: 1, isDelete: false, wantErr: false},
		{name: "mysql-获取没有逻辑删除临时节点数据", provider: "mysql", wantVersion: 2, nodeType: 2, isDelete: false, wantErr: false},
		{name: "mysql-获取没有逻辑删除序列节点数据", provider: "mysql", wantVersion: 3, nodeType: 3, isDelete: false, wantErr: false},
		{name: "mysql-获取逻辑删除永久节点数据", provider: "mysql", wantVersion: 1, nodeType: 1, isDelete: true, wantErr: true},
		{name: "mysql-获取逻辑删除临时节点数据", provider: "mysql", wantVersion: 1, nodeType: 2, isDelete: true, wantErr: true},
		{name: "mysql-获取逻辑删除序列节点数据", provider: "mysql", wantVersion: 1, nodeType: 3, isDelete: true, wantErr: true},

		// {name: "oracle-获取的节点不存在", provider: "oracle", isDelete: false, nodeType: 4, wantErr: false},
		// {name: "oracle-获取没有逻辑删除永久节点数据", provider: "oracle", wantVersion: 1, nodeType: 1, isDelete: false, wantErr: false},
		// {name: "oracle-获取没有逻辑删除临时节点数据", provider: "oracle", wantVersion: 2, nodeType: 2, isDelete: false, wantErr: false},
		// {name: "oracle-获取没有逻辑删除序列节点数据", provider: "oracle", wantVersion: 3, nodeType: 3, isDelete: false, wantErr: false},
		// {name: "oracle-获取逻辑删除永久节点数据", provider: "oracle", wantVersion: 1, nodeType: 1, isDelete: true, wantErr: true},
		// {name: "oracle-获取逻辑删除临时节点数据", provider: "oracle", wantVersion: 1, nodeType: 2, isDelete: true, wantErr: true},
		// {name: "oracle-获取逻辑删除序列节点数据", provider: "oracle", wantVersion: 1, nodeType: 3, isDelete: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/test/TestDBRGetValue/%s", data)
			if tt.nodeType == 4 {
				_, _, err := rgt.GetValue(path)
				assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("查询节点数据异常,%v", err))
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
			for i := 1; i < tt.wantVersion; i++ {
				data = types.GetString(time.Now().UnixNano())
				err = rgt.Update(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("更新节点数据异常,%v", err))
			}

			if tt.isDelete {
				rgt.Delete(path)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("删除节点异常,%v", err))
			}

			gotData, gotVersion, err := rgt.GetValue(path)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("查询节点数据异常,%v", err))
			if !tt.wantErr {
				assert.Equal(t, data, string(gotData), fmt.Sprintf("查询节点数据异常,%v", err))
				assert.Equal(t, int32(tt.wantVersion), gotVersion, fmt.Sprintf("查询节点数据异常,%v", err))
			}
			rgt.Close()
		})
	}
}

func TestDBR_GetChildren(t *testing.T) {

	tests := []struct {
		name        string
		provider    string
		wantVersion int
		nodeType    int //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete    bool
		wantErr     bool
	}{
		{name: "mysql-获取的节点不存在", provider: "mysql", isDelete: false, nodeType: 4, wantErr: true},
		{name: "mysql-节点存在没有子节点", provider: "mysql", wantVersion: 1, nodeType: 1, isDelete: false, wantErr: false},
		{name: "mysql-节点存在有一个字节点", provider: "mysql", wantVersion: 2, nodeType: 1, isDelete: false, wantErr: false},
		{name: "mysql-节点存在有多个子节点", provider: "mysql", wantVersion: 3, nodeType: 1, isDelete: false, wantErr: false},
		{name: "mysql-有子节点的父级节点被逻辑删除", provider: "mysql", wantVersion: 3, nodeType: 1, isDelete: true, wantErr: true},
		{name: "mysql-无字节点的父级节点被逻辑删除", provider: "mysql", wantVersion: 2, nodeType: 1, isDelete: true, wantErr: true},
		{name: "mysql-子节点存在逻辑删除子节点", provider: "mysql", wantVersion: 2, nodeType: 1, isDelete: true, wantErr: true},
		{name: "mysql-子节点全都被逻辑删除", provider: "mysql", wantVersion: 1, nodeType: 1, isDelete: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/test/TestDBRGetChildren/%s", data)
			//gotPaths, gotVersion, err :=
			rgt.GetChildren(path)

		})
	}
}
