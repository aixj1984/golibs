
@echo off
echo "start lint....."
:: 转中文输出
chcp 65001

echo 当前盘符和路径：%~dp0

set arg=%1
IF "%arg%"=="" (
  set arg=lint
)

IF NOT "%arg%"=="lint" IF NOT "%arg%"=="cover" (
  echo "-------> 参数错误，参数应为lint或cover <---------"
  exit /b 3
)

IF "%arg%"=="lint" (
    go mod tidy
    IF %errorlevel% NEQ 0 (
        Echo "-------> 构建环境失败 <---------"
        exit /b 1
    ) ELSE (
        Echo "-------> 执行LINT <---------"
    )

    golangci-lint run ||  (
        Echo "-------> Lint存在异常，请检查 <---------"
        exit /b 2
    ) 
    
    Echo "-------> Lint检测通过 <---------"
    
)

IF "%arg%"=="cover" (
    Echo "-------> 执行Go Test <---------"
    @REM go test -cover ./...
    :: go test -coverprofile="cover.out" ./...
    :: go test -coverprofile="cover.out" -coverpkg="$(go list ./... | grep -v components | grep -v utils)" -covermode=atomic ./...
    go test -coverprofile="cover.out" -coverpkg="./..." -covermode=atomic ./...
    go tool cover -html="cover.out" 
)
