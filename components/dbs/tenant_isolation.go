package dbs

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/micro-plat/hydra/conf/app"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
	libdb "github.com/micro-plat/lib4go/db"
)

// 本文件实现「按租户隔离模式自动路由 DB」的通用能力。
//   - ROW          → 主库（行隔离，由业务 SQL 的 tenant_id WHERE 过滤）
//   - SCHEMA       → 同实例独立库（DSN 模板生成 dbname）
//   - DEDICATED    → 独立实例（db_dsn 列显式指定）
// SCHEMA/DEDICATED 均复用 lib4go 的 indep_schema 通路（NewTenantDB）。
//
// 设计要点：
//   - TenantIsolationConfig 字段私有，包外无法用 struct literal 初始化，必须走 Builder。
//   - 主库 driver/DSN/池参数自动从 var 配置树读取（app.Cache.GetVarConf），无需业务显式注入。
//   - 路由结果用包级 sync.Map 缓存（进程内，无 TTL）；RefreshTenantIsolation 手动清。
//   - lib4go 不改，仅消费其 InitTenantManager/SetTenantConfigLoader/NewTenantDB。

// ========================= 配置（字段私有，仅 Builder 可构造）=========================

// TenantIsolationConfig 多租户混合隔离配置
// 字段私有：包外无法 { } 初始化，必须通过 NewTenantIsolationConfig(...).WithXxx(...).Build() 构造
type TenantIsolationConfig struct {
	mainDBName         string // 主库 var 名，默认 "db"
	querySQL           string // 必填：查询 SQL，含 @tenant_id，返回 modeField 与 dsnField 两列
	modeField          string // 模式列名，默认 "isolation_mode"
	dsnField           string // DSN 列名，默认 "db_dsn"（DEDICATED 用）
	schemaTemplate     string // SCHEMA 模式 DSN 模板，含一个 %s；为空则 SCHEMA 模式不可用
	schemaPrefix       string // SCHEMA 库名前缀（最终库名 = prefix + tenantID）
	rowModeValue       string // 行隔离标识值，默认 "ROW"
	schemaModeValue    string // 独立库-同实例标识值，默认 "SCHEMA"
	dedicatedModeValue string // 独立库-独立实例标识值，默认 "DEDICATED"
	maxOpen            int    // 独立池最大连接，默认从主库继承或 20
	maxIdle            int    // 默认从主库继承或 10
	maxLifeSec         int    // 默认从主库继承或 1800
	sharedDriver       string // 可选：覆盖主库 driver（一般无需，自动从 var 读）
	sharedDSN          string // 可选：覆盖主库 DSN（一般无需，自动从 var 读）
}

// ========================= Builder（必填进 New，可选进 With）=========================

// TenantIsolationConfigBuilder 配置构造器
type TenantIsolationConfigBuilder struct {
	cfg TenantIsolationConfig
}

// NewTenantIsolationConfig 创建配置构造器
// querySQL 是唯一必填项（含表名与条件，业务特定，无默认值）；其余参数均有默认值，按需用 WithXxx 覆盖
func NewTenantIsolationConfig(querySQL string) *TenantIsolationConfigBuilder {
	return &TenantIsolationConfigBuilder{cfg: TenantIsolationConfig{
		querySQL:           querySQL,
		mainDBName:         "db",
		modeField:          "isolation_mode",
		dsnField:           "db_dsn",
		rowModeValue:       "ROW",
		schemaModeValue:    "SCHEMA",
		dedicatedModeValue: "DEDICATED",
		maxOpen:            0, // 0 表示从主库继承
		maxIdle:            0,
		maxLifeSec:         0,
	}}
}

// WithMainDB 设置主库（共享库）var 名，默认 "db"
func (b *TenantIsolationConfigBuilder) WithMainDB(name string) *TenantIsolationConfigBuilder {
	b.cfg.mainDBName = name
	return b
}

