package aof

import (
	"io"
	"ledis/config"
	databaseface "ledis/interface/database"
	"ledis/lib/logger"
	"ledis/lib/utils"
	"ledis/resp/connection"
	"ledis/resp/parser"
	"ledis/resp/reply"
	"os"
	"strconv"
	"sync"
)

type CmdLine = [][]byte

const aofBufSize = 1 << 16 // 65535
type payload struct {
	cmdLine CmdLine
	dbIndex int
}

type AofHandler struct {
	db          databaseface.Database
	aofChan     chan *payload
	aofFile     *os.File
	aofFilename string
	// aof goroutine will send msg to main goroutine through this channel when aof tasks finished and ready to shutdown
	aofFinished chan struct{}
	// pause aof for start/finish aof rewrite progress
	pausingAof sync.RWMutex
	currentDB  int
}

func NewAofHandler(database databaseface.Database) (*AofHandler, error) {
	handler := &AofHandler{
		aofFilename: config.Properties.AppendFilename,
		db:          database,
	}

	// 恢复AOF
	handler.LoadAof()

	aofile, err := os.OpenFile(handler.aofFilename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	handler.aofFile = aofile

	// channel 初始化
	handler.aofChan = make(chan *payload, aofBufSize)
	go func() {
		handler.handleAof()
	}()

	return handler, nil
}

// AddAof 落盘动作放到channel中，由另外的协程去操作
func (handler *AofHandler) AddAof(dbIndex int, cmd CmdLine) {
	// 判断是否开启aof
	if config.Properties.AppendOnly && handler.aofChan != nil {
		handler.aofChan <- &payload{
			cmdLine: cmd,
			dbIndex: dbIndex,
		}
	}

}

func (handler *AofHandler) handleAof() {
	handler.currentDB = 0
	for p := range handler.aofChan {
		if p.dbIndex != handler.currentDB {
			line := utils.ToCmdLine("select", strconv.Itoa(p.dbIndex))
			// *2 $5 select $1 3
			data := reply.MakeMultiBulkReply(line).ToBytes()
			_, err := handler.aofFile.Write(data)
			if err != nil {
				logger.Error(err)
				continue
			}
			handler.currentDB = p.dbIndex
		}

		data := reply.MakeMultiBulkReply(p.cmdLine).ToBytes()
		_, err := handler.aofFile.Write(data)
		if err != nil {
			logger.Error(err)
		}
	}
}

// LoadAof ...
func (handler *AofHandler) LoadAof() {
	aofile, err := os.Open(handler.aofFilename)
	if err != nil {
		logger.Error(err)
		return
	}

	defer func() {
		aofile.Close()
	}()

	ch := parser.ParseStream(aofile)
	fakeConn := &connection.Connection{}
	for p := range ch {
		if p.Err != nil {
			if p.Err == io.EOF {
				// 文件已经结束
				break
			}
			logger.Error(p.Err)
			continue
		}

		if p.Data == nil {
			logger.Error("empty payload")
			continue
		}

		multiBulkReply, ok := p.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("need multi bulk reply")
			continue
		}

		rep := handler.db.Exec(fakeConn, multiBulkReply.Args)
		if reply.IsErrorReply(rep) {
			logger.Error("exec err", rep.ToBytes())
		}
	}
}
