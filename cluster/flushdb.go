package cluster

import (
	"ledis/interface/resp"
	"ledis/resp/reply"
)

func flushdb(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply {
	replies := cluster.broadcast(c, cmdArgs)
	var errReply reply.ErrorReply

	for _, re := range replies {
		if reply.IsErrorReply(re) {
			errReply = re.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return reply.MakeOKReply()
	}

	return reply.MakeStandardErrReply("error : " + errReply.Error())
}
