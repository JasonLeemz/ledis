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

	// subscribing channels
	subs map[string]bool

	// password may be changed by CONFIG command during runtime, so store the password
	password string

	// queued commands for `multi`
	multiState bool
	queue      [][][]byte
	watching   map[string]uint32

	// selected db
	selectedDB int
}

func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

func (c Connection) Write(bytes []byte) error {
	if len(bytes) == 0 {
		return nil
	}
	c.mu.Lock()
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
		c.mu.Unlock()
	}()

	_, err := c.conn.Write(bytes)
	return err
}

func (c Connection) GetDBIndex() int {
	return c.selectedDB
}

func (c Connection) SelectDB(dbNum int) {
	c.selectedDB = dbNum
}

// Close disconnect with the client
func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
