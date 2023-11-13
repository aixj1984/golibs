conf是对viper三方库的封装，主要是统一日志文件类型，统一配置文件存储的路径和文件名，便于工程的一致性。
同时该封装没有打破viper的特性和能力，封装提供了获取viper对象的接口；对于需求使用最多的配置文件解析和某段配置的解析，提供了泛型的配置参数获取。


## 特性

1. 无侵入，保留viper提供的能力：
2. 统一配置文件类型、路径和文件名：***./etc/conf.yaml***
3. 提供直接获取viper对象接口
4. 提供统一的泛型获取参数接口
5. 提供按某段来获取subviper；提供泛型方式，按段来获取配置

## 快速上手

现有配置文件内容如下：
```yaml
clothing:
  jacket: leather
  trousers: denim
  times: 10s

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
//获取整个配置文件的viper对象
viper := conf.GetViper()

//获取某段配置文件的viper对象
subViper := conf.GetSubViper("log")

//可以使用viper原生方法，获取某个字段
viper.GetString("clothing.jacket")

//如果想直接使用自己指定的文件初始化对象
viper.New("file path")

//


```


```go
// 泛型获取配置对象
type testConfig struct {
	Jacket   string        `json:"jacket"`
	Trousers string        `json:"trousers"`
	Times    time.Duration `json:"times"`
}

subv := conf.GetSubViper("clothing")
if subv == nil {
  log.Fatalf("No 'database' key in config")
}

var data testConfig

if err := subv.Unmarshal(&data); err != nil {
  log.Fatalf("Unable to decode into struct, %s", err.Error())
} else {
  log.Printf("config %v ", data)
}


```



## 开始使用

```SQL
 go get g3.io/golibs/conf
```


