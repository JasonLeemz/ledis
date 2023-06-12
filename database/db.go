package database

import (
	"ledis/datastruct/dict"
	"ledis/interface/database"
	"ledis/interface/resp"
	"ledis/resp/reply"
	"strings"
)

type DB struct {
	index int
	data  dict.Dict

	addAof func(CmdLine)
}

func makeDB() *DB {
	db := &DB{
		data:   dict.MakeSyncDict(),
		addAof: func(line CmdLine) {}, //空方法 防止报错
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
	if !validArity(cmd.arity, cmdLine) {
		return reply.MakeArgNumErrReply(cmdName)
	}

	fun := cmd.exector
	// SET K V
	return fun(db, cmdLine[1:])
}

func validArity(arity int, cmdArgs [][]byte) bool {
	argNum := len(cmdArgs)
	if arity >= 0 {
		return argNum == arity
	}
	return argNum >= -arity
}

/* ---- data Access ----- */
// Remove the given key from db
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exist := db.data.Get(key)
		if exist {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// GetEntity returns DataEntity bind to given key
func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	raw, exist := db.data.Get(key)
	if !exist {
		return nil, false
	}
	entity, _ := raw.(*database.DataEntity)
	return entity, true
}

// PutEntity a DataEntity into DB
func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}

// PutIfExists edit an existing DataEntity
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}

// PutIfAbsent insert an DataEntity only if the key not exists
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()

}
