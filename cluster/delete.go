package cluster

import (
	"ledis/interface/resp"
	"ledis/resp/reply"
)

// del k1 k2 k3 k4 ...

func del(cluster *ClusterDatabase, c resp.Connection, args [][]byte) resp.Reply {
	replies := cluster.broadcast(c, args)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, re := range replies {
		if reply.IsErrorReply(re) {
			errReply = re.(reply.ErrorReply)
			break
		}
		intReply, ok := re.(*reply.IntReply)
		if !ok {
			errReply = reply.MakeStandardErrReply("error")
		}
		deleted += intReply.Code
	}

	if errReply == nil {
		return reply.MakeIntReply(deleted)
	}
	return reply.MakeStandardErrReply("error occurs: " + errReply.Error())
}