// WithFields 设置模式列名与 DSN 列名，默认 "isolation_mode"/"db_dsn"
// 当 querySQL 返回列名与默认不同时使用（或在 SQL 内用 AS 对齐默认列名）
func (b *TenantIsolationConfigBuilder) WithFields(modeField, dsnField string) *TenantIsolationConfigBuilder {
	b.cfg.modeField = modeField
	b.cfg.dsnField = dsnField
	return b
}

// WithModeValues 设置三种模式的标识值，默认 "ROW"/"SCHEMA"/"DEDICATED"
func (b *TenantIsolationConfigBuilder) WithModeValues(row, schema, dedicated string) *TenantIsolationConfigBuilder {
	b.cfg.rowModeValue = row
	b.cfg.schemaModeValue = schema
	b.cfg.dedicatedModeValue = dedicated
	return b
}

// WithSchema 设置 SCHEMA 模式参数（配置后 SCHEMA 模式才可用）
//   - template: DSN 模板，含一个 %s（由 SchemaName 填充，lib4go tenant.go:182）
//   - prefix:   库名前缀（最终库名 = prefix + tenant_code，如 ismart_yfyl；tenant_code 缺失时回退 tenantID）
func (b *TenantIsolationConfigBuilder) WithSchema(template, prefix string) *TenantIsolationConfigBuilder {
	b.cfg.schemaTemplate = template
	b.cfg.schemaPrefix = prefix
	return b
}

// WithPool 设置独立库连接池参数；默认从主库 var 配置继承
func (b *TenantIsolationConfigBuilder) WithPool(maxOpen, maxIdle, maxLifeSec int) *TenantIsolationConfigBuilder {
	b.cfg.maxOpen = maxOpen
	b.cfg.maxIdle = maxIdle
	b.cfg.maxLifeSec = maxLifeSec
	return b
}

// WithSharedDSN 显式覆盖主库 driver/DSN
// 可选：默认自动从 var 配置读取；仅当主库未通过标准 DB var 注册时才需要
func (b *TenantIsolationConfigBuilder) WithSharedDSN(driver, dsn string) *TenantIsolationConfigBuilder {
	b.cfg.sharedDriver = driver
	b.cfg.sharedDSN = dsn
	return b
}

// Build 生成配置
// querySQL 为空时 panic（New 已要求传入）；SCHEMA 模板运行时按需校验（见 lib4goTenantLoader）
func (b *TenantIsolationConfigBuilder) Build() TenantIsolationConfig {
	if b.cfg.querySQL == "" {
		panic("TenantIsolationConfig: querySQL 为必填项（NewTenantIsolationConfig 参数）")
	}
	return b.cfg
}

// ========================= Provider 接口（自定义逃生通道）=========================

// TenantIsolationProvider 自定义隔离配置来源
// mode 对照 rowMode/schemaMode/dedicatedMode；dsn 仅供 DEDICATED 用
// code 为租户编码（tenant_code），SCHEMA 模式用于派生库名（prefix+code）；
// 配置来源无该列时返回空串，SCHEMA 库名回退 prefix+tenantID（兼容旧 querySQL）
// 不注册自定义 provider 时，使用内置 DB provider（查主库 + sync.Map 缓存）
type TenantIsolationProvider interface {
	Load(tenantID string) (mode string, dsn string, code string, err error)
}

// ========================= 全局状态 =========================

var (
	isolationCfg      TenantIsolationConfig
	isolationProvider TenantIsolationProvider
	sharedStandardDB  *StandardDB  // 由 NewStandardDB 设置，供内置 provider 查主库
	tenantMgrOnce     sync.Once    // ensureTenantManager 只初始化一次
	tenantModeCache   sync.Map     // tenantID -> tenantIsolationValue（进程内缓存，无 TTL）
)

// SetTenantIsolationConfig 注册配置（参数须由 Builder.Build() 生成）
func SetTenantIsolationConfig(cfg TenantIsolationConfig) {
	isolationCfg = cfg
	if isolationProvider == nil { // 未先注册自定义 provider 时启用内置
		isolationProvider = &dbTenantIsolationProvider{}
	}
}

