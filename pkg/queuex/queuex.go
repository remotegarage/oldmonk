// Package queuex provides primitives for communicating with all queue
package queuex

import (
	oldmonkv1 "github.com/remotegarage/oldmonk/api/v1"
	log "github.com/sirupsen/logrus"
)

var logger = log.WithFields(log.Fields{"queue_type": "Queue"})

// QueueFactory hold Interface for queue client
type QueueFactory interface {
	GetCount() int32
	Close() error
}

// NewQueueConnection create a new QueueFactory interface
// It returns the QueueFactory
func NewQueueConnection(queueType string, config *oldmonkv1.ListOptions) QueueFactory {
	switch queueType {
	case "SQS":
		clientInterface := NewSqsClient(config)
		return clientInterface
	case "RABBITMQ":
		clientInterface := NewRabbitmqClient(config)
		return clientInterface
	case "BEANSTALKD":
		clientInterface := NewBeanstalkClient(config)
		return clientInterface
	case "NATS":
		clientInterface := NewNatsClient(config)
		return clientInterface
	default:
		log.Printf("type undefined")
		return nil
	}
}
