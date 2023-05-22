package main

import (
	"github.com/lagolibs/amqpfx"
	"go.uber.org/fx"
)

func main() {
	var m amqpfx.SessionManager

	fx.New(
		fx.NopLogger,
		amqpfx.NewSimpleModule("amq", "amqp://10.107.126.117:61616"),
		fx.Populate(m),
	)

}
