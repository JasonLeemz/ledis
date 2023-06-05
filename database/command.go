package database

import "strings"

// 使用时都是只读场景，不存在并发问题
var cmdTable = make(map[string]*command)

type command struct {
	exector ExecFunc
	arity   int
}

func RegisterCommand(name string, exector ExecFunc, arity int) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		exector: exector,
		arity:   arity,
	}
}
