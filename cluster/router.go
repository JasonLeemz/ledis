package cluster

import "ledis/interface/resp"

func makeRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc // exists key
	routerMap["type"] = defaultFunc   // type key
	routerMap["set"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc

	routerMap["ping"] = ping
	routerMap["select"] = selectdb
	routerMap["rename"] = rename // rename k1 k2,可能会导致hash变化
	routerMap["renamenx"] = rename
	routerMap["flushdb"] = flushdb
	routerMap["del"] = del

	return routerMap
}

// get key ,set k v
func defaultFunc(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	// 找到转发给哪个节点
	key := string(cmdArgs[1]) // key
	// 一致性hash,获取节点ip
	peer := cluster.peerPicker.PickNode(key)
	return cluster.relay(peer, c, cmdArgs)
}
