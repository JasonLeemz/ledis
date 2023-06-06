package database

import (
	"ledis/interface/resp"
	"ledis/resp/reply"
)

// DEL
func del(db *DB, args [][]byte) resp.Reply {

	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}
	deleted := db.Removes(keys...)
	return reply.MakeIntReply(int64(deleted))
}

// EXISTS K1 K2 K3 ...
func exists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)

	for _, v := range args {
		key := string(v)
		_, exist := db.GetEntity(key)
		if exist {
			result++
		}
	}

	return reply.MakeIntReply(result)
}

//FLUSHDB

//TYPE

//RENAME

//RENAMENX

func init() {
	RegisterCommand("DEL", del, -2)       // -2代表最少2个参数
	RegisterCommand("EXISTS", exists, -2) // -2代表最少2个参数
}
