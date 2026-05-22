package db_test

import (
	"fmt"
	"log"

	"github.com/micro-plat/lib4go/db"
)

// ======================== 推荐模式：共享连接池 + SET search_path ========================

// Example_tenantSharedSchema 演示推荐的多租户模式：共享连接池 + SET search_path
//
// 适用场景：同一个 PostgreSQL 实例，不同 schema 隔离租户数据
// 核心优势：
//   - 无论多少租户，只有 1 个连接池
//   - 无需修改任何表结构，无需加 tenant_id 列
//   - public schema 中的公共表（字典、配置等）自动可见
//   - SQL 语句完全不用改
//
// 前置条件（PostgreSQL 中执行一次）：
//
//	-- 创建租户 schema
//	CREATE SCHEMA tenant_a;
//	CREATE SCHEMA tenant_b;
//
//	-- 在各 schema 中创建业务表
//	CREATE TABLE tenant_a.users (id serial, name text, email text);
//	CREATE TABLE tenant_b.users (id serial, name text, email text);
//
//	-- 公共表放在 public schema，所有租户自动可见
//	CREATE TABLE public.sys_dict (id serial, key text, value text);
func Example_tenantSharedSchema() {
	// 1. 初始化 — 只需一个连接池
	err := db.InitTenantManager(
		"postgres",                                            // 驱动
		"postgres://user:pass@localhost/myapp?sslmode=disable", // 连接串（指向同一个数据库）
		50, 10, 60,                                            // 连接池参数
		"", "",                                                // 独立模式不用，传空
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseTenantManager()

	// 2. 配置租户（SchemaName = PostgreSQL 中的 schema 名）
	db.SetTenantConfig("tenant_a", &db.TenantConfig{
		TenantID:   "tenant_a",
		Model:      db.ModelSharedSchema,
		SchemaName: "tenant_a",
	})
	db.SetTenantConfig("tenant_b", &db.TenantConfig{
		TenantID:   "tenant_b",
		Model:      db.ModelSharedSchema,
		SchemaName: "tenant_b",
	})

	// 3. 业务中使用 — 根据登录用户的租户信息获取数据库
	// 实际场景中 tenantID 从用户 session/token 中获取
	tenantID := "tenant_a" // 来自登录用户信息
	tenantDB, err := db.NewTenantDB(tenantID)
	if err != nil {
		log.Fatal(err)
	}
	defer tenantDB.Close()

	// 4. 所有查询自动在 tenant_a schema 中执行
	// SQL 无需任何修改，也不用加 WHERE tenant_id = xxx
	rows, err := tenantDB.Query(
		"SELECT id, name FROM users WHERE name LIKE @name",
		map[string]interface{}{"name": "%张%"},
	)
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < rows.Len(); i++ {
		row := rows.Get(i)
		id, _ := row.Get("id")
		name, _ := row.Get("name")
		fmt.Printf("ID: %v, Name: %v\n", id, name)
	}

	// 5. 插入 — 直接写 SQL，自动路由到 tenant_a.users
	affected, err := tenantDB.Execute(
		"INSERT INTO users (name, email) VALUES (@name, @email)",
		map[string]interface{}{
			"name":  "张三",
			"email": "zhangsan@example.com",
		},
	)
	fmt.Printf("插入 %d 行\n", affected)

	// 6. 查询 public schema 的公共表 — 自动可见，无需加 public. 前缀
	dictRows, err := tenantDB.Query("SELECT key, value FROM sys_dict", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("字典条目: %d\n", dictRows.Len())

	// 7. 事务 — 转账示例
	tx, err := tenantDB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.Execute(
		"UPDATE accounts SET balance = balance - @amount WHERE id = @id",
		map[string]interface{}{"id": 1, "amount": 100},
	)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	_, err = tx.Execute(
		"UPDATE accounts SET balance = balance + @amount WHERE id = @id",
		map[string]interface{}{"id": 2, "amount": 100},
	)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	tx.Commit()

	// 8. 切换租户 — 完全独立的数据空间
	tenantB, _ := db.NewTenantDB("tenant_b")
	defer tenantB.Close()
	rows, _ = tenantB.Query("SELECT id, name FROM users", nil)
	fmt.Printf("tenant_b 用户数: %d\n", rows.Len())
}

// ======================== 模式二：独立 Schema（独立连接池） ========================

// Example_tenantIndependentSchema 演示独立 schema 模式
// 每个租户创建独立的连接池，适用于租户数量较少（<30）且需要完全物理隔离的场景
func Example_tenantIndependentSchema() {
	err := db.InitTenantManager(
		"postgres",                                            // 共享模式占位
		"postgres://user:pass@localhost/main?sslmode=disable", // 共享模式占位
		10, 5, 60,
		"postgres",                                              // 独立模式驱动
		"postgres://user:pass@localhost/mydb?search_path=%s",    // DSN 模板
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseTenantManager()

	db.SetTenantConfig("tenant_a", &db.TenantConfig{
		TenantID:   "tenant_a",
		Model:      db.ModelIndependentSchema,
		SchemaName: "tenant_a",
	})

	tenantDB, _ := db.NewTenantDB("tenant_a")
	defer tenantDB.Close()

	rows, _ := tenantDB.Query("SELECT id, name FROM users", nil)
	fmt.Printf("用户数: %d\n", rows.Len())
}

// ======================== 模式三：共享表 + RLS ========================

// Example_tenantSharedRLS 演示共享表 + PostgreSQL RLS 模式
// 所有租户共享同一组表，通过 RLS 策略自动过滤数据
//
// 前置条件：
//
//	ALTER TABLE users ADD COLUMN tenant_id TEXT NOT NULL DEFAULT '';
//	ALTER TABLE users ENABLE ROW LEVEL SECURITY;
//	CREATE POLICY tenant_isolation ON users
//	  USING (tenant_id = current_setting('app.current_tenant')::text);
func Example_tenantSharedRLS() {
	err := db.InitTenantManager(
		"postgres", "postgres://user:pass@localhost/myapp?sslmode=disable",
		20, 10, 60, "", "",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseTenantManager()

	db.SetTenantConfig("company_a", &db.TenantConfig{
		TenantID: "company_a",
		Model:    db.ModelSharedRLS,
	})

	tenantDB, _ := db.NewTenantDB("company_a")
	defer tenantDB.Close()

	// 查询自动只返回 company_a 的数据（RLS 策略过滤）
	rows, _ := tenantDB.Query("SELECT id, name FROM users", nil)
	fmt.Printf("company_a 用户: %d\n", rows.Len())

	// 插入需显式传入 tenant_id
	tenantDB.Execute(
		"INSERT INTO users (name, tenant_id) VALUES (@name, @tenant_id)",
		map[string]interface{}{"name": "李四", "tenant_id": "company_a"},
	)
}

// ======================== 批量操作示例 ========================

func Example_tenantBatchOps() {
	_ = db.InitTenantManager(
		"postgres", "postgres://user:pass@localhost/myapp?sslmode=disable",
		10, 5, 60, "", "",
	)
	defer db.CloseTenantManager()

	db.SetTenantConfig("tenant_x", &db.TenantConfig{
		TenantID:   "tenant_x",
		Model:      db.ModelSharedSchema,
		SchemaName: "tenant_x",
	})

	tenantDB, _ := db.NewTenantDB("tenant_x")
	defer tenantDB.Close()

	// 批量插入
	affected, _ := tenantDB.InsertBatch(
		"INSERT INTO products (name, price) VALUES (@name, @price)",
		[]map[string]interface{}{
			{"name": "商品A", "price": 99.9},
			{"name": "商品B", "price": 199.0},
			{"name": "商品C", "price": 299.0},
		},
	)
	fmt.Printf("批量插入 %d 条\n", affected)

	// 批量更新（事务保护，任一失败全部回滚）
	affected, _ = tenantDB.UpdateBatch(
		"UPDATE products SET price = @price WHERE name = @name",
		[]map[string]interface{}{
			{"name": "商品A", "price": 88.0},
			{"name": "商品B", "price": 188.0},
		},
	)
	fmt.Printf("批量更新 %d 条\n", affected)
}

// ======================== 自定义配置加载器 ========================

type myConfigLoader struct{}

func (l *myConfigLoader) Load(tenantID string) (*db.TenantConfig, error) {
	// 实际场景：从数据库、Redis、配置中心动态加载
	configs := map[string]*db.TenantConfig{
		"dynamic_tenant": {
			TenantID:   "dynamic_tenant",
			Model:      db.ModelSharedSchema,
			SchemaName: "schema_dynamic",
		},
	}
	cfg, ok := configs[tenantID]
	if !ok {
		return nil, fmt.Errorf("租户 %s 未找到", tenantID)
	}
	return cfg, nil
}

func Example_tenantCustomLoader() {
	_ = db.InitTenantManager(
		"postgres", "postgres://user:pass@localhost/myapp?sslmode=disable",
		10, 5, 60, "", "",
	)
	db.SetTenantConfigLoader(&myConfigLoader{})

	tenantDB, err := db.NewTenantDB("dynamic_tenant")
	if err != nil {
		log.Fatal(err)
	}
	defer tenantDB.Close()

	rows, _ := tenantDB.Query("SELECT 1", nil)
	fmt.Printf("结果: %d 行\n", rows.Len())
}
