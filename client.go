package amqpfx

import (
	"context"
	"github.com/Azure/go-amqp"
	"time"
)

type clientConfig struct {
	connectionOption *amqp.ConnOptions

	uri            string
	connectTimeout time.Duration
}

func newClientConfig() clientConfig {
	return clientConfig{
		connectTimeout: 10 * time.Second,
	}
}

type Client struct {
	config *clientConfig
}

func newClient(config *clientConfig) (*Client, error) {
	sm := Client{
		config: config,
	}
	return &sm, nil
}

func (s *Client) Close() error {
	return nil
}

func (s *Client) NewConn() (*Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.connectTimeout)
	defer cancel()
	conn, err := amqp.Dial(ctx, s.config.uri, s.config.connectionOption)
	if err != nil {
		return nil, err
	}

	return &Conn{Conn: conn, ConnectTimeout: s.config.connectTimeout}, nil
}
