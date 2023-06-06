package database

import "ledis/config"

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	mdb := &Database{}

	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}

	mdb.dbSet = make([]*DB, config.Properties.Databases)
	for i := range mdb.dbSet {
		singleDB := makeDB()
		singleDB.index = i
		mdb.dbSet[i] = singleDB
	}
	return mdb
}