// RegisterTenantIsolationProvider 注册自定义 provider（优先于内置 DB provider）
// 适合配置来源非主库的场景（如配置中心、外部服务）
func RegisterTenantIsolationProvider(p TenantIsolationProvider) {
	isolationProvider = p
}

// RefreshTenantIsolation 清除某租户的路由缓存（管理端改模式后调用）
// 注意：lib4go 的独立库连接池（indepSysDBs）无按租户释放 API，已建连接保留到进程重启；
// 此处仅清路由判断缓存，使下次请求重新查主库并按新模式路由。
func RefreshTenantIsolation(tenantID int64) {
	tenantModeCache.Delete(strconv.FormatInt(tenantID, 10))
}

// ========================= 内置 DB provider =========================

// dbTenantIsolationProvider 内置实现：查主库 + sync.Map 缓存
type dbTenantIsolationProvider struct{}

func (p *dbTenantIsolationProvider) Load(tenantID string) (string, string, string, error) {
	// 1. 缓存命中
	if v, ok := tenantModeCache.Load(tenantID); ok {
		iv := v.(tenantIsolationValue)
		return iv.mode, iv.dsn, iv.code, nil
	}
	// 2. 经 sharedStandardDB 查主库（防递归：GetDB 走共享库，不可走 NewTenantDB）
	if sharedStandardDB == nil {
		return "", "", "", fmt.Errorf("tenant isolation: StandardDB 未初始化")
	}
	db, err := sharedStandardDB.GetDB(isolationCfg.mainDBName)
	if err != nil {
		return "", "", "", fmt.Errorf("获取主库 %s 失败: %w", isolationCfg.mainDBName, err)
	}
	rows, err := db.Query(isolationCfg.querySQL, map[string]interface{}{"tenant_id": tenantID})
	if err != nil {
		return "", "", "", fmt.Errorf("query tenant isolation: %w", err)
	}
	if rows.IsEmpty() {
		return isolationCfg.rowModeValue, "", "", nil // 租户不存在 → 兜底 ROW，避免请求中断
	}
	mode := rows.Get(0).GetString(isolationCfg.modeField)
	if mode == "" {
		mode = isolationCfg.rowModeValue
	}
	dsn := rows.Get(0).GetString(isolationCfg.dsnField)
	// querySQL 未 SELECT tenant_code 列时 GetString 返回空串（xmap 缺键安全），SCHEMA 库名回退 tenantID
	code := rows.Get(0).GetString("tenant_code")
	tenantModeCache.Store(tenantID, tenantIsolationValue{mode: mode, dsn: dsn, code: code})
	return mode, dsn, code, nil
}

type tenantIsolationValue struct {
	mode string
	dsn  string
	code string // tenant_code，SCHEMA 模式派生库名（prefix+code）用
}

// ========================= lib4go 桥接 =========================

// lib4goTenantLoader 把 hydra provider 适配为 lib4go.TenantConfigLoader
// 由 libdb.NewTenantDB(tenantID) 在创建独立库时回调（lib4go tenant.go:360）
type lib4goTenantLoader struct{}

func (l *lib4goTenantLoader) Load(tenantID string) (*libdb.TenantConfig, error) {
	// 从 sync.Map 取（GetTenantDBByMode 已 Store）；miss 则回退查 provider
	var mode, dsn, code string
	if v, ok := tenantModeCache.Load(tenantID); ok {
		iv := v.(tenantIsolationValue)
		mode, dsn, code = iv.mode, iv.dsn, iv.code
	} else if isolationProvider != nil {
		m, d, c, err := isolationProvider.Load(tenantID)
		if err != nil {
			return nil, err
		}
		mode, dsn, code = m, d, c
	}
	cfg := &libdb.TenantConfig{TenantID: tenantID}
	switch mode {
	case isolationCfg.schemaModeValue:
		if isolationCfg.schemaTemplate == "" {
			return nil, fmt.Errorf("tenant %s 为 SCHEMA 模式但未配置模板，请用 WithSchema 设置", tenantID)
		}
		// 同实例独立库：DSN 留空，lib4go 用 SchemaName 经模板生成（tenant.go:182）
		// 库名 = prefix + tenant_code（如 ismart_yfyl）；code 为空回退 tenantID（兼容未返回该列的旧 querySQL）
		suffix := tenantID
		if code != "" {
			suffix = strings.ToLower(code)
		}
		cfg.Model = libdb.ModelIndependentSchema
		cfg.SchemaName = isolationCfg.schemaPrefix + suffix
		cfg.DSN = ""
	case isolationCfg.dedicatedModeValue:
		// 独立实例：用 db_dsn 列的显式连接串（tenant.go:177 走 cfg.DSN 分支）
		cfg.Model = libdb.ModelIndependentSchema
		cfg.DSN = dsn
	default: // ROW / 未知
		cfg.Model = libdb.ModelSharedSchema
	}
	return cfg, nil
}

