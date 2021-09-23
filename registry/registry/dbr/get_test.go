/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 16:51:16
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-22 14:14:27
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
		name         string
		provider     string
		subNodeCount int  //需要添加的子节点数
		nodeType     int  //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete     bool //是否删除父级节点
		isDeleteSub  int  //需要删除的子节点数
	}{
		{name: "mysql-永久节点-获取的节点不存在", provider: "mysql", isDelete: false, nodeType: 4},
		{name: "mysql-永久节点-节点存在没有子节点", provider: "mysql", nodeType: 1, subNodeCount: 0, isDelete: false, isDeleteSub: 0},
		{name: "mysql-永久节点-节点存在有一个字节点", provider: "mysql", nodeType: 1, subNodeCount: 1, isDelete: false, isDeleteSub: 0},
		{name: "mysql-永久节点-节点存在有多个子节点", provider: "mysql", nodeType: 1, subNodeCount: 3, isDelete: false, isDeleteSub: 0},
		{name: "mysql-永久节点-有子节点的父级节点被逻辑删除", provider: "mysql", nodeType: 1, subNodeCount: 1, isDelete: true, isDeleteSub: 0},
		{name: "mysql-永久节点-无子节点的父级节点被逻辑删除", provider: "mysql", nodeType: 1, subNodeCount: 0, isDelete: true, isDeleteSub: 0},
		{name: "mysql-永久节点-子节点存在被逻辑删除子节点", provider: "mysql", nodeType: 1, subNodeCount: 2, isDelete: false, isDeleteSub: 1},
		{name: "mysql-永久节点-子节点全都被逻辑删除", provider: "mysql", nodeType: 1, subNodeCount: 1, isDelete: false, isDeleteSub: 1},

		{name: "mysql-临时节点-获取的节点不存在", provider: "mysql", isDelete: false, nodeType: 4},
		{name: "mysql-临时节点-节点存在没有子节点", provider: "mysql", nodeType: 2, subNodeCount: 0, isDelete: false, isDeleteSub: 0},
		{name: "mysql-临时节点-节点存在有一个字节点", provider: "mysql", nodeType: 2, subNodeCount: 1, isDelete: false, isDeleteSub: 0},
		{name: "mysql-临时节点-节点存在有多个子节点", provider: "mysql", nodeType: 2, subNodeCount: 3, isDelete: false, isDeleteSub: 0},
		{name: "mysql-临时节点-有子节点的父级节点被逻辑删除", provider: "mysql", nodeType: 2, subNodeCount: 1, isDelete: true, isDeleteSub: 0},
		{name: "mysql-临时节点-无子节点的父级节点被逻辑删除", provider: "mysql", nodeType: 2, subNodeCount: 0, isDelete: true, isDeleteSub: 0},
		{name: "mysql-临时节点-子节点存在被逻辑删除子节点", provider: "mysql", nodeType: 2, subNodeCount: 2, isDelete: false, isDeleteSub: 1},
		{name: "mysql-临时节点-子节点全都被逻辑删除", provider: "mysql", nodeType: 2, subNodeCount: 1, isDelete: false, isDeleteSub: 1},

		// {name: "oracle-永久节点-获取的节点不存在", provider: "oracle", isDelete: false, nodeType: 4},
		// {name: "oracle-永久节点-节点存在没有子节点", provider: "oracle", nodeType: 1, subNodeCount: 0, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-永久节点-节点存在有一个字节点", provider: "oracle", nodeType: 1, subNodeCount: 1, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-永久节点-节点存在有多个子节点", provider: "oracle", nodeType: 1, subNodeCount: 3, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-永久节点-有子节点的父级节点被逻辑删除", provider: "oracle", nodeType: 1, subNodeCount: 1, isDelete: true, isDeleteSub: 0},
		// {name: "oracle-永久节点-无子节点的父级节点被逻辑删除", provider: "oracle", nodeType: 1, subNodeCount: 0, isDelete: true, isDeleteSub: 0},
		// {name: "oracle-永久节点-子节点存在被逻辑删除子节点", provider: "oracle", nodeType: 1, subNodeCount: 2, isDelete: false, isDeleteSub: 1},
		// {name: "oracle-永久节点-子节点全都被逻辑删除", provider: "oracle", nodeType: 1, subNodeCount: 1, isDelete: false, isDeleteSub: 1},

		// {name: "oracle-临时节点-获取的节点不存在", provider: "oracle", isDelete: false, nodeType: 4},
		// {name: "oracle-临时节点-节点存在没有子节点", provider: "oracle", nodeType: 2, subNodeCount: 0, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-临时节点-节点存在有一个字节点", provider: "oracle", nodeType: 2, subNodeCount: 1, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-临时节点-节点存在有多个子节点", provider: "oracle", nodeType: 2, subNodeCount: 3, isDelete: false, isDeleteSub: 0},
		// {name: "oracle-临时节点-有子节点的父级节点被逻辑删除", provider: "oracle", nodeType: 2, subNodeCount: 1, isDelete: true, isDeleteSub: 0},
		// {name: "oracle-临时节点-无子节点的父级节点被逻辑删除", provider: "oracle", nodeType: 2, subNodeCount: 0, isDelete: true, isDeleteSub: 0},
		// {name: "oracle-临时节点-子节点存在被逻辑删除子节点", provider: "oracle", nodeType: 2, subNodeCount: 2, isDelete: false, isDeleteSub: 1},
		// {name: "oracle-临时节点-子节点全都被逻辑删除", provider: "oracle", nodeType: 2, subNodeCount: 1, isDelete: false, isDeleteSub: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			defer rgt.Close()
			time.Sleep(1100 * time.Millisecond)
			data := fmt.Sprintf("%d", time.Now().Unix())
			path := fmt.Sprintf("/TestGetChildren/P%s", data)
			if tt.nodeType == 4 {
				gotPaths, gotVersion, err := rgt.GetChildren(path)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("获取子节点列表失败11,%v", err))
				assert.Equal(t, int32(0), gotVersion, fmt.Sprintf("节点版本号错误11,%v", gotVersion))
				assert.Equal(t, 0, len(gotPaths), fmt.Sprintf("子节点数量错误11,%v", gotVersion))
				return
			}

			subNodeMap := []string{}
			//创建父级节点
			if tt.nodeType == 1 {
				err = rgt.CreatePersistentNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建永久节点异常11,%v", err))
			}
			if tt.nodeType == 2 {
				err = rgt.CreateTempNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建临时节点异常11,%v", err))
			}
			if tt.nodeType == 3 {
				path, err = rgt.CreateSeqNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建序列节点异常11,%v", err))
			}

			//创建子节点
			for i := 0; i < tt.subNodeCount; i++ {
				subdata := fmt.Sprintf("%d", time.Now().UnixNano())
				subpath := fmt.Sprintf("%s/%s", path, subdata)
				if tt.nodeType == 1 {
					err = rgt.CreatePersistentNode(subpath, subdata)
					assert.Equal(t, false, (err != nil), fmt.Sprintf("创建永久节点异常11,%v", err))
				}
				if tt.nodeType == 2 {
					err = rgt.CreateTempNode(subpath, subdata)
					assert.Equal(t, false, (err != nil), fmt.Sprintf("创建临时节点异常11,%v", err))
				}
				if tt.nodeType == 3 {
					realpath, err := rgt.CreateSeqNode(subpath, subdata)
					assert.Equal(t, false, (err != nil), fmt.Sprintf("创建序列节点异常11,%v", err))
					subpath = realpath
				}
				subNodeMap = append(subNodeMap, subpath)
			}

			//删除父级节点
			if tt.isDelete {
				err = rgt.Delete(path)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("删除主节点异常11,%v", err))
			}

			//删除子节点
			for i := 0; i < tt.isDeleteSub; i++ {
				err = rgt.Delete(subNodeMap[i])
				assert.Equal(t, false, (err != nil), fmt.Sprintf("删除子节点异常11,%v", err))
			}

			//获取子节点列表
			gotPaths, gotVersion, err := rgt.GetChildren(path)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取子节点列表失败2,%v", err))
			if tt.isDelete {
				assert.Equal(t, int32(0), gotVersion, fmt.Sprintf("节点版本号错误2,%v", gotVersion))
				assert.Equal(t, 0, len(gotPaths), fmt.Sprintf("子节点数量错误2,%v", len(gotPaths)))
				return
			}
			assert.Equal(t, int32(1), gotVersion, fmt.Sprintf("节点版本号错误3,%v", gotVersion))
			assert.Equal(t, (tt.subNodeCount - tt.isDeleteSub), len(gotPaths), fmt.Sprintf("子节点数量错误3,%v", gotPaths))
		})
	}
}
