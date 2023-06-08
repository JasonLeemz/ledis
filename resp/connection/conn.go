package connection

import (
	"ledis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn net.Conn

	// waiting until reply finished
	waitingReply wait.Wait

	// lock while handler sending response
	mu sync.Mutex

	// selected db
	selectedDB int
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(time.Second * 10)
	_ = c.conn.Close()
	return nil
}
func (c *Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}

	// 要保证同一时刻只有一个协程写数据
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()
	_, err := c.conn.Write(bytes)

	return err

}

func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c *Connection) SelectDB(dbNUM int) {
	c.selectedDB = dbNUM
}
