package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"ledis/resp/client"
)

type connectionFactory struct {
	Peer string
}

func (c *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	cli, err := client.MakeClient(c.Peer)
	if err != nil {
		return nil, err
	}

	cli.Start()

	p := pool.NewPooledObject(cli)
	return p, nil
}

func (c *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	cli, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}

	cli.Close()
	return nil
}

func (c *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (c *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
