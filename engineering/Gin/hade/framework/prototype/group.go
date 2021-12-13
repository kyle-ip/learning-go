package framework

// IGroup 代表前缀分组
type IGroup interface {

	// 实现 HttpMethod 方法

	Get(string, ...ControllerHandler)
	Post(string, ...ControllerHandler)
	Put(string, ...ControllerHandler)
	Delete(string, ...ControllerHandler)

	// Group 实现嵌套
	Group(string) IGroup

	// Use 嵌套中间件
	Use(middlewares ...ControllerHandler)
}

// Group struct 实现了IGroup
type Group struct {
	core        *Core               // 指向 core 结构
	parent      *Group              // 指向上一个 Group，如果有的话
	prefix      string              // group 的通用前缀
	middlewares []ControllerHandler // 存放中间件
}

// NewGroup 初始化 Group
func NewGroup(core *Core, prefix string) *Group {
	return &Group{
		core:        core,
		parent:      nil,
		prefix:      prefix,
		middlewares: []ControllerHandler{},
	}
}

func (g *Group) Get(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Get(uri, allHandlers...)
}

func (g *Group) Post(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Post(uri, allHandlers...)
}

func (g *Group) Put(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Put(uri, allHandlers...)
}

func (g *Group) Delete(uri string, handlers ...ControllerHandler) {
	uri = g.getAbsolutePrefix() + uri
	allHandlers := append(g.getMiddlewares(), handlers...)
	g.core.Delete(uri, allHandlers...)
}

// getAbsolutePrefix 获取当前 group 的绝对路径
func (g *Group) getAbsolutePrefix() string {
	if g.parent == nil {
		return g.prefix
	}
	return g.parent.getAbsolutePrefix() + g.prefix
}

// getMiddlewares 获取 group 的 middleware，即除 Get/Post/Put/Delete 外设置的 middleware
func (g *Group) getMiddlewares() []ControllerHandler {
	if g.parent == nil {
		return g.middlewares
	}
	return append(g.parent.getMiddlewares(), g.middlewares...)
}

func (g *Group) Group(uri string) IGroup {
	cgroup := NewGroup(g.core, uri)
	cgroup.parent = g
	return cgroup
}

// Use 注册中间件
func (g *Group) Use(middlewares ...ControllerHandler) {
	g.middlewares = middlewares
}
