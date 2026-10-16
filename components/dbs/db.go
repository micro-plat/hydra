package dbs

import (
	"fmt"
	"strconv"

	"github.com/micro-plat/hydra/components/container"
	libdb "github.com/micro-plat/lib4go/db"
	"github.com/micro-plat/lib4go/types"

	"github.com/micro-plat/hydra/conf"
	"github.com/micro-plat/hydra/conf/app"
	xdb "github.com/micro-plat/hydra/conf/vars/db"
)

const (
	//typeNode DB在var配置中的类型名称
	dbTypeNode = "db"

	//nameNode DB名称在var配置中的末节点名称
	dbNameNode = "db"
)

//StandardDB db
type StandardDB struct {
	c container.IContainer
}

//NewStandardDB 创建DB
func NewStandardDB(c container.IContainer) *StandardDB {
	s := &StandardDB{c: c}
	sharedStandardDB = s // 供内置 tenant isolation provider 查主库（tenant_isolation.go）
	return s
}

//GetRegularDB 获取正式的没有异常数据库实例
func (s *StandardDB) GetRegularDB(names ...string) (d IDB) {
	d, err := s.GetDB(names...)
	if err != nil {
		panic(err)
	}
	return d
}

//GetDB 获取数据库操作对象
func (s *StandardDB) GetDB(names ...string) (d IDB, err error) {
	name := types.GetStringByIndex(names, 0, dbNameNode)
	obj, err := s.c.GetOrCreate(dbTypeNode, name, func(conf *conf.RawConf, keys ...string) (obj interface{}, err error) {
		if conf.IsEmpty() {
			return nil, fmt.Errorf("节点/%s/%s未配置，或不可用", dbTypeNode, name)
		}
		var dbConf xdb.DB
		if err = conf.ToStruct(&dbConf); err != nil {
			return nil, fmt.Errorf("数据库[%s/%s]配置有误：%w", dbTypeNode, name, err)
		}
		return libdb.NewDB(dbConf.Provider, dbConf.ConnString, dbConf.MaxOpen, dbConf.MaxIdle, dbConf.LifeTime)
	})
	if err != nil {
		return nil, err
	}
	return obj.(IDB), nil
}

//GetRegularTenantDB 获取绑定租户的数据库实例，出错时 panic
func (s *StandardDB) GetRegularTenantDB(tenantID string, names ...string) (d IDB) {
	d, err := s.GetTenantDB(tenantID, names...)
	if err != nil {
		panic(err)
	}
	return d
}

//GetTenantDB 获取绑定租户的数据库实例
// tenantID 为空或配置未启用多租户时，自动降级到普通 GetDB
func (s *StandardDB) GetTenantDB(tenantID string, names ...string) (d IDB, err error) {
	if tenantID == "" {
		return s.GetDB(names...)
	}

	name := types.GetStringByIndex(names, 0, dbNameNode)

	varConf, err := app.Cache.GetVarConf()
	if err != nil {
		return nil, fmt.Errorf("无法获取var配置: %w", err)
	}
	if !varConf.Has(dbTypeNode, name) {
		return nil, fmt.Errorf("节点/%s/%s未配置", dbTypeNode, name)
	}
	jconf, err := varConf.GetConf(dbTypeNode, name)
	if err != nil {
		return nil, err
	}
	var dbConf xdb.DB
	if err := jconf.ToStruct(&dbConf); err != nil {
		return nil, fmt.Errorf("数据库[%s/%s]配置有误: %w", dbTypeNode, name, err)
	}

	if dbConf.TenantMode == "" {
		return s.GetDB(names...)
	}

	if err := libdb.InitTenantManager(
		dbConf.Provider, dbConf.ConnString,
		dbConf.MaxOpen, dbConf.MaxIdle, dbConf.LifeTime,
		"", "",
	); err != nil {
		return nil, fmt.Errorf("初始化多租户管理器失败: %w", err)
	}

	libdb.SetTenantConfig(tenantID, &libdb.TenantConfig{
		TenantID:   tenantID,
		Model:      libdb.TenantModel(dbConf.TenantMode),
		SchemaName: tenantID,
	})

	return libdb.NewTenantDB(tenantID)
}

//GetTenantDBByMode 按租户隔离模式路由数据库实例
//   - 未配置隔离（querySQL 空）、tenantID<=0、模式为 ROW/未知 → 返回共享库 GetDB（行隔离，由业务 SQL 的 tenant_id WHERE 保证）
//   - 模式为 SCHEMA/DEDICATED → 返回独立库 libdb.NewTenantDB（indep_schema 通路）
//
// tenantID 为 int64（ismart 业务侧普遍类型），内部转 string 供缓存键与 lib4go 使用，调用方无需 strconv
// names 仅对共享库路径生效（同 GetDB）；独立库路径固定用 mainDBName 的配置派生连接池
func (s *StandardDB) GetTenantDBByMode(tenantID int64, names ...string) (d IDB, err error) {
	tid := strconv.FormatInt(tenantID, 10)
	// 未配置隔离或无 tenantID → 共享库（灰度安全：未接入的 service 行为不变）
	if isolationCfg.querySQL == "" || tenantID <= 0 {
		return s.GetDB(names...)
	}
	mode, dsn, code, err := isolationProvider.Load(tid)
	if err != nil {
		return nil, fmt.Errorf("加载租户 %s 隔离配置失败: %w", tid, err)
	}
	// ROW / 未知值 → 共享库（未知值安全降级，避免误走 SET search_path）
	if mode != isolationCfg.schemaModeValue && mode != isolationCfg.dedicatedModeValue {
		return s.GetDB(names...)
	}
	// SCHEMA / DEDICATED → 独立库
	if err := ensureTenantManager(); err != nil {
		return nil, fmt.Errorf("初始化租户管理器失败: %w", err)
	}
	// 预存路由结果，供 lib4goTenantLoader 在 NewTenantDB 内部读取（自定义 provider 亦覆盖，避免二次查询）
	tenantModeCache.Store(tid, tenantIsolationValue{mode: mode, dsn: dsn, code: code})
	return libdb.NewTenantDB(tid)
}

//GetRegularTenantDBByMode 按租户隔离模式路由数据库实例，出错时 panic（与 GetRegularDB 风格一致）
func (s *StandardDB) GetRegularTenantDBByMode(tenantID int64, names ...string) (d IDB) {
	d, err := s.GetTenantDBByMode(tenantID, names...)
	if err != nil {
		panic(err)
	}
	return d
}
