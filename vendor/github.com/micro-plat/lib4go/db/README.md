# db 模块使用手册

## 1. 模块介绍

db模块是lib4go库中的数据库操作模块，提供了统一的数据库操作接口，支持多种数据库类型，并简化了数据库操作的复杂性。

主要特性：
- 支持多种数据库（SQLite3、Oracle、MySQL、PostgreSQL等）
- 提供统一的数据库操作接口
- 支持事务处理
- 支持SQL模板和参数替换
- 支持批量操作
- 提供友好的查询结果处理

## 2. 安装和配置

### 2.1 安装

```bash
go get github.com/micro-plat/lib4go/db
```

### 2.2 导入

```go
import "github.com/micro-plat/lib4go/db"
```

### 2.3 支持的数据库驱动

- SQLite3: `sqlite3`
- Oracle: `oci8`、`ora`
- MySQL: `mysql`
- PostgreSQL: `postgres`
- TDengine: `taosws`、`taossql`、`taosrestful`

## 3. 基本用法

### 3.1 创建数据库连接

使用`NewDB`函数创建数据库连接：

```go
// 创建SQLite3数据库连接
dbInstance, err := db.NewDB("sqlite3", ":memory:", 10, 5, 60)
if err != nil {
    log.Fatalf("创建数据库连接失败: %v", err)
}
defer dbInstance.Close()

// 创建MySQL数据库连接
dbInstance, err := db.NewDB("mysql", "user:password@tcp(localhost:3306)/database", 10, 5, 60)
if err != nil {
    log.Fatalf("创建数据库连接失败: %v", err)
}
defer dbInstance.Close()
```

参数说明：
- `provider`: 数据库类型
- `connString`: 连接字符串
- `maxOpen`: 最大打开连接数
- `maxIdle`: 最大空闲连接数
- `maxLifeTime`: 连接最大生命周期（秒）

### 3.2 执行查询

使用`Query`方法执行查询并获取结果：

```go
// 执行查询
sql := "SELECT * FROM users WHERE id = @id"
input := map[string]interface{}{"id": 1}
rows, err := dbInstance.Query(sql, input)
if err != nil {
    log.Fatalf("查询失败: %v", err)
}

// 处理查询结果
for _, row := range rows {
    id := row.GetInt("id")
    name := row.GetString("name")
    age := row.GetInt("age")
    fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
}
```

### 3.3 执行单条SQL操作

使用`Execute`方法执行插入、更新、删除等操作：

```go
// 执行插入
sql := "INSERT INTO users (name, age) VALUES (@name, @age)"
input := map[string]interface{}{"name": "张三", "age": 30}
rowsAffected, err := dbInstance.Execute(sql, input)
if err != nil {
    log.Fatalf("插入失败: %v", err)
}
fmt.Printf("影响行数: %d\n", rowsAffected)

// 执行更新
sql := "UPDATE users SET age = @age WHERE id = @id"
input := map[string]interface{}{"id": 1, "age": 31}
rowsAffected, err := dbInstance.Execute(sql, input)
if err != nil {
    log.Fatalf("更新失败: %v", err)
}
fmt.Printf("影响行数: %d\n", rowsAffected)

// 执行删除
sql := "DELETE FROM users WHERE id = @id"
input := map[string]interface{}{"id": 1}
rowsAffected, err := dbInstance.Execute(sql, input)
if err != nil {
    log.Fatalf("删除失败: %v", err)
}
fmt.Printf("影响行数: %d\n", rowsAffected)
```

### 3.4 获取单个值

使用`Scalar`方法获取单个值：

```go
sql := "SELECT COUNT(*) FROM users"
count, err := dbInstance.Scalar(sql, nil)
if err != nil {
    log.Fatalf("查询失败: %v", err)
}
fmt.Printf("用户总数: %v\n", count)
```

## 4. 事务处理

使用`Begin`方法创建事务，事务对象支持与数据库实例相同的查询和执行方法：

```go
// 开始事务
trans, err := dbInstance.Begin()
if err != nil {
    log.Fatalf("创建事务失败: %v", err)
}

// 执行事务操作
sql1 := "UPDATE users SET balance = balance - 100 WHERE id = @id1"
sql2 := "UPDATE users SET balance = balance + 100 WHERE id = @id2"

_, err = trans.Execute(sql1, map[string]interface{}{"id1": 1})
if err != nil {
    trans.Rollback()
    log.Fatalf("事务操作失败: %v", err)
}

_, err = trans.Execute(sql2, map[string]interface{}{"id2": 2})
if err != nil {
    trans.Rollback()
    log.Fatalf("事务操作失败: %v", err)
}

// 提交事务
if err := trans.Commit(); err != nil {
    log.Fatalf("提交事务失败: %v", err)
}
```

