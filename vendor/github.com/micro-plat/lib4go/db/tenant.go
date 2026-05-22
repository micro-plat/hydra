// Package db 多租户扩展：支持三种隔离模式
//   - ModelSharedSchema:   共享连接池 + SET search_path 切换 schema（推荐，同实例最优）
//   - ModelSharedRLS:      共享连接池 + PostgreSQL RLS 行级安全
//   - ModelIndependentSchema: 独立连接池 + 独立 schema（租户少时可用）
package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/micro-plat/lib4go/db/tpl"
)

// ========================= 租户模式定义 =========================

// TenantModel 租户数据隔离模式
type TenantModel string

const (
	// ModelSharedSchema 共享连接池 + SET search_path（推荐）
	// 同一数据库实例，所有租户共享连接池，每次操作前 SET search_path TO schema_name, public
	// 无需修改表结构，无需加 tenant_id 列，public schema 中的公共表自动可见
	ModelSharedSchema TenantModel = "shared_schema"

	// ModelSharedRLS 共享连接池 + PostgreSQL RLS 行级安全
	// 需要每张表加 tenant_id 列并创建 RLS 策略，侵入性较强
	ModelSharedRLS TenantModel = "shared_rls"

	// ModelIndependentSchema 独立连接池 + 独立 schema
	// 每个租户创建独立的连接池，租户数量多时会耗尽数据库连接数
	ModelIndependentSchema TenantModel = "indep_schema"
)

// TenantConfig 租户配置
type TenantConfig struct {
	TenantID   string
	Model      TenantModel
	SchemaName string // schema 名称（ModelSharedSchema 和 ModelIndependentSchema 使用）
	DSN        string // 独立模式专用：自定义连接串（为空则使用模板生成）
}

// TenantConfigLoader 租户配置加载接口
type TenantConfigLoader interface {
	Load(tenantID string) (*TenantConfig, error)
}

// ========================= 租户管理器 =========================

var (
	tenantMgr     *tenantManager
	tenantMgrOnce sync.Once
)

// tenantManager 管理多租户数据源（单例）
type tenantManager struct {
	sharedSysDB *SysDB
	sharedTpl   tpl.ITPLContext

	indepSysDBs        map[string]*SysDB
	indepTpls          map[string]tpl.ITPLContext
	configLoader       TenantConfigLoader
	defaultDriver      string
	defaultDSNTemplate string

	mu sync.RWMutex
}

// InitTenantManager 初始化多租户管理器（必须在程序启动时调用一次）
func InitTenantManager(sharedDriver, sharedDSN string, maxOpen, maxIdle, maxLifeTime int,
	defaultDriver, defaultDSNTemplate string) error {
	if tenantMgr != nil {
		return nil
	}
	sysDB, err := NewSysDB(sharedDriver, sharedDSN, maxOpen, maxIdle, time.Duration(maxLifeTime)*time.Second)
	if err != nil {
		return fmt.Errorf("创建共享数据源失败: %w", err)
	}
	tplCtx, err := tpl.GetDBContext(sharedDriver)
	if err != nil {
		return fmt.Errorf("获取模板上下文失败: %w", err)
	}
	tenantMgrOnce.Do(func() {
		tenantMgr = &tenantManager{
			sharedSysDB:        sysDB,
			sharedTpl:          tplCtx,
			indepSysDBs:        make(map[string]*SysDB),
			indepTpls:          make(map[string]tpl.ITPLContext),
			configLoader:       &defaultConfigLoader{configs: make(map[string]*TenantConfig)},
			defaultDriver:      defaultDriver,
			defaultDSNTemplate: defaultDSNTemplate,
		}
	})
	return nil
}

// SetTenantConfig 设置租户配置（如果使用默认加载器）
func SetTenantConfig(tenantID string, cfg *TenantConfig) {
	if tenantMgr == nil {
		panic("tenant manager not initialized, call InitTenantManager first")
	}
	if loader, ok := tenantMgr.configLoader.(*defaultConfigLoader); ok {
		loader.set(tenantID, cfg)
	} else {
		panic("cannot set tenant config: config loader is not default")
	}
}

// SetTenantConfigLoader 设置自定义配置加载器
func SetTenantConfigLoader(loader TenantConfigLoader) {
	if tenantMgr == nil {
		panic("tenant manager not initialized, call InitTenantManager first")
	}
	tenantMgr.mu.Lock()
	defer tenantMgr.mu.Unlock()
	tenantMgr.configLoader = loader
}

