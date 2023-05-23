package amqpfx

import (
	"context"
	"github.com/Azure/go-amqp"
	"time"
)

type Session struct {
	*amqp.Session
	ConnectTimeout time.Duration
}

func (s *Session) NewReceiver(source string, opts ...ReceiverOption) (*Receiver, error) {
	var conf *amqp.ReceiverOptions = nil
	if len(opts) != 0 {
		conf = &amqp.ReceiverOptions{}
		for _, otp := range opts {
			otp(conf)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ConnectTimeout)
	defer cancel()

	r, err := s.Session.NewReceiver(ctx, source, conf)
	if err != nil {
		return nil, err
	}
	return &Receiver{Receiver: r}, nil
}

type Receiver struct {
	*amqp.Receiver
}

type ReceiveOption = func(conf *amqp.ReceiveOptions)

type ReceiverOption = func(conf *amqp.ReceiverOptions)

func (r *Receiver) Receive(ctx context.Context, opts ...ReceiveOption) (*amqp.Message, error) {
	var conf *amqp.ReceiveOptions = nil
	if len(opts) != 0 {
		conf = &amqp.ReceiveOptions{}
		for _, otp := range opts {
			otp(conf)
		}
	}
	return r.Receiver.Receive(ctx, conf)
}

func (s *Session) NewSender(source string, opts ...SenderOption) (*Sender, error) {
	var conf *amqp.SenderOptions = nil
	if len(opts) != 0 {
		conf = &amqp.SenderOptions{}
		for _, otp := range opts {
			otp(conf)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), s.ConnectTimeout)
	defer cancel()

	r, err := s.Session.NewSender(ctx, source, conf)
	if err != nil {
		return nil, err
	}
	return &Sender{Sender: r}, nil
}

type Sender struct {
	*amqp.Sender
}

type SendOption = func(conf *amqp.SendOptions)

type SenderOption = func(conf *amqp.SenderOptions)

func (s *Sender) Send(ctx context.Context, msg *amqp.Message, opts ...SendOption) error {
	var conf *amqp.SendOptions = nil
	if len(opts) != 0 {
		conf = &amqp.SendOptions{}
		for _, otp := range opts {
			otp(conf)
		}
	}
	return s.Sender.Send(ctx, msg, conf)
}
