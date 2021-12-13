package codec

// ========== Gob ==========

import (
	"bufio"
	"encoding/gob"
	"io"
	"log"
)

// GobCodec 编解码器
type GobCodec struct {

	// conn 连接
	conn io.ReadWriteCloser

	// buf 写缓冲区
	buf *bufio.Writer

	// dec 解码器
	dec *gob.Decoder

	// enc 编码器
	enc *gob.Encoder
}

var _ Codec = (*GobCodec)(nil)

// NewGobCodec GobCodec 构造函数
func NewGobCodec(conn io.ReadWriteCloser) Codec {
	buf := bufio.NewWriter(conn)
	return &GobCodec{
		conn: conn,
		buf:  buf,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(buf),
	}
}

// ReadHeader 解码消息头
func (c *GobCodec) ReadHeader(h *Header) error {
	return c.dec.Decode(h)
}

// ReadBody 解码消息体
func (c *GobCodec) ReadBody(body interface{}) error {
	return c.dec.Decode(body)
}

// Write 编码消息
func (c *GobCodec) Write(h *Header, body interface{}) (err error) {
	defer func() {
		_ = c.buf.Flush()
		if err != nil {
			_ = c.Close()
		}
	}()
	if err = c.enc.Encode(h); err != nil {
		log.Println("rpc: gob error encoding header:", err)
		return
	}
	if err = c.enc.Encode(body); err != nil {
		log.Println("rpc: gob error encoding body:", err)
		return
	}
	return
}

// Close 关闭连接：io.Closer
func (c *GobCodec) Close() error {
	return c.conn.Close()
}
