// Package queuex provides primitives for communicating with all queue
package queuex

import (
	oldmonkv1 "github.com/remotegarage/oldmonk/pkg/apis/oldmonk/v1"
	"github.com/streadway/amqp"
)

// RabbitmqController hold configration and client for rabbitmq queue
type RabbitmqController struct {
	Config  *oldmonkv1.ListOptions
	Channel *amqp.Channel
}

// NewRabbitmqClient create a new RabbitmqController object
// It returns the RabbitmqController
func NewRabbitmqClient(config *oldmonkv1.ListOptions) *RabbitmqController {
	conn, err := amqp.Dial(config.Uri)
	if err != nil {
		logger.Error("unable to dial", err)
		return &RabbitmqController{}
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Error("unable to create channel", err)
		return &RabbitmqController{}
	}
	rabbitmqController := RabbitmqController{
		Channel: ch,
		Config:  config,
	}
	return &rabbitmqController
}

// GetCount count the number of message in a queue
// It returns the number of Messages in a queue
func (r *RabbitmqController) GetCount() int32 {

	if err := r.Channel.ExchangeDeclare(
		r.Config.Exchange, // name of the exchange
		r.Config.Type,     // type
		true,              // durable
		false,             // delete when complete
		false,             // internal
		false,             // noWait
		nil,               // arguments
	); err != nil {
		return -1
	}

	err := r.Channel.Qos(
		10,    // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return -1
	}

	state, err := r.Channel.QueueDeclare(
		r.Config.Queue, // name of the queue
		true,           // durable
		false,          // delete when usused
		false,          // exclusive
		false,          // noWait
		nil,            // arguments
	)

	if err != nil {
		logger.Error("Queue Declare: ", err)
		return -1
	}

	return int32(state.Messages)
}

// Close will close  rabbitmq channel
// It returns the error
func (r *RabbitmqController) Close() error {
	return r.Channel.Close()
}
