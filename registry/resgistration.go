package registry

type ServiceName string

// 服务信息
type Registration struct {
	ServiceName ServiceName // 服务的名称
	ServiceURL  string      // 服务的地址

	RequiredServices []ServiceName // 服务所依赖的其他项 (例如grades依赖的log服务)
	ServiceUpdateURL string        // 接收服务状态信息

	HeartbeatURL string // 心跳检查
}

// 服务更新 每一条
type patchEntry struct {
	Name ServiceName
	URL  string
}

// 每次服务变更 增加或者减少
type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}

// 服务列表
const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	ProtalService  = ServiceName("Protald")
)