## 5. 批量操作

### 5.1 批量插入

使用`InsertBatch`方法批量插入数据：

```go
sql := "INSERT INTO users (name, age) VALUES (@name, @age)"
inputs := []map[string]interface{}{
    {"name": "张三", "age": 30},
    {"name": "李四", "age": 28},
    {"name": "王五", "age": 35},
}

rowsAffected, err := dbInstance.InsertBatch(sql, inputs)
if err != nil {
    log.Fatalf("批量插入失败: %v", err)
}
fmt.Printf("插入行数: %d\n", rowsAffected)
```

### 5.2 批量更新

使用`UpdateBatch`方法批量更新数据：

```go
sql := "UPDATE users SET age = @age WHERE id = @id"
inputs := []map[string]interface{}{
    {"id": 1, "age": 31},
    {"id": 2, "age": 29},
    {"id": 3, "age": 36},
}

rowsAffected, err := dbInstance.UpdateBatch(sql, inputs)
if err != nil {
    log.Fatalf("批量更新失败: %v", err)
}
fmt.Printf("更新行数: %d\n", rowsAffected)
```

### 5.3 批量执行SQL语句

使用`ExecuteBatch`方法批量执行SQL语句：

```go
sqls := []string{
    "INSERT INTO users (name, age) VALUES (@name, @age)",
    "SELECT id FROM users WHERE name = @name",
    "UPDATE users SET age = @newage WHERE id = @id",
}

input := map[string]interface{}{
    "name": "张三",
    "age": 30,
    "newage": 31,
}

result, err := dbInstance.ExecuteBatch(sqls, input)
if err != nil {
    log.Fatalf("批量执行失败: %v", err)
}
```

## 6. SQL模板

### 6.1 参数替换

db模块支持使用`@参数名`的形式在SQL语句中定义参数，执行时会自动替换为对应的数据库参数格式：

```go
// MySQL: SELECT * FROM users WHERE id = ?
// Oracle: SELECT * FROM users WHERE id = :id
// PostgreSQL: SELECT * FROM users WHERE id = $1
sql := "SELECT * FROM users WHERE id = @id AND name = @name"
input := map[string]interface{}{"id": 1, "name": "张三"}
rows, err := dbInstance.Query(sql, input)
```

### 6.2 手动替换参数

如果需要手动替换参数，可以使用`Replace`方法：

```go
sql := "SELECT * FROM users WHERE id = ? AND name = ?"
args := []interface{}{1, "张三"}
newSql := dbInstance.Replace(sql, args)
fmt.Println(newSql) // 输出: SELECT * FROM users WHERE id = 1 AND name = '张三'
```

## 7. 查询结果处理

查询结果`QueryRows`提供了丰富的方法来处理查询结果：

```go
// 获取结果行数
count := rows.Len()

// 获取第一行
firstRow := rows.Get(0)

// 遍历结果
for i := 0; i < rows.Len(); i++ {
    row := rows.Get(i)
    
    // 获取字段值
    id := row.GetInt("id")
    name := row.GetString("name")
    age := row.GetInt("age")
    createdAt := row.GetString("created_at")
    
    // 获取原始值
    rawValue := row["id"]
}
```

## 8. 完整示例

### 8.1 基本CRUD示例

```go
package main

import (
    "fmt"
    "log"
    "github.com/micro-plat/lib4go/db"
)

func main() {
    // 创建数据库连接
    dbInstance, err := db.NewDB("sqlite3", ":memory:", 10, 5, 60)
    if err != nil {
        log.Fatalf("创建数据库连接失败: %v", err)
    }
    defer dbInstance.Close()

    // 创建表
    createTableSQL := `
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        age INTEGER
    )
    `
    _, err = dbInstance.Execute(createTableSQL, nil)
    if err != nil {
        log.Fatalf("创建表失败: %v", err)
    }

    // 插入数据
    insertSQL := "INSERT INTO users (name, age) VALUES (@name, @age)"
    insertInput := map[string]interface{}{"name": "张三", "age": 30}
    _, err = dbInstance.Execute(insertSQL, insertInput)
    if err != nil {
        log.Fatalf("插入数据失败: %v", err)
    }

    // 查询数据
    querySQL := "SELECT * FROM users"
    rows, err := dbInstance.Query(querySQL, nil)
    if err != nil {
        log.Fatalf("查询数据失败: %v", err)
    }

    // 输出结果
    for i := 0; i < rows.Len(); i++ {
        row := rows.Get(i)
        id := row.GetInt("id")
        name := row.GetString("name")
        age := row.GetInt("age")
        fmt.Printf("ID: %d, Name: %s, Age: %d\n", id, name, age)
    }

    // 更新数据
    updateSQL := "UPDATE users SET age = @age WHERE id = @id"
    updateInput := map[string]interface{}{"id": 1, "age": 31}
    _, err = dbInstance.Execute(updateSQL, updateInput)
    if err != nil {
        log.Fatalf("更新数据失败: %v", err)
    }

    // 删除数据
    deleteSQL := "DELETE FROM users WHERE id = @id"
    deleteInput := map[string]interface{}{"id": 1}
    _, err = dbInstance.Execute(deleteSQL, deleteInput)
    if err != nil {
        log.Fatalf("删除数据失败: %v", err)
    }
}
```