// CloseTenantManager 关闭多租户管理器，释放所有连接池资源
func CloseTenantManager() {
	if tenantMgr == nil {
		return
	}
	tenantMgr.mu.Lock()
	defer tenantMgr.mu.Unlock()
	if tenantMgr.sharedSysDB != nil {
		tenantMgr.sharedSysDB.Close()
	}
	for _, sysDB := range tenantMgr.indepSysDBs {
		sysDB.Close()
	}
	tenantMgr.indepSysDBs = nil
	tenantMgr.indepTpls = nil
	tenantMgr = nil
}

// defaultConfigLoader 内存配置加载器
type defaultConfigLoader struct {
	configs map[string]*TenantConfig
	mu      sync.RWMutex
}

func (l *defaultConfigLoader) Load(tenantID string) (*TenantConfig, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	cfg, ok := l.configs[tenantID]
	if !ok {
		return nil, fmt.Errorf("租户 %s 不存在", tenantID)
	}
	return cfg, nil
}

func (l *defaultConfigLoader) set(tenantID string, cfg *TenantConfig) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.configs[tenantID] = cfg
}

// getOrCreateIndepDB 获取或创建独立 schema 的底层数据库
func (m *tenantManager) getOrCreateIndepDB(cfg *TenantConfig) (*SysDB, tpl.ITPLContext, error) {
	m.mu.RLock()
	db, ok := m.indepSysDBs[cfg.TenantID]
	m.mu.RUnlock()
	if ok {
		return db, m.indepTpls[cfg.TenantID], nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if db, ok := m.indepSysDBs[cfg.TenantID]; ok {
		return db, m.indepTpls[cfg.TenantID], nil
	}

	dsn := cfg.DSN
	if dsn == "" {
		if m.defaultDSNTemplate == "" {
			return nil, nil, errors.New("独立模式缺少 DSN 或默认模板")
		}
		dsn = fmt.Sprintf(m.defaultDSNTemplate, cfg.SchemaName)
	}
	sysDB, err := NewSysDB(m.defaultDriver, dsn, 20, 10, 30*time.Minute)
	if err != nil {
		return nil, nil, err
	}
	tplCtx, err := tpl.GetDBContext(m.defaultDriver)
	if err != nil {
		sysDB.Close()
		return nil, nil, err
	}
	m.indepSysDBs[cfg.TenantID] = sysDB
	m.indepTpls[cfg.TenantID] = tplCtx
	return sysDB, tplCtx, nil
}

// getSharedExecutor 从共享连接池获取一个设置了租户会话的连接
func (m *tenantManager) getSharedExecutor(cfg *TenantConfig) (*sessionExecutor, error) {
	conn, err := m.sharedSysDB.db.Conn(context.Background())
	if err != nil {
		return nil, err
	}

	var setSQL string
	var setArgs []interface{}
	switch cfg.Model {
	case ModelSharedSchema:
		setSQL = "SET search_path TO " + cfg.SchemaName + ", public"
	case ModelSharedRLS:
		setSQL = "SET app.current_tenant = $1"
		setArgs = []interface{}{cfg.TenantID}
	default:
		conn.Close()
		return nil, fmt.Errorf("不支持的模式: %s", cfg.Model)
	}

	_, err = conn.ExecContext(context.Background(), setSQL, setArgs...)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return &sessionExecutor{
		conn: conn,
		tpl:  m.sharedTpl,
	}, nil
}

// ========================= 共享会话执行器 =========================

// sessionExecutor 封装一个已设置租户会话的 *sql.Conn，实现 IBaseDB
type sessionExecutor struct {
	conn *sql.Conn
	tpl  tpl.ITPLContext
}

func (r *sessionExecutor) Close() error {
	return r.conn.Close()
}

func (r *sessionExecutor) Execute(query string, args ...interface{}) (int64, error) {
	result, err := r.conn.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *sessionExecutor) Executes(query string, args ...interface{}) (int64, int64, error) {
	result, err := r.conn.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	lastID, _ := result.LastInsertId()
	return lastID, affected, nil
}

func (r *sessionExecutor) Query(query string, args ...interface{}) (QueryRows, error) {
	rows, err := r.conn.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return resolveFullRows(rows)
}

func (r *sessionExecutor) FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	return r.conn.QueryContext(context.Background(), query, args...)
}

