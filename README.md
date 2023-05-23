# AMQP Fx

Fx Module for amqp client. This project provide a thin wrapper around https://github.com/Azure/go-amqp.

## Installation

Use Go modules to install amqpfx.

```shell
go get -u github.com/lagolibs/amqpfx
```

Example usage with multiple queue and viper:

## Usage

```yaml
amq.uris:
  clienta: amqp://localhost:61616
  clientb: amqp://localhost:61617
```

```go
package main

import (
	"github.com/lagolibs/amqpfx"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"os"
)

func init() {
	viper.AddConfigPath(lo.Must(os.Getwd()))
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	if err := viper.SafeWriteConfig(); err != nil {
		lo.Must0(viper.ReadInConfig())
	}
}

func main() {
	app := fx.New(
		amqpfx.NewModule("amq", amqpfx.WithURIs(viper.GetStringMapString("amq.uris"))),
		fx.Invoke(fx.Annotate(func(client *amqpfx.Client, client2 *amqpfx.Client) {}, fx.ParamTags(`name:"amq_clienta"`, `name:"amq_clienb"`))),
	)

	app.Run()
}
```

See [tests](amqpfx_test.go) for more usage.