package amqpfx

import (
	"context"
	"fmt"
	"go.uber.org/fx"
)

// NewSimpleModule construct a module contain single client.
// Does not register group namespace.
// The name of the mongo client is the same as the name space.
func NewSimpleModule(namespace string, uri string) fx.Option {
	otp := newSessionManagerConfig()
	otp.URI = uri
	return fx.Module(namespace,
		fx.Provide(
			fx.Annotate(
				sessionManagerFactory(&otp),
				fx.ResultTags(
					fmt.Sprintf(`name:"%s"`, namespace),
				),
			),
		),
	)
}

// NewModule construct a new fx Module for mongodb, using configuration options
// Each mongo client will be named as <namespace>_<name>
// Also register a <namespace> group
func NewModule(namespace string, opts ...ModuleOption) fx.Option {
	conf := moduleConfig{}
	for i := range opts {
		opts[i](&conf)
	}
	return newModule(namespace, conf)
}

func newModule(namespace string, conf moduleConfig) fx.Option {
	configs := conf.configs
	if configs == nil || len(configs) == 0 {
		return fx.Module(namespace)
	}
	provides := make([]fx.Option, 0, len(configs))
	for name, option := range configs {
		provides = append(provides,
			fx.Provide(
				fx.Annotate(
					sessionManagerFactory(option),
					fx.ResultTags(
						fmt.Sprintf(`name:"%s_%s"`, namespace, name),
						fmt.Sprintf(`group:"%s"`, namespace),
					),
				),
			),
		)
	}
	return fx.Module(namespace, provides...)
}

func sessionManagerFactory(config *sessionManagerConfig) func(lc fx.Lifecycle) (SessionManager, error) {
	return func(lc fx.Lifecycle) (SessionManager, error) {
		sm, err := newSimpleSessionManager(config)
		if err != nil {
			return nil, err
		}
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return sm.Close()
			},
		})
		return sm, nil
	}
}

// ModuleOption applies an option to moduleConfig
type ModuleOption func(conf *moduleConfig)

// WithURIs create ModuleOption that parse a map of uris into moduleConfig.
// This help integrate with configuration library such as vipers
func WithURIs(uris map[string]string) ModuleOption {
	return func(conf *moduleConfig) {
		for key, uri := range uris {
			c := newSessionManagerConfig()
			c.URI = uri
			conf.configs[key] = &c
		}
	}
}

type moduleConfig struct {
	configs map[string]*sessionManagerConfig
}
