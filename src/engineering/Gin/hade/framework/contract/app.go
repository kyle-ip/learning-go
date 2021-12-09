package contract

// AppKey 字符串凭证
const AppKey = "hade:app"

// App 定义接口
type App interface {

	// Version 当前版本
	Version() string

	//BaseFolder 项目基础地址
	BaseFolder() string

	// ConfigFolder 配置文件的路径
	ConfigFolder() string

	// LogFolder 日志所在路径
	LogFolder() string

	// ProviderFolder 业务自己的服务提供者地址
	ProviderFolder() string

	// MiddlewareFolder 业务自己定义的中间件
	MiddlewareFolder() string

	// CommandFolder 业务定义的命令
	CommandFolder() string

	// RuntimeFolder 业务的运行中间态信息
	RuntimeFolder() string

	// TestFolder 测试所需要的信息
	TestFolder() string
}