### 8.2 事务示例

```go
package main

import (
    "fmt"
    "log"
    "github.com/micro-plat/lib4go/db"
)

func main() {
    // 创建数据库连接
    dbInstance, err := db.NewDB("sqlite3", ":memory:", 10, 5, 60)
    if err != nil {
        log.Fatalf("创建数据库连接失败: %v", err)
    }
    defer dbInstance.Close()

    // 创建表
    createTableSQL := `
    CREATE TABLE accounts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        balance INTEGER DEFAULT 0
    )
    `
    _, err = dbInstance.Execute(createTableSQL, nil)
    if err != nil {
        log.Fatalf("创建表失败: %v", err)
    }

    // 插入测试数据
    insertSQL := "INSERT INTO accounts (name, balance) VALUES (@name, @balance)"
    _, err = dbInstance.Execute(insertSQL, map[string]interface{}{"name": "张三", "balance": 1000})
    if err != nil {
        log.Fatalf("插入数据失败: %v", err)
    }
    _, err = dbInstance.Execute(insertSQL, map[string]interface{}{"name": "李四", "balance": 1000})
    if err != nil {
        log.Fatalf("插入数据失败: %v", err)
    }

    // 转账操作（使用事务）
    trans, err := dbInstance.Begin()
    if err != nil {
        log.Fatalf("创建事务失败: %v", err)
    }

    // 张三转出100元
    sql1 := "UPDATE accounts SET balance = balance - 100 WHERE name = @name1"
    _, err = trans.Execute(sql1, map[string]interface{}{"name1": "张三"})
    if err != nil {
        trans.Rollback()
        log.Fatalf("转账失败: %v", err)
    }

    // 李四转入100元
    sql2 := "UPDATE accounts SET balance = balance + 100 WHERE name = @name2"
    _, err = trans.Execute(sql2, map[string]interface{}{"name2": "李四"})
    if err != nil {
        trans.Rollback()
        log.Fatalf("转账失败: %v", err)
    }

    // 提交事务
    if err := trans.Commit(); err != nil {
        log.Fatalf("提交事务失败: %v", err)
    }

    // 查看结果
    querySQL := "SELECT * FROM accounts"
    rows, err := dbInstance.Query(querySQL, nil)
    if err != nil {
        log.Fatalf("查询数据失败: %v", err)
    }

    for i := 0; i < rows.Len(); i++ {
        row := rows.Get(i)
        name := row.GetString("name")
        balance := row.GetInt("balance")
        fmt.Printf("%s的余额: %d\n", name, balance)
    }
}
```

## 9. 注意事项

1. **连接管理**：使用完数据库连接后，应调用`Close`方法关闭连接，释放资源。

2. **事务处理**：事务操作过程中如果发生错误，应及时调用`Rollback`回滚事务，避免数据不一致。

3. **SQL注入防护**：db模块自动处理参数转义，应始终使用参数化查询，避免直接拼接SQL语句。

4. **错误处理**：执行数据库操作时应始终检查返回的错误。

5. **性能优化**：
   - 批量操作比单条操作更高效
   - 合理设置连接池参数
   - 避免在循环中执行大量数据库操作

6. **Oracle数据库**：使用Oracle数据库时，需要按照sys.db.go文件中的说明安装Oracle客户端并配置环境变量。

## 10. API参考

### 核心接口

- `IDB`: 数据库操作接口
- `IDBTrans`: 事务操作接口
- `IDBExecuter`: 数据库执行接口

### 主要方法

- `NewDB`: 创建数据库连接
- `Query`: 执行查询
- `Execute`: 执行SQL操作
- `Scalar`: 获取单个值
- `InsertBatch`: 批量插入
- `UpdateBatch`: 批量更新
- `Begin`: 创建事务
- `Close`: 关闭连接

更多详细API请参考源代码注释。