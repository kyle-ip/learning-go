package codec

// ========== 序列化与反序列化 ==========

import "io"

// Header RPC 请求协议头
type Header struct {
	// ServiceMethod 服务名和方法名，格式："Service.Method"
	ServiceMethod string
	// Seq 请求序号
	Seq uint64
	// Error 错误信息（仅用于服务端）
	Error string
}

// Codec 消息编解码接口
type Codec interface {
	// Closer 可关闭的
	io.Closer

	// ReadHeader 解码消息头
	ReadHeader(*Header) error

	// ReadBody 解码消息体
	ReadBody(interface{}) error

	// Write 编码消息
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

// NewCodecFuncMap 编解码工厂，传入编解码类型，传出对应的构造函数
var NewCodecFuncMap map[Type]NewCodecFunc

// init 初始化
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
	NewCodecFuncMap[JsonType] = NewJsonCodec
}
