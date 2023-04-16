# go-distributed

### 实现过程

- 服务注册
  - 日志服务 
    - 启动命令 `cd cmd/logservice` `go run main.go`
    - 入口文件 [`main.go`](cmd/logservice/main.go)
    - 业务逻辑 [`server.go`](log/server.go)
    - 模拟请求
        ```
        curl --location --request POST 'http://localhost:4000/log' \
        --header 'Content-Type: text/plain' \
        --data-raw '测试日志记录'
        ```