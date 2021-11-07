package server

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
)

// command of message
const (
	CommandPing = 100
	CommandPong = 101
	writeWait   = time.Second * 10
	readWait    = time.Minute * 2
)

type Option struct {
	id     string
	listen string
}

// Server is a websocket implement of the Server
type Server struct {
	once    sync.Once
	id      string
	address string
	sync.Mutex
	// 会话列表
	users map[string]net.Conn
}

// Start server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "Server",
		"listen": s.address,
		"id":     s.id,
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			_ = conn.Close()
			return
		}
		// 读取 userId
		user := r.URL.Query().Get("user")
		if user == "" {
			_ = conn.Close()
			return
		}

		// 添加到会话管理中
		old, ok := s.addUser(user, conn)
		if ok {
			// 断开旧的连接
			_ = old.Close()
			log.Infof("close old connection %v", old.RemoteAddr())
		}
		log.Infof("user %s in from %v", user, conn.RemoteAddr())

		go func(user string, conn net.Conn) {
			err := s.readLoop(user, conn)
			if err != nil {
				log.Warn("readLoop - ", err)
			}
			_ = conn.Close()
			// 删除用户
			s.delUser(user)

			log.Infof("connection of %s closed", user)
		}(user, conn)
	})
	log.Infoln("started")
	return http.ListenAndServe(s.address, mux)
}

func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()
	old, ok := s.users[user]
	s.users[user] = conn
	return old, ok
}

func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.users, user)
}

func (s *Server) readLoop(user string, conn net.Conn) error {
	for {
		// 要求：客户端必须在指定时间内发送一条消息过来，可以是ping，也可以是正常数据包
		_ = conn.SetReadDeadline(time.Now().Add(readWait))

		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}
		if frame.Header.OpCode == ws.OpPing {
			// 返回一个pong消息
			_ = wsutil.WriteServerMessage(conn, ws.OpPong, nil)
			logrus.Info("write a pong...")
			continue
		}
		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote side close the conn")
		}
		logrus.Info(frame.Header)

		if frame.Header.Masked {
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}
		// 接收文本帧内容
		if frame.Header.OpCode == ws.OpText {
			go s.handle(user, string(frame.Payload))
		} else if frame.Header.OpCode == ws.OpBinary {
			go s.handleBinary(user, frame.Payload)
		}
	}
}

// 广播消息
func (s *Server) handle(user string, message string) {
	logrus.Infof("recv message %s from %s", message, user)
	s.Lock()
	defer s.Unlock()
	broadcast := fmt.Sprintf("%s -- FROM %s", message, user)
	for u, conn := range s.users {
		if u == user { // 不发给自己
			continue
		}
		logrus.Infof("send to %s : %s", u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}

func (s *Server) handleBinary(user string, message []byte) {
	logrus.Infof("recv message %v from %s", message, user)
	s.Lock()
	defer s.Unlock()
	// handle ping request
	i := 0
	command := binary.BigEndian.Uint16(message[i : i+2])
	i += 2
	payloadLen := binary.BigEndian.Uint32(message[i : i+4])
	logrus.Infof("command: %v payloadLen: %v", command, payloadLen)
	if command == CommandPing {
		u := s.users[user]
		// return pong
		err := wsutil.WriteServerBinary(u, []byte{0, CommandPong, 0, 0, 0, 0})
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}

func (s *Server) writeText(conn net.Conn, message string) error {
	// 创建文本帧数据
	f := ws.NewTextFrame([]byte(message))
	err := conn.SetWriteDeadline(time.Now().Add(writeWait))
	if err != nil {
		return err
	}
	return ws.WriteFrame(conn, f)
}

// NewServerStartCmd creates a new http server command
func NewServerStartCmd(ctx context.Context, version string) *cobra.Command {
	opts := &Option{}

	cmd := &cobra.Command{
		Use:   "chat",
		Short: "Starts a chat server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts, version)
		},
	}
	cmd.PersistentFlags().StringVarP(&opts.id, "serverid", "i", "demo", "server id")
	cmd.PersistentFlags().StringVarP(&opts.listen, "listen", "l", ":8000", "listen address")
	return cmd
}

func RunServerStart(ctx context.Context, opts *Option, version string) error {

	// init
	server := &Server{
		id:      opts.id,
		address: opts.listen,
		users:   make(map[string]net.Conn, 100),
	}

	// shutdown
	defer server.once.Do(func() {
		server.Lock()
		defer server.Unlock()
		for _, conn := range server.users {
			_ = conn.Close()
		}
	})

	// start
	return server.Start()
}
