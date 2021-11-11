/*
 * @Description:
 * @Autor: taoshouyin
 * @Date: 2021-09-18 11:08:07
 * @LastEditors: taoshouyin
 * @LastEditTime: 2021-09-22 17:14:22
 */
package dbr

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/micro-plat/lib4go/assert"
	"github.com/micro-plat/lib4go/types"
)

func TestDBR_CreatePersistentNode(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		path     string
		data     string
		wantErr  bool
	}{
		// {name: "mysql添加节点-新节点添加", provider: "mysql", path: "/test/CreatePersistentNode", data: `{"node":1}`, wantErr: false},
		// {name: "mysql添加节点-重复节点更新", provider: "mysql", path: "/test/CreatePersistentNode", data: `{"node":2}`, wantErr: true},
		// {name: "mysql添加节点-节点和值都相同", provider: "mysql", path: "/test/CreatePersistentNode", data: `{"node":1}`, wantErr: true},
		// {name: "mysql添加节点-特殊字符", provider: "mysql", path: "/test/CreatePersistentNode1", data: `{"node":"~!@#$%^&*(){}|":?><\[]"}`, wantErr: false},
		{name: "oracle添加节点-新节点添加", provider: "oracle", path: "/test/CreatePersistentNode", data: `{"node":1}`, wantErr: false},
		{name: "oracle添加节点-重复节点更新", provider: "oracle", path: "/test/CreatePersistentNode", data: `{"node":2}`, wantErr: true},
		{name: "oracle添加节点-节点和值都相同", provider: "oracle", path: "/test/CreatePersistentNode", data: `{"node":1}`, wantErr: true},
		{name: "oracle添加节点-特殊字符", provider: "oracle", path: "/test/CreatePersistentNode1", data: `{"node":"~!@#$%^&*(){}|":?><\[]"}`, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// var (
			// 	xdb    *xdb.DB
			// 	sqlStr *sqltexture
			// )
			r, err := getRegistryForTest(tt.provider)

			// if tt.provider == "mysql" {
			// 	xdb = mysql.NewBy("hbsv2x_dev", "123456dev", "192.168.0.36:3306", "hbsv2x_dev", db.WithConnect(10, 6, 900))
			// 	sqlStr = &mysqltexture
			// }
			// if tt.provider == "oracle" {
			// 	xdb = oracle.NewBy("ims17_v1_dev", "123456dev", "orcl136", db.WithConnect(10, 6, 900))
			// 	sqlStr = &oracletexture
			// }
			//r, err := NewDBR(xdb, sqlStr, nil)
			assert.Equal(t, false, (err != nil), "获取数据库对象失败", err)
			err = r.CreatePersistentNode(tt.path, tt.data)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("创建节点失败,%v", err))
			if !tt.wantErr {
				data, _, err := r.GetValue(tt.path)
				assert.Equal(t, tt.wantErr, (err != nil), "获取创建的节点失败", err)
				assert.Equal(t, tt.data, string(data), "获取的节点信息不想等")
			}
		})
	}
}

func TestDBR_CreatePersistentNode1(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		path     string
		data     string
		wantErr  bool
	}{
		{name: "mysql添加节点-断开链接重连后节点是否依然存在", provider: "mysql", wantErr: false},
		{name: "oracle添加节点-断开链接重连后节点是否依然存在", provider: "oracle", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), "获取数据库注册中心失败", err)
			tt.path = fmt.Sprintf("%d", time.Now().Unix())
			tt.data = time.Now().String()
			err = rgt.CreatePersistentNode(tt.path, tt.data)
			assert.Equal(t, tt.wantErr, (err != nil), "创建节点失败", err)
			rgt.Close()
			time.Sleep(10 * time.Second)
			rgt1, _ := getRegistryForTest(tt.provider)
			data, _, err := rgt1.GetValue(tt.path)
			assert.Equal(t, tt.wantErr, (err != nil), "获取创建的永久节点失败", err)
			assert.Equal(t, tt.data, string(data), "获取的节点信息不想等")
			rgt1.Close()
		})
	}
}

func TestDBR_CreateTempNode(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		path     string
		data     string
		wantErr  bool
	}{
		{name: "mysql添加临时节点-断开链接重连后节点是否被删除", provider: "mysql", wantErr: false},
		{name: "oracle添加临时节点-断开链接重连后节点是否被删除", provider: "oracle", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgt, err := getRegistryForTest(tt.provider)
			assert.Equal(t, false, (err != nil), "获取数据库注册中心失败", err)
			tt.path = fmt.Sprintf("/test/CreateTempNode/%d", time.Now().UnixNano())
			tt.data = tt.path
			err = rgt.CreateTempNode(tt.path, tt.data)
			assert.Equal(t, tt.wantErr, (err != nil), "创建临时节点失败", err)
			rgt.Close()
			time.Sleep(10 * time.Second)
			rgt1, _ := getRegistryForTest(tt.provider)
			b, err := rgt1.Exists(tt.path)
			assert.Equal(t, tt.wantErr, (err != nil), fmt.Sprintf("检查创建的临时节点是否存在失败,%v", err))
			assert.Equal(t, false, b, "创建的临时节点没有被删除")
			rgt1.Close()
		})
	}
}

