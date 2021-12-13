package codec

// ========== JSON ==========

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
)

// JsonCodec 编解码器
type JsonCodec struct {

	// conn 连接
	conn io.ReadWriteCloser

	// buf 写缓冲区
	buf *bufio.Writer

	// dec 解码器
	dec *json.Decoder

	// enc 编码器
	enc *json.Encoder
}

var _ Codec = (*JsonCodec)(nil)

// NewJsonCodec JsonCodec 构造函数
func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &JsonCodec{
		conn: conn,
		buf:  buf,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(buf),
	}
}

// ReadHeader 解码消息头
func (c *JsonCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 解码消息体
func (c *JsonCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// Write 编码消息
func (c *JsonCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: json error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

// Close 关闭连接：io.Closer
func (c *JsonCodec) Close() error {
	return c.conn.Close()
}
