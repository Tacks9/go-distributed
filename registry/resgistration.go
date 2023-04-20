package registry

type ServiceName string

// 服务信息
type Registration struct {
	ServiceName ServiceName
	ServiceURL  string
}

// 服务列表
const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
)