func TestDBR_CreateSeqNode(t *testing.T) {
	//mysql顺序创建序列节点
	rgt, err := getRegistryForTest(MYSQL)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败y,%v", err))
	createPath := map[string]bool{}
	defer func() {
		for k, _ := range createPath {
			if k != "" {
				rgt.Delete(k)
			}
		}
		rgt.Close()
	}()

	data := fmt.Sprintf("%d", time.Now().UnixNano())
	path := fmt.Sprintf("/test/CreateSeqNode/%s", data)
	realPath := path + "/xyz"
	gotRpath, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
	createPath[gotRpath] = true
	pArry := strings.Split(gotRpath, "_")
	assert.Equal(t, 2, len(pArry), "创建的序列节点路径错误")
	nid := pArry[1]
	gotRpath1, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败1,%v", err))
	createPath[gotRpath1] = true
	pArry = strings.Split(gotRpath1, "_")
	assert.Equal(t, 2, len(pArry), "创建的序列节点路径错误")
	nid1 := pArry[1]
	assert.Equal(t, types.GetInt(nid1), types.GetInt(nid)+1, "创建的序列顺序不相同")
	cPaths, _, err := rgt.GetChildren(path)
	assert.Equal(t, 2, len(cPaths), "创建序列的子节点数不正确")

	//msyql并发创建序列节点
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(wg1 *sync.WaitGroup) {
				gpath, err := rgt.CreateSeqNode(realPath, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("并发获取数据库注册中心失败,%v", err))
				createPath[gpath] = true
				wg1.Done()
			}(&wg)
		}
		wg.Wait()
	}

	gotRpath2, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("并发获取数据库注册中心失败,%v", err))
	createPath[gotRpath2] = true
	pArry = strings.Split(gotRpath2, "_")
	assert.Equal(t, 2, len(pArry), "并发创建的序列节点路径错误")
	nid2 := pArry[1]
	assert.Equal(t, types.GetInt(nid2), types.GetInt(nid)+102, "创建的序列顺序不相同xx")
	cPaths, _, err = rgt.GetChildren(path)
	assert.Equal(t, 103, len(cPaths), "并发创建序列的子节点数不正确")
}

func TestDBROracleCreateSeqNode(t *testing.T) {
	//mysql顺序创建序列节点
	rgt, err := getRegistryForTest(ORACLE)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败y,%v", err))
	createPath := map[string]bool{}
	defer func() {
		for k, _ := range createPath {
			if k != "" {
				rgt.Delete(k)
			}
		}
		rgt.Close()
	}()

	data := fmt.Sprintf("%d", time.Now().UnixNano())
	path := fmt.Sprintf("/test/CreateSeqNode/%s", data)
	realPath := path + "/xyz"
	gotRpath, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败,%v", err))
	createPath[gotRpath] = true
	pArry := strings.Split(gotRpath, "_")
	assert.Equal(t, 2, len(pArry), "创建的序列节点路径错误")
	nid := pArry[1]
	gotRpath1, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("获取数据库注册中心失败1,%v", err))
	createPath[gotRpath1] = true
	pArry = strings.Split(gotRpath1, "_")
	assert.Equal(t, 2, len(pArry), "创建的序列节点路径错误")
	nid1 := pArry[1]
	assert.Equal(t, types.GetInt(nid1), types.GetInt(nid)+1, "创建的序列顺序不相同")
	cPaths, _, err := rgt.GetChildren(path)
	assert.Equal(t, 2, len(cPaths), "创建序列的子节点数不正确")

	//msyql并发创建序列节点
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			wg.Add(1)
			go func(wg1 *sync.WaitGroup) {
				gpath, err := rgt.CreateSeqNode(realPath, data)
				assert.Equal(t, false, (err != nil), fmt.Sprintf("并发获取数据库注册中心失败,%v", err))
				createPath[gpath] = true
				wg1.Done()
			}(&wg)
		}
		wg.Wait()
	}

	gotRpath2, err := rgt.CreateSeqNode(realPath, data)
	assert.Equal(t, false, (err != nil), fmt.Sprintf("并发获取数据库注册中心失败,%v", err))
	createPath[gotRpath2] = true
	pArry = strings.Split(gotRpath2, "_")
	assert.Equal(t, 2, len(pArry), "并发创建的序列节点路径错误")
	nid2 := pArry[1]
	assert.Equal(t, types.GetInt(nid2), types.GetInt(nid)+102, "创建的序列顺序不相同xx")
	cPaths, _, err = rgt.GetChildren(path)
	assert.Equal(t, 103, len(cPaths), "并发创建序列的子节点数不正确")
}
