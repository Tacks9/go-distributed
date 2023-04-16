# go-distributed

### 实现过程

- 服务注册
  - 日志服务 
    - 启动命令 `cd cmd/logservice` `go run main.go`
    - 入口文件 [`main.go`](../cmd/logservice/main.go)
    - 业务逻辑 [`server.go`](../log/server.go)
    - 模拟请求
        ```shell
        curl --location --request POST 'http://localhost:4000/log' \
        --header 'Content-Type: text/plain' \
        --data-raw '测试日志记录'
        ```
    - 流程图
      - cmd 
        - `logservice main.go`
          - pkg
            - log
            - service
  - 注册服务
    - 启动命令 `cd cmd/registryservice` `go run main.go`
    - 入口文件 [`main.go`](../cmd/registryservice/main.go)
    - 业务逻辑 [`server.go`](../registry/server.go)
    - 类型定义 [`registration.go`](../registry/resgistration.go)
    - 模拟请求
      ```shell
      curl --location --request POST 'http://localhost:3000/services' \
      --header 'Content-Type: application/json' \
      --data-raw '{
          "serviceName": "Test Service",
          "serviceURL": "http://localhost:5000"
      }'
      ```
    - 流程图
      - cmd
        - 日志服务 `logservice main.go`
          - pkg
            - log
            - service
        - 注册服务 `registryservice main.go`
          - pkg
            - registry