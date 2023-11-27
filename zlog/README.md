zlog 是对zap三方库的封装，主要是统一日志文件的引入和调用


## 特性

1. 无侵入，保留zap提供的能力：
2. 统一配置文件类型、路径和文件名：***./etc/conf.yaml***
3. 统一配置文件中的段配置
```yaml
log:
  appName: admin-server
  level: -1
  debug: true
  logPath: ./log/admin-server.log
  maxSize: 128
  maxAge: 7
  maxBackups: 30
  compress: false
```
4. 提供直接获取zap对象接口
5. 每个错误种类，提供三种不同类型的日志输出：Debug/DebugF/DebugO


## 快速上手

现有配置文件内容如下：
```yaml
log:
  appName: test-server
  level: -1
  debug: true
  logPath: ./log/test-server.log
  maxSize: 128
  maxAge: 7
  maxBackups: 30
  compress: false
```

```go
//如果有自定义的配置信息，可进行重新初始化

InitLogger(&Config{
  LogPath:    "./log/test.log",
  AppName:    "log-sample",
  Level:      -1,
  MaxSize:    0,
  MaxAge:     0,
  MaxBackups: 0,
  Compress:   true,
})



```


```go
  // 按map输出
  Debug("info log", Fields{"abc": 11})
  Info("info log", Fields{"abc": 11})
  Warn("info log", Fields{"abc": 11})
  Error("info log", Fields{"abc": 11})
	

  // 按字符串输出
  Debugf("format log : %d", 12)
  Infof("format log : %d", 12)
  Warnf("format log : %d", 12)
  Errorf("format log : %d", 12)

  // 输出任意对象
  DebugO("object log ", time.Now())
  InfoO("object log ", time.Now())
  WarnO("object log ", time.Now())
  ErrorO("object log ", time.Now())

  // Panic和Fatal的输出都会直接导致程序退出，所以慎用

  //FatalO("object log ", time.Now())
  //DPanic("info log", Fields{"abc": 11})
  //Panic("info log", Fields{"abc": 11})
  //Fatal("info log", Fields{"abc": 11})

  Logger().Debug("info log", Fields{"abc": 11})


```


## 开始使用

```go
 go get github.com/aixj1984/golibs/zlog
```


## 代码覆盖率
go test -cover ./...

go test -coverprofile=coverage ./...

go tool cover -html=coverage