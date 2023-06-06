package database

import (
	"ledis/interface/resp"
	"ledis/lib/wildcard"
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

// FLUSHDB
func flushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()

	return reply.MakeOKReply()
}

// TYPE k1
func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeStatusReply("none")
	}

	// TODO other type
	switch entity.Data.(type) {
	case []byte:
		return reply.MakeStatusReply("string")
	}

	return reply.MakeUnknowErrReply()
}

// RENAME k1 k2
func rename(db *DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])  // k1
	destKey := string(args[1]) // k2
	entity, exist := db.GetEntity(srcKey)
	if !exist {
		return reply.MakeStandardErrReply("no such key")
	}
	// 插入k2
	db.PutEntity(destKey, entity)
	// 删除k1
	db.Remove(srcKey)

	return reply.MakeOKReply()
}

// RENAMENX(检查k2是否会被覆盖)
func renameNX(db *DB, args [][]byte) resp.Reply {
	srcKey := string(args[0])  // k1
	destKey := string(args[1]) // k2

	// 检查k2是否存在
	_, ok := db.GetEntity(destKey)
	if ok {
		return reply.MakeIntReply(0)
	}

	// 检查k1
	entity, exist := db.GetEntity(srcKey)
	if !exist {
		return reply.MakeStandardErrReply("no such key")
	}
	// 插入k2
	db.PutEntity(destKey, entity)
	// 删除k1
	db.Remove(srcKey)

	return reply.MakeIntReply(1)
}

// KEYS *
func keys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0) // 存放所有的key

	db.data.ForEach(func(key string, val interface{}) bool {
		match := pattern.IsMatch(key)
		if match {
			result = append(result, []byte(key))
		}
		return true
	})

	return reply.MakeMultiBulkReply(result)
}

func init() {
	RegisterCommand("DEL", del, -2)          // -2代表最少2个参数
	RegisterCommand("EXISTS", exists, -2)    // 最少2个参数
	RegisterCommand("FLUSHDB", flushDB, -1)  // 最少1个参数 flushdb a b c
	RegisterCommand("TYPE", execType, 2)     // 参数个数必须是2 type k1
	RegisterCommand("RENAME", rename, 3)     // 参数个数必须是3 rename k1 k2
	RegisterCommand("RENAMENX", renameNX, 3) // 参数个数必须是3 rename k1 k2
	RegisterCommand("KEYS", keys, 2)         // 参数个数必须是2 keys *
}
