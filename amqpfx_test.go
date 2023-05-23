package amqpfx_test

import (
	"github.com/lagolibs/amqpfx"
	"go.uber.org/fx"
)

func ExampleNewModule() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "amqp://localhost:61616"
	configs["clientb"] = "amqp://127.0.0.1:61616"

	fx.New(
		amqpfx.NewModule("amq", amqpfx.WithURIs(configs)),
		fx.Invoke(
			fx.Annotate(func(client *amqpfx.Client) {},
				fx.ParamTags(`name:"amq_clienta"`),
			),
		),
	).Run()
}

func ExampleNewModule_singleClient() {
	configs := make(map[string]string, 2)
	configs["clienta"] = "amqp://localhost:61616"

	fx.New(
		amqpfx.NewModule("amq", amqpfx.WithURIs(configs)),
		fx.Invoke(
			fx.Annotate(func(client *amqpfx.Client) {},
				fx.ParamTags(`name:"amq_clienta"`),
			),
		),
	).Run()
}

func ExampleNewSimpleModule() {
	fx.New(
		amqpfx.NewSimpleModule("amq", "amqp://localhost:61616"),
		fx.Invoke(
			fx.Annotate(func(client *amqpfx.Client) {},
				fx.ParamTags(`name:"amq"`),
			),
		),
	).Run()
}