func (r *sessionExecutor) Begin(cfg *TenantConfig) (ISysDBTrans, error) {
	tx, err := r.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	var setSQL string
	var setArgs []interface{}
	switch cfg.Model {
	case ModelSharedSchema:
		setSQL = "SET search_path TO " + cfg.SchemaName + ", public"
	case ModelSharedRLS:
		setSQL = "SET app.current_tenant = $1"
		setArgs = []interface{}{cfg.TenantID}
	}
	_, err = tx.ExecContext(context.Background(), setSQL, setArgs...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &sessionTransaction{tx: tx, tpl: r.tpl}, nil
}

// sessionTransaction 事务对象
type sessionTransaction struct {
	tx  *sql.Tx
	tpl tpl.ITPLContext
}

func (t *sessionTransaction) Execute(query string, args ...interface{}) (int64, error) {
	result, err := t.tx.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (t *sessionTransaction) Executes(query string, args ...interface{}) (int64, int64, error) {
	result, err := t.tx.ExecContext(context.Background(), query, args...)
	if err != nil {
		return 0, 0, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, 0, err
	}
	lastID, _ := result.LastInsertId()
	return lastID, affected, nil
}

func (t *sessionTransaction) Query(query string, args ...interface{}) (QueryRows, error) {
	rows, err := t.tx.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return resolveFullRows(rows)
}

func (t *sessionTransaction) FetchRows(query string, args ...interface{}) (*sql.Rows, error) {
	return t.tx.QueryContext(context.Background(), query, args...)
}

func (t *sessionTransaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *sessionTransaction) Commit() error {
	return t.tx.Commit()
}

// ========================= 租户数据库对象（实现 IDB） =========================

// TenantDB 实现了 IDB 接口
type TenantDB struct {
	tenantID   string
	cfg        *TenantConfig
	indepSysDB *SysDB
	indepTpl   tpl.ITPLContext
}

// NewTenantDB 创建一个绑定租户的数据库操作对象
func NewTenantDB(tenantID string) (IDB, error) {
	if tenantMgr == nil {
		return nil, errors.New("tenant manager not initialized, call InitTenantManager first")
	}
	cfg, err := tenantMgr.configLoader.Load(tenantID)
	if err != nil {
		return nil, err
	}
	td := &TenantDB{tenantID: tenantID, cfg: cfg}
	if cfg.Model == ModelIndependentSchema {
		sysDB, tplCtx, err := tenantMgr.getOrCreateIndepDB(cfg)
		if err != nil {
			return nil, err
		}
		td.indepSysDB = sysDB
		td.indepTpl = tplCtx
	}
	return td, nil
}

func (td *TenantDB) isSharedMode() bool {
	return td.cfg.Model == ModelSharedSchema || td.cfg.Model == ModelSharedRLS
}

func (td *TenantDB) getExecutor() (*sessionExecutor, error) {
	return tenantMgr.getSharedExecutor(td.cfg)
}

// ========== IDB 接口实现 ==========

func (td *TenantDB) Query(sql string, input map[string]interface{}) (QueryRows, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return nil, err
		}
		defer exec.Close()
		query, args := exec.tpl.GetSQLContext(sql, input)
		return exec.Query(query, args...)
	}
	return query(td.indepSysDB, td.indepTpl, sql, input)
}

func (td *TenantDB) Execute(sql string, input map[string]interface{}) (int64, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return 0, err
		}
		defer exec.Close()
		query, args := exec.tpl.GetSQLContext(sql, input)
		return exec.Execute(query, args...)
	}
	return execute(td.indepSysDB, td.indepTpl, sql, input)
}

func (td *TenantDB) Executes(sql string, input map[string]interface{}) (int64, int64, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return 0, 0, err
		}
		defer exec.Close()
		query, args := exec.tpl.GetSQLContext(sql, input)
		return exec.Executes(query, args...)
	}
	return executes(td.indepSysDB, td.indepTpl, sql, input)
}

func (td *TenantDB) Scalar(sql string, input map[string]interface{}) (interface{}, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return nil, err
		}
		defer exec.Close()
		query, args := exec.tpl.GetSQLContext(sql, input)
		rows, err := exec.Query(query, args...)
		if err != nil {
			return nil, err
		}
		if rows.Len() == 0 || rows.Get(0).IsEmpty() {
			return nil, nil
		}
		data, _ := rows.Get(0).Get(rows.Get(0).Keys()[0])
		return data, nil
	}
	return scalar(td.indepSysDB, td.indepTpl, sql, input)
}

func (td *TenantDB) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(td, sqls, input)
}

func (td *TenantDB) QueryBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return queryBatch(td, sqls, input)
}

func (td *TenantDB) InsertBatch(sql string, inputs []map[string]interface{}) (int64, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return 0, err
		}
		defer exec.Close()
		return insertBatch(exec, exec.tpl, sql, inputs)
	}
	return insertBatch(td.indepSysDB, td.indepTpl, sql, inputs)
}

