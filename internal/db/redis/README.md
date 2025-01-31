# DB/Redis

数据库连接。双检锁单例模式。返回自定义结构体Redis

```go
type Redis struct {
    Connect  *redis.Client
    Ctx      context.Context
    Remains  int
    addr     string
    password string
    db       int
}
```

