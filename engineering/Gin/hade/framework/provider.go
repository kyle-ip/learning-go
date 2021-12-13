package framework

// NewInstance 创建新实例，所有服务容器的创建服务
type NewInstance func(...interface{}) (interface{}, error)

// ServiceProvider 服务提供者需要实现的接口
type ServiceProvider interface {

	// Register 在服务容器中注册实例化服务的方法，是否在注册时就实例化，需参考 IsDefer 接口。
	Register(Container) NewInstance

	// Boot 在调用实例化服务时调用，可放入准备工作（基础配置，初始化参数），返回 error 表示服务实例化就会实例化失败。
	Boot(Container) error

	// IsDefer 决定是否在注册时实例化，true 即在第一次 make 时进行实例化操作，false 则在注册时实例化。
	IsDefer() bool

	// Params 定义传递给 NewInstance 的参数，可以自定义多个，建议将 container 作为第一个参数。
	Params(Container) []interface{}

	// Name 代表了这个服务提供者的凭证。
	Name() string
}
