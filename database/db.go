package database

import (
	"ledis/datastruct/dict"
	"ledis/interface/resp"
	"ledis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict
}

func MakeDB() *DB {
	db := &DB{
		data: dict.MakeSyncDict(),
	}
	return db
}

type ExecFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine = [][]byte

func (db *DB) Exec(c resp.Connection, cmdLine CmdLine) resp.Reply {
	// PING SET SETNX ...
	cmdName := strings.ToLower(string(cmdLine[0]))

	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.MakeStandardErrReply("Err unknown command: " + cmdName)
	}

	// 参数个数
	if validArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}

	fun := cmd.exector
	// SET K V
	return fun(db, cmdLine[1:])
}

func validArity(arity int, args [][]byte) bool {
	return true
}
