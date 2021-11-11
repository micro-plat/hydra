/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-22 14:56:12
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-24 15:57:36
 */
package dbr

import (
	"fmt"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
)

func TestDBR_WatchValue(t *testing.T) {

	tests := []struct {
		name     string
		provider string
		nodeType int //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete bool
		isNotify bool
	}{
		{name: "mysql-值监控节点不存在", provider: "mysql", nodeType: 4, isDelete: false, isNotify: false},
		{name: "mysql-永久节点-监控已经被逻辑删除的节点", provider: "mysql", nodeType: 1, isDelete: true, isNotify: false},
		{name: "mysql-永久节点-正常节点监控", provider: "mysql", nodeType: 1, isDelete: false, isNotify: true},
		{name: "mysql-临时节点-监控已经被逻辑删除的节点", provider: "mysql", nodeType: 2, isDelete: true, isNotify: false},
		{name: "mysql-临时节点-正常节点监控", provider: "mysql", nodeType: 2, isDelete: false, isNotify: true},
		{name: "mysql-序列节点-监控已经被逻辑删除的节点", provider: "mysql", nodeType: 3, isDelete: true, isNotify: false},
		{name: "mysql-序列节点-正常节点监控", provider: "mysql", nodeType: 3, isDelete: false, isNotify: true},

		// {name: "oracle-值监控节点不存在", provider: "oracle", nodeType: 4, isDelete: false, isNotify: false},
		// {name: "oracle-永久节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 1, isDelete: true, isNotify: false},
		// {name: "oracle-永久节点-正常节点监控", provider: "oracle", nodeType: 1, isDelete: false, isNotify: true},
		// {name: "oracle-临时节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 2, isDelete: true, isNotify: false},
		// {name: "oracle-临时节点-正常节点监控", provider: "oracle", nodeType: 2, isDelete: false, isNotify: true},
		// {name: "oracle-序列节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 3, isDelete: true, isNotify: false},
		// {name: "oracle-序列节点-正常节点监控", provider: "oracle", nodeType: 3, isDelete: false, isNotify: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			defer rgt.Close()
			data := fmt.Sprintf("%d", time.Now().UnixNano())
			path := fmt.Sprintf("/TestDBRWatchValue/%s", data)
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

			gotData, err := rgt.WatchValue(path)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("把节点添加到监控列表失败,%v", err))

			data = fmt.Sprintf("%d", time.Now().UnixMicro())
			err = rgt.Update(path, data)
			if tt.isNotify {
				assert.Equal(t, false, (err != nil), fmt.Sprintf("更新节点错误,%v", err))
			} else {
				assert.Equal(t, true, (err != nil), fmt.Sprintf("更新节点错误1,%v", err))
			}
			var (
				rPath string
				rData []byte
				// versionNo int32
			)
			forbool := true
			tk := time.NewTimer(time.Second * 3)
			for forbool {
				select {
				case <-tk.C:
					forbool = false
				case content, ok := <-gotData:

					if !ok {
						assert.Equal(t, true, ok, fmt.Sprintf("把节点添加到监控列表失败,%v", err))
						return
					}
					if err = content.GetError(); err != nil {
						assert.Equal(t, true, ok, fmt.Sprintf("节点变更通知内容错误,%v", err))
						return
					}
					rPath = content.GetPath()
					rData, _ = content.GetValue()
					forbool = false
				}
			}

			if tt.isNotify {
				assert.Equal(t, path, rPath, fmt.Sprintf("通知路进错误错误,%v", err))
				assert.Equal(t, data, string(rData), fmt.Sprintf("通知数据错误,%v", err))
				// assert.Equal(t, int32(2), versionNo, fmt.Sprintf("通知的版本号错误,%v", err))
			}
		})
	}
}