// ensureTenantManager 首次库隔离路由时懒初始化 lib4go TenantManager + 注册 loader
// sync.Once 保证全局一次；InitTenantManager 内部也自带 nil 守卫（lib4go tenant.go:75-77）
func ensureTenantManager() error {
	var initErr error
	tenantMgrOnce.Do(func() {
		driver, dsn := isolationCfg.sharedDriver, isolationCfg.sharedDSN
		maxOpen, maxIdle, maxLife := isolationCfg.maxOpen, isolationCfg.maxIdle, isolationCfg.maxLifeSec
		// 未显式覆盖时从 var 配置树读主库
		if driver == "" || dsn == "" || maxOpen == 0 || maxIdle == 0 || maxLife == 0 {
			d, dn, mo, mi, ml, e := readMainDBConf(isolationCfg.mainDBName)
			if e != nil {
				initErr = fmt.Errorf("tenant isolation: 读主库 %s 配置失败: %w（可用 WithSharedDSN 显式注入）", isolationCfg.mainDBName, e)
				return
			}
			if driver == "" {
				driver = d
			}
			if dsn == "" {
				dsn = dn
			}
			if maxOpen == 0 {
				maxOpen = mo
			}
			if maxIdle == 0 {
				maxIdle = mi
			}
			if maxLife == 0 {
				maxLife = ml
			}
		}
		// 兜底默认
		if maxOpen == 0 {
			maxOpen = 20
		}
		if maxIdle == 0 {
			maxIdle = 10
		}
		if maxLife == 0 {
			maxLife = 1800
		}
		if driver == "" || dsn == "" {
			initErr = fmt.Errorf("tenant isolation: 主库 driver/DSN 不可用")
			return
		}
		// 第 6 参数 defaultDriver 固定 mysql；第 7 参数 = SCHEMA 模板（含一个 %s）
		if err := libdb.InitTenantManager(driver, dsn, maxOpen, maxIdle, maxLife, "mysql", isolationCfg.schemaTemplate); err != nil {
			initErr = err
			return
		}
		libdb.SetTenantConfigLoader(&lib4goTenantLoader{})
	})
	return initErr
}

// readMainDBConf 从 hydra var 配置树读主库的 driver/DSN/池参数
// 复用 db.go GetTenantDB 的读取模式（app.Cache.GetVarConf）
func readMainDBConf(name string) (driver, dsn string, maxOpen, maxIdle, maxLife int, err error) {
	varConf, e := app.Cache.GetVarConf()
	if e != nil {
		return "", "", 0, 0, 0, e
	}
	if !varConf.Has(dbTypeNode, name) {
		return "", "", 0, 0, 0, fmt.Errorf("节点/%s/%s 未配置", dbTypeNode, name)
	}
	jconf, e := varConf.GetConf(dbTypeNode, name)
	if e != nil {
		return "", "", 0, 0, 0, e
	}
	var c xdb.DB
	if e := jconf.ToStruct(&c); e != nil {
		return "", "", 0, 0, 0, e
	}
	return c.Provider, c.ConnString, c.MaxOpen, c.MaxIdle, c.LifeTime, nil
}
