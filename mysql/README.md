# 1. 项目简介

1.基于gorm库封装，增加tracing
2.接入统一日志输出

# 2. 详细文档地址

https://jasperxu.github.io/gorm-zh

# 3. 配置文件参数

```yaml
mysqlDB:
  alias: admin-server
  type: mysql  #当前代码里面预制的是mysql
  server: 127.0.0.1
  port: 3306
  user: root
  password: my-secret-pw
  charset: utf8mb4
  database: ets_db
  maxIdleConns: 3
  maxOpenConns: 10
  maxLeftTime: 10s
  timezone: Local
```

# 4. 快速开始

```go
db := gorm.New(&config.Config{
    Type:     "mysql",  
    Server:   "127.0.0.1",
    Port:     3306,
    Database: "db",
    User:     "db",
    Password: "123456",
})

//应用中使用，sess的具体方法查看 https://jasperxu.github.io/gorm-zh
sess = db.Context(context.Background())

```

# 5. 使用示例
参见gorm_test.go