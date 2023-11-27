# 1. 项目简介

1.基于gorm库封装，增加tracing
2.接入统一日志输出

# 2. 详细文档地址

https://jasperxu.github.io/gorm-zh

# 3. 配置文件参数

```yaml
gorm:
  alias: admin-server
  driver: mysql  #当前代码里面预制的是mysql
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

// 如果配置文件中有gorm配置文件，则会默认初始化default

// 如果只希望获取自己定义的初始化DB，可通过一下接口获取
db := gorm.NewEngine(&config.Config{
    Driver:     "mysql",  
    Server:   "127.0.0.1",
    Port:     3306,
    Database: "db",
    User:     "db",
    Password: "123456",
})

// 如果想管理多个DB连接，则使用注册方式进行
gorm.RegisterDataBase("test-mysql",&config.Config{
    Driver:     "mysql",  
    Server:   "127.0.0.1",
    Port:     3306,
    Database: "db",
    User:     "db",
    Password: "123456",
})
// 获取注册的DB
gorm.GetEngine("test-mysql")

// 获取默认DB
gorm.GetEngine()
gorm.GetEngine("defalut")

//应用中使用，sess的具体方法查看 https://jasperxu.github.io/gorm-zh
sess = db.Context(context.Background())

```

# 5. 使用示例
参见gorm_test.go

# 6. 如果需要使用sqlite数据库，则需要环境中包括gcc编译环境

windows推荐直接下载https://github.com/niXman/mingw-builds-binaries/releases对应的版本， 如x86_64-13.2.0-release-posix-seh-msvcrt-rt_v11-rev0.7z，下载后，解压到某个目录，配置环境变量，指向该目录下的bin目录
