package amqpfx

import (
	"context"
	"github.com/Azure/go-amqp"
	"time"
)

type Conn struct {
	*amqp.Conn
	ConnectTimeout time.Duration
}

func (c *Conn) NewSession() (*Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.ConnectTimeout)
	defer cancel()
	session, err := c.Conn.NewSession(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Session{Session: session, ConnectTimeout: c.ConnectTimeout}, nil
}
