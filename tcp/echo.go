package tcp

import (
	"bufio"
	"context"
	"io"
	"ledis/lib/logger"
	"net"
	"sync"
	"time"

	"ledis/lib/sync/atomic"
	"ledis/lib/sync/wait"
)

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	_ = e.Conn.Close()

	return nil
}

type EchoHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
}

func NewHandler() *EchoHandler {
	return &EchoHandler{}
}

func (handler *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		_ = conn.Close()
	}

	client := &EchoClient{
		Conn: conn,
	}

	handler.activeConn.Store(client, struct{}{})

	reader := bufio.NewReader(conn)
	for true {
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("Connection closed")
				handler.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		_, _ = conn.Write([]byte(readString))
		client.Waiting.Done()
	}
}

func (handler *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Conn.Close()
		return true
	})

	return nil
}
