package queuex

import (
	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"
	nats "github.com/nats-io/nats.go"
)

// NatsController hold configration and client for sqs
type NatsController struct {
	Client *nats.Conn
	Config *oldmonkv1.ListOptions
}

// NewNatsClient create a new NatsController object
// It returns the NatsController
func NewNatsClient(config *oldmonkv1.ListOptions) *NatsController {
	nc, err := nats.Connect(config.Uri)
	if err != nil {
		logger.Error("unable to dial", err)
		return &NatsController{}
	}
	natsController := NatsController{
		Config: config,
		Client: nc,
	}
	return &natsController
}

// GetCount count the number of message in a queue
// It returns the number of Messages in a queue
func (r *NatsController) GetCount() int32 {
	sub, err := r.Client.SubscribeSync(r.Config.Queue)
	if err != nil {
		return -1
	}
	count, _, err := sub.Pending()
	if err != nil {
		return -1
	}
	return int32(count)
}

// Close will close  sqs connection
// It returns the error
func (r *NatsController) Close() error {
	return nil
}
