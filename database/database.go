package database

import (
	"ledis/aof"
	"ledis/config"
	"ledis/interface/resp"
	"ledis/lib/logger"
	"ledis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet      []*DB // 数据库
	aofHandler *aof.AofHandler
}

func NewDatabase() *Database {
	database := &Database{}

	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	database.dbSet = make([]*DB, config.Properties.Databases)
	for i := range database.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		database.dbSet[i] = singleDB
	}

	// 初始化aof
	if config.Properties.AppendOnly {
		handler, err := aof.NewAofHandler(database)
		if err != nil {
			panic(err)
		}
		database.aofHandler = handler
		for _, db := range database.dbSet {
			sdb := db
			sdb.addAof = func(line CmdLine) {
				database.aofHandler.AddAof(sdb.index, line)
			}
		}
	}
	return database
}

func (database *Database) Exec(client resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.MakeArgNumErrReply(cmdName)
		}

		return execSelect(client, database, args[1:])
	}

	dbIndex := client.GetDBIndex()
	dbSet := database.dbSet[dbIndex]

	return dbSet.Exec(client, args)
}

// SELECT 1
func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {

	dbNum, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.MakeStandardErrReply("ERR invalid DB index")
	}
	if dbNum >= len(database.dbSet) {
		return reply.MakeStandardErrReply("ERR DB index is out of range")
	}

	c.SelectDB(dbNum)
	return reply.MakeOKReply()
}

func (database *Database) AfterClientClose(c resp.Connection) {
}

func (database *Database) Close() {
}
