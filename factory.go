package amqpfx

import (
	"context"
	"github.com/Azure/go-amqp"
	"time"
)

type sessionManagerConfig struct {
	ConnectionOption *amqp.ConnOptions

	URI            string
	ConnectTimeout time.Duration
}

func newSessionManagerConfig() sessionManagerConfig {
	return sessionManagerConfig{
		ConnectTimeout: 10 * time.Second,
	}
}

type SessionManager interface {
	// NewConn ALWAYS create new connection
	// The connection created by this method is not managed, so you need to manually close it
	NewConn() (*amqp.Conn, error)

	// Run execute code on a managed session
	Run(func(session *amqp.Session) error) error

	// Close all managed connections
	Close() error
}

type simpleSessionManager struct {
	config *sessionManagerConfig
	conn   *amqp.Conn
}

func newSimpleSessionManager(config *sessionManagerConfig) (SessionManager, error) {
	sm := simpleSessionManager{
		config: config,
	}

	conn, err := sm.NewConn()
	if err != nil {
		return nil, err
	}
	sm.conn = conn
	return &sm, nil
}

func (s *simpleSessionManager) Run(executor func(session *amqp.Session) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ConnectTimeout)
	defer cancel()
	session, err := s.conn.NewSession(ctx, nil)
	if err != nil {
		return err
	}
	defer session.Close(context.Background())
	return executor(session)
}

func (s *simpleSessionManager) Close() error {
	return s.conn.Close()
}

func (s *simpleSessionManager) NewConn() (*amqp.Conn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.config.ConnectTimeout)
	defer cancel()
	return amqp.Dial(ctx, s.config.URI, s.config.ConnectionOption)
}

var _ SessionManager = (*simpleSessionManager)(nil)
