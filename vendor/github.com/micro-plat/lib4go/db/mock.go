// Package db 提供 Mock 数据库实现用于测试
package db

import (
	"database/sql"
	"errors"

	"github.com/stretchr/testify/mock"
)

// MockDB Mock 数据库实现
// 用于单元测试，模拟数据库操作
type MockDB struct {
	mock.Mock

	// 事务支持
	parent     *MockDB
	committed  bool
	rolledBack bool
}

// 确保 MockDB 实现 IDB 接口
var _ IDB = (*MockDB)(nil)
var _ IDBTrans = (*MockDB)(nil)

// NewMockDB 创建 Mock 数据库
func NewMockDB() *MockDB {
	return &MockDB{}
}

// ========== IDBExecuter 接口实现 ==========

// Query 执行查询操作
func (m *MockDB) Query(sql string, input map[string]interface{}) (QueryRows, error) {
	args := m.Mock.Called(sql, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(QueryRows), args.Error(1)
}

// QueryBatch 执行批量查询操作
func (m *MockDB) QueryBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	args := m.Mock.Called(sqls, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(QueryRows), args.Error(1)
}

// Scalar 执行标量查询
func (m *MockDB) Scalar(sql string, input map[string]interface{}) (interface{}, error) {
	args := m.Mock.Called(sql, input)
	return args.Get(0), args.Error(1)
}

// Execute 执行非查询 SQL
func (m *MockDB) Execute(sql string, input map[string]interface{}) (int64, error) {
	args := m.Mock.Called(sql, input)
	return args.Get(0).(int64), args.Error(1)
}

// Executes 执行 SQL 并返回自增 ID 和影响行数
func (m *MockDB) Executes(sql string, input map[string]interface{}) (int64, int64, error) {
	args := m.Mock.Called(sql, input)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(1)
}

// ExecuteBatch 执行批量 SQL
func (m *MockDB) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	args := m.Mock.Called(sqls, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(QueryRows), args.Error(1)
}

// InsertBatch 批量插入
func (m *MockDB) InsertBatch(sql string, inputs []map[string]interface{}) (int64, error) {
	args := m.Mock.Called(sql, inputs)
	return args.Get(0).(int64), args.Error(1)
}

// UpdateBatch 批量更新
func (m *MockDB) UpdateBatch(sql string, inputs []map[string]interface{}) (int64, error) {
	args := m.Mock.Called(sql, inputs)
	return args.Get(0).(int64), args.Error(1)
}

// FetchRows 获取原始数据行
// 注意：Mock 模式下不支持此方法，请使用 Query 方法代替
func (m *MockDB) FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error) {
	// 在 Mock 测试中，应该使用 Query 方法而不是 FetchRows
	// FetchRows 用于获取原始 *sql.Rows，这在 Mock 场景下没有意义
	return nil, errors.New("FetchRows is not supported in Mock mode, use Query instead")
}

// ========== IDB 接口实现 ==========

// Begin 开始事务
func (m *MockDB) Begin() (IDBTrans, error) {
	args := m.Mock.Called()
	trans := &MockDB{parent: m}
	if args.Get(0) != nil {
		return args.Get(0).(IDBTrans), args.Error(1)
	}
	return trans, args.Error(1)
}

// Close 关闭数据库连接
func (m *MockDB) Close() {
	m.Mock.Called()
}

// ExecuteSP 执行存储过程
func (m *MockDB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (int64, error) {
	args := m.Mock.Called(procName, input, output)
	return args.Get(0).(int64), args.Error(1)
}

// ========== IDBTrans 接口实现 ==========

// Rollback 回滚事务
func (m *MockDB) Rollback() error {
	if m.parent != nil {
		m.parent.rolledBack = true
	}
	args := m.Mock.Called()
	return args.Error(0)
}

// Commit 提交事务
func (m *MockDB) Commit() error {
	if m.parent != nil {
		m.parent.committed = true
	}
	args := m.Mock.Called()
	return args.Error(0)
}

// ========== 辅助方法 ==========

// IsCommitted 检查事务是否已提交
func (m *MockDB) IsCommitted() bool {
	return m.committed
}

// IsRolledBack 检查事务是否已回滚
func (m *MockDB) IsRolledBack() bool {
	return m.rolledBack
}
