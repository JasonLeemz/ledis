package cluster

import (
	"ledis/interface/resp"
	"ledis/resp/reply"
)

// rename k1 k2
func rename(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.MakeStandardErrReply("ERR Wrong number args")
	}

	k1 := string(cmdArgs[1])
	k2 := string(cmdArgs[2])
	peer1 := cluster.peerPicker.PickNode(k1)
	peer2 := cluster.peerPicker.PickNode(k2)
	if peer1 != peer2 {
		return reply.MakeStandardErrReply("ERR rename must within one peer")
	}
	return cluster.relay(peer1, c, cmdArgs)
}
