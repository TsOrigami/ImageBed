# DB/Sql

连接数据库。双检锁单例模式，返回自定义结构体DB

```go
type DB struct {
    Dsn     string
    Connect *sql.DB
}
```