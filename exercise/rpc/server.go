package geerpc

// ========== 服务端 ==========

import (
	"encoding/json"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

// MagicNumber 魔数，用于标识协议
const MagicNumber = 0x3bef5c

// Option 选项
type Option struct {
	// MagicNumber 魔术
	MagicNumber int

	// CodecType 编解码类型
	CodecType codec.Type // client may choose different Codec to encode body
}

// DefaultOption 默认选项
var DefaultOption = &Option{
	MagicNumber: MagicNumber,
	CodecType:   codec.GobType,
}

// Server RPC 服务端
type Server struct{}

// NewServer Server 构造函数
func NewServer() *Server {
	return &Server{}
}

// DefaultServer 默认服务端指针
var DefaultServer = NewServer()

// ServeConn 与客户端建立连接，处理请求
func (server *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option

	// 通过 JSON 反序列化得到 Option 实例、判断魔数、获取编解码类型。
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}
	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number %x", opt.MagicNumber)
		return
	}
	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
		return
	}

	// 读取、处理、回复请求。
	server.serveCodec(f(conn))
}

// invalidRequest 预留用于错误请求的响应参数
var invalidRequest = struct{}{}

// serveCodec 处理请求：一次连接中，允许接收多个请求。
func (server *Server) serveCodec(cc codec.Codec) {

	// 响应请求时使用的锁。
	sending := new(sync.Mutex)

	// 等待所有请求处理完成。
	wg := new(sync.WaitGroup)

	for {

		// 读取请求头、参数等。
		req, err := server.readRequest(cc)
		if err != nil {

			// 在 header 解析失败时终止循环。
			if req == nil {
				break
			}
			req.h.Error = err.Error()

			// 响应请求时锁定：回复请求的报文必须逐个发送，并发容易导致多个回复报文交织在一起，客户端无法解析。
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)

		// 创建协程并发处理请求。
		go server.handleRequest(cc, req, sending, wg)
	}
	wg.Wait()
	_ = cc.Close()
}

// request 请求信息
type request struct {

	// 请求头
	h *codec.Header

	// 请求参数、响应参数
	argv, replyv reflect.Value
}

// readRequestHeader 读取请求头
func (server *Server) readRequestHeader(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

// readRequest 读取请求
func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHeader(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

// sendResponse 发送响应
func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error:", err)
	}
}

// handleRequest 处理请求
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// Accept 监听和接收连接
func (server *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
			return
		}
		// 每监听到一个客户端连接，交由 ServeConn 处理。
		go server.ServeConn(conn)
	}
}

// Accept 监听和接收连接
func Accept(lis net.Listener) { DefaultServer.Accept(lis) }