func TestDBR_WatchChildren1(t *testing.T) {
	tests := []struct {
		index    int64
		name     string
		provider string
		nodeType int //1.永久节点 2.临时节点  3.序列节点 4.物理不存在的节点
		isDelete bool
		isNotify bool
	}{
		{index: 1, name: "mysql-监控节点不存在", provider: "mysql", nodeType: 4, isDelete: false, isNotify: false},
		{index: 2, name: "mysql-永久节点-监控被逻辑删除的节点", provider: "mysql", nodeType: 1, isDelete: true, isNotify: false},
		{index: 3, name: "mysql-永久节点-正常节点监控", provider: "mysql", nodeType: 1, isDelete: false, isNotify: true},
		{index: 4, name: "mysql-临时节点-监控被逻辑删除的节点", provider: "mysql", nodeType: 2, isDelete: true, isNotify: false},
		{index: 5, name: "mysql-临时节点-正常节点监控", provider: "mysql", nodeType: 2, isDelete: false, isNotify: true},
		{index: 6, name: "mysql-序列节点-监控被逻辑删除的节点", provider: "mysql", nodeType: 3, isDelete: true, isNotify: false},
		{index: 7, name: "mysql-序列节点-正常节点监控", provider: "mysql", nodeType: 3, isDelete: false, isNotify: true},

		// {index: 8,name: "oracle-值监控节点不存在", provider: "oracle", nodeType: 4, isDelete: false, isNotify: false},
		// {index: 9,name: "oracle-永久节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 1, isDelete: true, isNotify: false},
		// {index: 10,name: "oracle-永久节点-正常节点监控", provider: "oracle", nodeType: 1, isDelete: false, isNotify: true},
		// {index: 11,name: "oracle-临时节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 2, isDelete: true, isNotify: false},
		// {index: 12,name: "oracle-临时节点-正常节点监控", provider: "oracle", nodeType: 2, isDelete: false, isNotify: true},
		// {index: 13,name: "oracle-序列节点-监控已经被逻辑删除的节点", provider: "oracle", nodeType: 3, isDelete: true, isNotify: false},
		// {index: 14,name: "oracle-序列节点-正常节点监控", provider: "oracle", nodeType: 3, isDelete: false, isNotify: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
			defer func() {
				rgt.Close()
			}()
			data := fmt.Sprintf("%d-%d", time.Now().UnixMicro(), tt.index)
			path := fmt.Sprintf("/WatchChildren/%s", data)
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

			gotData, err := rgt.WatchChildren(path)
			assert.Equal(t, false, (err != nil), fmt.Sprintf("把节点添加到监控列表失败,%v", err))
			timeD := time.Now().UnixMicro()
			for i := 0; i < 2; i++ {
				data1 := fmt.Sprintf("%d-%d", timeD, i)
				path1 := fmt.Sprintf("%s/%s", path, data1)
				go func(p string, nodeType int) {
					if nodeType == 1 {
						err = rgt.CreatePersistentNode(p, data1)
						if err != nil {
							fmt.Printf("创建永久节点异常2,%v\n", err)
						}
					}
					if nodeType == 2 {
						err = rgt.CreateTempNode(p, data1)
						if err != nil {
							fmt.Printf("创建临时节点异常2,%v\n", err)
						}
					}
					if nodeType == 3 {
						_, err = rgt.CreateSeqNode(p, data1)
						if err != nil {
							fmt.Printf("创建序列节点异常2,%v\n", err)
						}
					}
				}(path1, tt.nodeType)
			}

			var (
				rPath     string
				versionNo int32
			)
			forbool := true
			tk := time.NewTimer(time.Second * 3)
			for forbool {
				select {
				case <-tk.C:
					forbool = false
				case content, ok := <-gotData:
					if !ok {
						break
					}
					if err = content.GetError(); err != nil {
						assert.Equal(t, true, ok, fmt.Sprintf("节点变更通知内容错误,%v", err))
						return
					}
					rPath = content.GetPath()
					_, versionNo = content.GetValue()
					assert.Equal(t, path, rPath, fmt.Sprintf("通知路进错误错误,%v", err))
					assert.Equal(t, int32(1), versionNo, fmt.Sprintf("通知的版本号错误,%v", err))
				}
			}
		})
	}
}
