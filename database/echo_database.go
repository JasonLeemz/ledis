package database

import (
	"ledis/interface/resp"
	"ledis/resp/reply"
)

type EchoDatabase struct {
}

var theEchoDatabase = new(EchoDatabase)

func NewEchoDatabase() *EchoDatabase {
	return theEchoDatabase
}

func (e EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e EchoDatabase) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}

func (e EchoDatabase) Close() {
	//TODO implement me
	panic("implement me")
}
