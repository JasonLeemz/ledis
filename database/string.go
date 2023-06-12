package database

import (
	"ledis/interface/database"
	"ledis/interface/resp"
	"ledis/lib/utils"
	"ledis/resp/reply"
)

// GET key
func get(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}
	bytes, ok := entity.Data.([]byte)
	if !ok {
		return reply.MakeStandardErrReply("Data not a []byte type")
	}
	return reply.MakeBulkReply(bytes)

}

// SET key value
func set(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	entity := &database.DataEntity{
		Data: val,
	}
	db.PutEntity(key, entity)

	db.addAof(utils.ToCmdLine3("SET", args...))

	return reply.MakeOKReply()
}

// SETNX key value 如果key存在返回0，不存在返回1
func setnx(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	entity := &database.DataEntity{
		Data: val,
	}
	result := db.PutIfAbsent(key, entity)

	db.addAof(utils.ToCmdLine3("SETNX", args...))

	return reply.MakeIntReply(int64(result))
}

// GETSET k1 v1 ,相当于先执行set k1 v1,然后返回原始的value
func getset(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	val := args[1]

	oldEntity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}

	newEntity := &database.DataEntity{
		Data: val,
	}
	db.PutEntity(key, newEntity)

	db.addAof(utils.ToCmdLine3("GETSET", args...))

	return reply.MakeBulkReply(oldEntity.Data.([]byte))
}

// STRLEN key , 获取key对应value的长度
func strlen(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])

	entity, exist := db.GetEntity(key)
	if !exist {
		return reply.MakeNullBulkReply()
	}

	bytes := entity.Data.([]byte)
	return reply.MakeIntReply(int64(len(bytes)))
}

func init() {
	RegisterCommand("GET", get, 2)       // get key
	RegisterCommand("SET", set, 3)       // set k v
	RegisterCommand("SETNX", setnx, 3)   // setnx k v
	RegisterCommand("GETSET", getset, 3) // getset k v
	RegisterCommand("STRLEN", strlen, 2) // strlen k
}
