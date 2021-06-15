package best_practise

// Self referential functions and design    http://commandcenter.blogspot.com.au/2014/01/self-referential-functions-and-design.html

import (
    "crypto/tls"
    "fmt"
    "testing"
    "time"
)

// 服务器配置结构

type Server struct {
    Addr     string        // required
    Port     int           // required
    Protocol string        // not null, default TCP
    Timeout  time.Duration // not null, default 30
    MaxConns int           // not null, default 1024
    TLS      *tls.Config   //
}

// ========== 多构造函数 ==========

// 针对以上配置，有多种创建不同配置 Server 的函数签名（Go 不支持重载）：

func NewDefaultServer(addr string, port int) (*Server, error) {
    return &Server{addr, port, "tcp", 30 * time.Second, 100, nil}, nil
}

func NewTLSServer(addr string, port int, tls *tls.Config) (*Server, error) {
    return &Server{addr, port, "tcp", 30 * time.Second, 100, tls}, nil
}

func NewServerWithTimeout(addr string, port int, timeout time.Duration) (*Server, error) {
    return &Server{addr, port, "tcp", timeout, 100, nil}, nil
}

func NewTLSServerWithMaxConnAndTimeout(addr string, port int, maxconns int, timeout time.Duration, tls *tls.Config) (*Server, error) {
    return &Server{addr, port, "tcp", 30 * time.Second, maxconns, tls}, nil
}

// ========== 配置对象 ==========

type OptionConfig struct {
    Protocol string
    Timeout  time.Duration
    Maxconns int
    TLS      *tls.Config
}

// 可以把非必需的选项都放在 OptionConfig 中，只需要一个 NewServer() 的函数即可，在使用前需要构造 Config 对象。

type ConfigServer struct {
    Addr string
    Port int
    Conf *OptionConfig
}

func NewServer(addr string, port int, conf *OptionConfig) (*ConfigServer, error) {
    //...
    return nil, nil
}

func Case12() {
    //Using the default configuratrion
    srv1, _ := NewServer("localhost", 9000, nil)

    conf := OptionConfig{Protocol: "tcp", Timeout: 60}
    srv2, _ := NewServer("locahost", 9000, &conf)

    fmt.Println(srv1, srv2)
}

// ========== Builder 模式 ==========

type ServerBuilder struct {
    Server
}

func (sb *ServerBuilder) Create(addr string, port int) *ServerBuilder {
    sb.Server.Addr = addr
    sb.Server.Port = port
    // 其它代码设置其它成员的默认值
    return sb
}

func (sb *ServerBuilder) WithProtocol(protocol string) *ServerBuilder {
    sb.Server.Protocol = protocol
    return sb
}

func (sb *ServerBuilder) WithMaxConn(maxconn int) *ServerBuilder {
    sb.Server.MaxConns = maxconn
    return sb
}

func (sb *ServerBuilder) WithTimeOut(timeout time.Duration) *ServerBuilder {
    sb.Server.Timeout = timeout
    return sb
}

func (sb *ServerBuilder) WithTLS(tls *tls.Config) *ServerBuilder {
    sb.Server.TLS = tls
    return sb
}

func (sb *ServerBuilder) Build() (Server, error) {
    return sb.Server, nil
}

func Test13(t *testing.T) {
    sb := ServerBuilder{}
    server, _ := sb.Create("127.0.0.1", 8080).
        WithProtocol("udp").
        WithMaxConn(1024).
        WithTimeOut(30 * time.Second).
        Build()
    fmt.Println(server)
}

// ========== Functional Options ==========

// 高阶函数：传入一个参数返回一个函数，返回函数会设置自己的 Server 参数。
// 例如调用其中的一个函数 MaxConns(30) 时，其返回值是一个 func(s* Server) { s.MaxConns = 30 } 的函数。

// Functional Options 的优势：
// - 直觉式的编程。
// - 高度的可配置化。
// - 很容易维护和扩展。
// - 自文档。
// - 新来的人很容易上手。
// - 直观，无困惑（nil？空？）。

type Option func(*Server)

func Protocol(p string) Option {
    return func(s *Server) {
        s.Protocol = p
    }
}
func Timeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.Timeout = timeout
    }
}
func MaxConns(maxconns int) Option {
    return func(s *Server) {
        s.MaxConns = maxconns
    }
}
func TLS(tls *tls.Config) Option {
    return func(s *Server) {
        s.TLS = tls
    }
}

func NewServerFP(addr string, port int, options ...Option) (*Server, error) {

    // 有一个可变参数 options 可以传出多个上面的函数，然后使用 for-loop 来设置 Server 对象。
    srv := Server{
        Addr:     addr,
        Port:     port,
        Protocol: "tcp",
        Timeout:  30 * time.Second,
        MaxConns: 1000,
        TLS:      nil,
    }
    for _, option := range options {
        option(&srv)
    }
    //...
    return &srv, nil
}

func Test14(t *testing.T) {
    s1, _ := NewServerFP("localhost", 1024)
    s2, _ := NewServerFP("localhost", 2048, Protocol("udp"))
    s3, _ := NewServerFP("0.0.0.0", 8080, Timeout(300*time.Second), MaxConns(1000))
    fmt.Println(s1, s2, s3)
}
