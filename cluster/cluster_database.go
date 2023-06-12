package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"ledis/config"
	database2 "ledis/database"
	"ledis/interface/database"
	"ledis/interface/resp"
	"ledis/lib/consistenthash"
	"ledis/lib/logger"
	"ledis/resp/reply"
	"strings"
)

type ClusterDatabase struct {
	self string // 自己的名称地址

	nodes          []string // 整个集群的节点
	peerPicker     *consistenthash.NodeMap
	peerConnection map[string]*pool.ObjectPool
	db             database.Database
}

type CmdFunc func(cluster *ClusterDatabase, c resp.Connection, cmdArgs [][]byte) resp.Reply

var router = makeRouter()

func MakeClusterDatabase() *ClusterDatabase {
	clusterDB := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database2.NewStandAloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
		nodes:          nil,
	}

	nodes := make([]string, 0, len(config.Properties.Peers)+1) //len cap

	for _, peer := range config.Properties.Peers {
		nodes = append(nodes, peer)
	}

	nodes = append(nodes, config.Properties.Self)
	clusterDB.peerPicker.AddNode(nodes...)

	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		clusterDB.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{
			Peer: peer,
		})
	}

	clusterDB.nodes = nodes
	return clusterDB
}

func (cluster *ClusterDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {

	var result resp.Reply
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = reply.MakeUnknowErrReply()
			return
		}
	}()

	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		result = reply.MakeStandardErrReply("not supported cmd")
	}
	result = cmdFunc(cluster, client, args)

	return result
}

func (cluster *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	cluster.db.AfterClientClose(conn)
}

func (cluster *ClusterDatabase) Close() {
	cluster.db.Close()
}
