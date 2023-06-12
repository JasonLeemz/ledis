package handler

import (
	"context"
	"io"
	"ledis/cluster"
	"ledis/config"
	"ledis/database"
	databaseface "ledis/interface/database"
	"ledis/lib/logger"
	"ledis/lib/sync/atomic"
	"ledis/resp/connection"
	"ledis/resp/parser"
	"ledis/resp/reply"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	closing    atomic.Boolean
	db         databaseface.Database
}

func MakeHandler() *RespHandler {
	var db databaseface.Database
	if config.Properties.Self != "" && len(config.Properties.Peers) > 0 {
		db = cluster.MakeClusterDatabase()
	} else {
		db = database.NewStandAloneDatabase()
	}
	return &RespHandler{
		db: db,
	}
}

// 关闭单个客户端
func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}

	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	//死循环监听管道
	for payload := range ch {
		if payload.Err != nil {
			// 错误逻辑
			// EOF 代表做四次挥手
			// use of closed network connection 使用了未关闭的连接
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {

				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
			}

			// 协议错误
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		// exec
		if payload.Data == nil {
			continue
		}
		rr, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require mulit bulk reply")
			continue
		}
		result := r.db.Exec(client, rr.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

// Close 关闭所有客户端
func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(
		func(key, value any) bool {
			client := key.(*connection.Connection)
			_ = client.Close()

			return true
		})

	r.db.Close()
	return nil
}
