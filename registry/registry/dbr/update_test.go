/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-22 14:21:51
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-22 14:52:30
 */
package dbr

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
)

func TestDBR_Update(t *testing.T) {
	tests := []struct {
		name        string
		provider    string
		wantVersion int32
		nodeType    int //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete    bool
		isChange    bool
	}{
		{name: "mysql-更新不存在的节点", provider: "mysql", nodeType: 4, wantVersion: 0},
		{name: "mysql-永久节点-更新已经被逻辑删除的节点", provider: "mysql", nodeType: 1, isDelete: true, isChange: true, wantVersion: 0},
		{name: "mysql-永久节点-更新节点,但数据没有变化", provider: "mysql", nodeType: 1, isDelete: false, isChange: false, wantVersion: 2},
		{name: "mysql-永久节点-更新节点,数据发生变化", provider: "mysql", nodeType: 1, isDelete: false, isChange: true, wantVersion: 2},
		{name: "mysql-永久节点-多次更新节点", provider: "mysql", nodeType: 1, isDelete: false, isChange: true, wantVersion: 4},
		{name: "mysql-临时节点-更新已经被逻辑删除的节点", provider: "mysql", nodeType: 2, isDelete: true, isChange: true, wantVersion: 0},
		{name: "mysql-临时节点-更新节点,但数据没有变化", provider: "mysql", nodeType: 2, isDelete: false, isChange: false, wantVersion: 2},
		{name: "mysql-临时节点-更新节点,数据发生变化", provider: "mysql", nodeType: 2, isDelete: false, isChange: true, wantVersion: 2},
		{name: "mysql-临时节点-多次更新节点", provider: "mysql", nodeType: 2, isDelete: false, isChange: true, wantVersion: 4},
		{name: "mysql-序列节点-更新已经被逻辑删除的节点", provider: "mysql", nodeType: 3, isDelete: true, isChange: true, wantVersion: 0},
		{name: "mysql-序列节点-更新节点,但数据没有变化", provider: "mysql", nodeType: 3, isDelete: false, isChange: false, wantVersion: 2},
		{name: "mysql-序列节点-更新节点,数据发生变化", provider: "mysql", nodeType: 3, isDelete: false, isChange: true, wantVersion: 2},
		{name: "mysql-序列节点-多次更新节点", provider: "mysql", nodeType: 3, isDelete: false, isChange: true, wantVersion: 4},

		// {name: "oracle-更新不存在的节点", provider: "oracle", nodeType: 4, wantVersion: 0},
		// {name: "oracle-永久节点-更新已经被逻辑删除的节点", provider: "oracle", nodeType: 1, isDelete: true, isChange: true, wantVersion: 0},
		// {name: "oracle-永久节点-更新节点,但数据没有变化", provider: "oracle", nodeType: 1, isDelete: false, isChange: false, wantVersion: 2},
		// {name: "oracle-永久节点-更新节点,数据发生变化", provider: "oracle", nodeType: 1, isDelete: false, isChange: true, wantVersion: 2},
		// {name: "oracle-永久节点-多次更新节点", provider: "oracle", nodeType: 1, isDelete: false, isChange: true, wantVersion: 4},
		// {name: "oracle-临时节点-更新已经被逻辑删除的节点", provider: "oracle", nodeType: 2, isDelete: true, isChange: true, wantVersion: 0},
		// {name: "oracle-临时节点-更新节点,但数据没有变化", provider: "oracle", nodeType: 2, isDelete: false, isChange: false, wantVersion: 2},
		// {name: "oracle-临时节点-更新节点,数据发生变化", provider: "oracle", nodeType: 2, isDelete: false, isChange: true, wantVersion: 2},
		// {name: "oracle-临时节点-多次更新节点", provider: "oracle", nodeType: 2, isDelete: false, isChange: true, wantVersion: 4},
		// {name: "oracle-序列节点-更新已经被逻辑删除的节点", provider: "oracle", nodeType: 3, isDelete: true, isChange: true, wantVersion: 0},
		// {name: "oracle-序列节点-更新节点,但数据没有变化", provider: "oracle", nodeType: 3, isDelete: false, isChange: false, wantVersion: 2},
		// {name: "oracle-序列节点-更新节点,数据发生变化", provider: "oracle", nodeType: 3, isDelete: false, isChange: true, wantVersion: 2},
		// {name: "oracle-序列节点-多次更新节点", provider: "oracle", nodeType: 3, isDelete: false, isChange: true, wantVersion: 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/TestUpdate/%s", data)
			if tt.nodeType == 4 {
				err := rgt.Update(path, data)
				assert.Equal(t, true, (err != nil), fmt.Sprintf("更新节点数据异常,%v", err))
				return
			}

			if tt.nodeType == 1 {
				err = rgt.CreatePersistentNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建永久节点异常1,%v", err))
			}
			if tt.nodeType == 2 {
				err = rgt.CreateTempNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建临时节点异常1,%v", err))
			}
			if tt.nodeType == 3 {
				path, err = rgt.CreateSeqNode(path, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("创建序列节点异常1,%v", err))
			}

			if tt.isDelete {
				err = rgt.Delete(path)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("删除节点异常1,%v", err))
			}

			data1 := data
			var i int32
			for i = 0; i < tt.wantVersion-1; i++ {
				if tt.isChange {
					data1 = fmt.Sprintf("%d", time.Now().UnixNano())
				}
				err := rgt.Update(path, data1)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("更新节点数据异常2,%v", err))
			}

			//获取节点信息
			content, versionNo, err := rgt.GetValue(path)
			if tt.isDelete {
				assert.Equal(t, true, (err != nil), fmt.Sprintf("更新节点数据异常2,%v", err))
				assert.Equal(t, tt.wantVersion, versionNo, fmt.Sprint("节点版本号错误2"))
				return
			}
			assert.Equal(t, false, (err != nil), fmt.Sprintf("更新节点数据异常2,%v", err))
			assert.Equal(t, tt.wantVersion, versionNo, fmt.Sprint("节点版本号错误2"))
			assert.Equal(t, data1, string(content), fmt.Sprint("节点内容错误"))
		})
	}
}
