package framework

import (
	"log"
	"net/http"
	"strings"
)

// Core represent core struct
type Core struct {
	router      map[string]*Tree    // all routers
	middlewares []ControllerHandler // 设置中间件。
}

// NewCore 初始化框架核心结构
func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	core := &Core{router: router}
	return core
}

// Use 注册中间件
func (c *Core) Use(middlewares ...ControllerHandler) {
	// c.middlewares = append(c.middlewares, middlewares...)
	c.middlewares = middlewares
}

// === http method wrap

func (c *Core) Get(url string, handlers ...ControllerHandler) {
	// 将 core 的 middleware 和 handlers 结合起来。
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["GET"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Post(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["POST"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Put(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["PUT"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (c *Core) Delete(url string, handlers ...ControllerHandler) {
	allHandlers := append(c.middlewares, handlers...)
	if err := c.router["DELETE"].AddRouter(url, allHandlers); err != nil {
		log.Fatal("add router error: ", err)
	}
}

// ==== http method wrap end

func (c *Core) Group(prefix string) IGroup {
	return NewGroup(c, prefix)
}

// FindRouteNodeByRequest 匹配路由，如果没有匹配到，返回 nil
func (c *Core) FindRouteNodeByRequest(request *http.Request) *node {
	// uri 和 method 全部转换为大写，保证大小写不敏感
	uri := request.URL.Path
	method := request.Method
	upperMethod := strings.ToUpper(method)

	if methodHandlers, ok := c.router[upperMethod]; ok {
		return methodHandlers.root.matchNode(uri)
	}
	return nil
}

// ServeHTTP 路由分发
func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// 封装自定义 context
	ctx := NewContext(request, response)

	// 寻找路由
	n := c.FindRouteNodeByRequest(request)
	if n == nil {
		// 如果没有找到，这里打印日志
		ctx.SetStatus(404).Json("not found")
		return
	}

	ctx.SetHandlers(n.handlers)

	params := n.parseParamsFromEndNode(request.URL.Path)
	ctx.SetParams(params)

	// 调用路由函数，如果返回 err 代表存在内部错误，返回 500 状态码
	if err := ctx.Next(); err != nil {
		ctx.Json(500).Json("inner error")
		return
	}
}
