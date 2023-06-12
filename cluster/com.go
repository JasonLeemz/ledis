package cluster

import (
	"context"
	"errors"
	"ledis/interface/resp"
	"ledis/lib/utils"
	"ledis/resp/client"
	"ledis/resp/reply"
	"strconv"
)

// communication

// 从连接池拿一个连接
func (cluster *ClusterDatabase) getPeerClient(peer string) (*client.Client, error) {

	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return nil, errors.New("connection not found")
	}

	ctx := context.Background()
	object, err := pool.BorrowObject(ctx)
	if err != nil {
		return nil, err
	}
	cli, ok := object.(*client.Client)
	if !ok {
		return nil, errors.New("wrong type")
	}
	return cli, nil
}

// 还回一个连接
func (cluster *ClusterDatabase) returnPeerClient(peer string, peerClient *client.Client) error {
	pool, ok := cluster.peerConnection[peer]
	if !ok {
		return errors.New("connection not found")
	}

	ctx := context.Background()
	return pool.ReturnObject(ctx, peerClient)
}

// 转发
func (cluster *ClusterDatabase) relay(peer string, conn resp.Connection, args [][]byte) resp.Reply {

	if peer == cluster.self {
		return cluster.db.Exec(conn, args)
	}

	peerClient, err := cluster.getPeerClient(peer)
	if err != nil {
		return reply.MakeStandardErrReply(err.Error())
	}

	defer func() {
		_ = cluster.returnPeerClient(peer, peerClient)
	}()

	// select db
	peerClient.Send(utils.ToCmdLine2("SELECT", strconv.Itoa(conn.GetDBIndex())))
	// other cmd
	return peerClient.Send(args)

}

// 广播 类似于执行 flushdb
func (cluster *ClusterDatabase) broadcast(conn resp.Connection, args [][]byte) map[string]resp.Reply {

	results := make(map[string]resp.Reply)
	for _, node := range cluster.nodes {
		result := cluster.relay(node, conn, args)
		results[node] = result
	}

	return results
}