func (td *TenantDB) UpdateBatch(sql string, inputs []map[string]interface{}) (row int64, err error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return 0, err
		}
		defer exec.Close()
		tx, txErr := exec.Begin(td.cfg)
		if txErr != nil {
			return 0, txErr
		}
		defer func() {
			if err == nil {
				err = tx.Commit()
				return
			}
			tx.Rollback()
		}()
		return updateSave(tx, exec.tpl, sql, inputs)
	}
	tx, err := td.indepSysDB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
			return
		}
		tx.Rollback()
	}()
	return updateSave(tx, td.indepTpl, sql, inputs)
}

func (td *TenantDB) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (int64, error) {
	if td.isSharedMode() {
		return 0, errors.New("存储过程在共享模式下暂不支持")
	}
	return executeSP(td.indepSysDB, td.indepTpl, procName, input, output...)
}

func (td *TenantDB) Begin() (IDBTrans, error) {
	if td.isSharedMode() {
		exec, err := td.getExecutor()
		if err != nil {
			return nil, err
		}
		tx, err := exec.Begin(td.cfg)
		if err != nil {
			exec.Close()
			return nil, err
		}
		return &tenantTransactionWrapper{
			trans: tx.(*sessionTransaction),
			exec:  exec,
			tpl:   exec.tpl,
		}, nil
	}
	sysTx, err := td.indepSysDB.Begin()
	if err != nil {
		return nil, err
	}
	return &DBTrans{
		tpl: td.indepTpl,
		tx:  sysTx,
	}, nil
}

func (td *TenantDB) Close() {}

func (td *TenantDB) FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error) {
	if td.isSharedMode() {
		return nil, errors.New("FetchRows 在共享模式下不支持，请使用 Query 代替")
	}
	return fetchRows(td.indepSysDB, td.indepTpl, sql, input)
}

// ========================= 事务包装器（共享模式） =========================

type tenantTransactionWrapper struct {
	trans *sessionTransaction
	exec  *sessionExecutor
	tpl   tpl.ITPLContext
}

func (w *tenantTransactionWrapper) Query(sql string, input map[string]interface{}) (QueryRows, error) {
	query, args := w.tpl.GetSQLContext(sql, input)
	return w.trans.Query(query, args...)
}

func (w *tenantTransactionWrapper) Execute(sql string, input map[string]interface{}) (int64, error) {
	query, args := w.tpl.GetSQLContext(sql, input)
	return w.trans.Execute(query, args...)
}

func (w *tenantTransactionWrapper) Executes(sql string, input map[string]interface{}) (int64, int64, error) {
	query, args := w.tpl.GetSQLContext(sql, input)
	return w.trans.Executes(query, args...)
}

func (w *tenantTransactionWrapper) Scalar(sql string, input map[string]interface{}) (interface{}, error) {
	rows, err := w.Query(sql, input)
	if err != nil {
		return nil, err
	}
	if rows.Len() == 0 || rows.Get(0).IsEmpty() {
		return nil, nil
	}
	data, _ := rows.Get(0).Get(rows.Get(0).Keys()[0])
	return data, nil
}

func (w *tenantTransactionWrapper) ExecuteBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return executeBatch(w, sqls, input)
}

func (w *tenantTransactionWrapper) QueryBatch(sqls []string, input map[string]interface{}) (QueryRows, error) {
	return queryBatch(w, sqls, input)
}

func (w *tenantTransactionWrapper) InsertBatch(sql string, inputs []map[string]interface{}) (int64, error) {
	return insertBatch(w.trans, w.tpl, sql, inputs)
}

func (w *tenantTransactionWrapper) UpdateBatch(sql string, inputs []map[string]interface{}) (int64, error) {
	return updateSave(w.trans, w.tpl, sql, inputs)
}

func (w *tenantTransactionWrapper) ExecuteSP(procName string, input map[string]interface{}, output ...interface{}) (int64, error) {
	return executeSP(w.trans, w.tpl, procName, input, output...)
}

func (w *tenantTransactionWrapper) FetchRows(sql string, input map[string]interface{}) (*sql.Rows, error) {
	query, args := w.tpl.GetSQLContext(sql, input)
	return w.trans.FetchRows(query, args...)
}

func (w *tenantTransactionWrapper) Rollback() error {
	err := w.trans.Rollback()
	w.exec.Close()
	return err
}

func (w *tenantTransactionWrapper) Commit() error {
	err := w.trans.Commit()
	w.exec.Close()
	return err
}
