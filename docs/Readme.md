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
  - 注册服务-服务端
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
  - 注册服务-客户端
    - 流程图
      - cmd
        - 启动注册服务 `cd cmd/registryservice` `go run main.go` 
          - 监听 http://localhost:3000/services
        - 启动日志服务 `cd cmd/logservice` `go run main.go`
          - 调用 service/service.go 启动日志服务
          - 监听 http://localhost:4000/log
          - 设置日志服务 handle 进行日志记录
          - 向服务中心注册日志服务 调用 http://localhost:3000/services 进行 add
  - 移除服务-客户端
    - 流程图
      - 启动注册服务 `cd cmd/registryservice` `go run main.go` 
      - 注册日志服务 `cd cmd/logservice` `go run main.go`
      - 关闭日志服务
      - 触发注册中心移除日志服务
  - 学生成绩服务
    - 流程图
      - cmd 
        - 1、启动注册服务  `cd cmd/registryservice` `go run main.go` 
        - 2、启动日志服务  `cd cmd/logservice` `go run main.go` 
        - 3、启动学生服务  `cd cmd/gradingservice` `go run main.go` 
  - 服务的发现
    - 服务之间想要相互依赖的时候可以在 `Registration.RequiredServices` 中
  - 服务的更新
    - 服务依赖变化的时候进行通知 notify
  - Web 学生服务 `protal`
    - 增加 学生 HTML 模版，实现 WEB 界面
    - 依赖
      - grading
      - log
  - 服务的检测
    - 心跳检查，启动 go routinue 对每一个服务进行 get 请求