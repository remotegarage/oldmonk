// Package queuex provides primitives for communicating with all queue
package queuex

//
// import (
// 	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"
// 	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
// )
//
// // KafkaController hold configration and client for kafka topic
// type KafkaController struct {
// 	Config *oldmonkv1.ListOptions
// 	Client *kafka.Consumer
// }
//
// // NewKafkaClient create a new KafkaController object
// // It returns the KafkaController
// func NewKafkaClient(config *oldmonkv1.ListOptions) *KafkaController {
// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers":      config.Uri,
// 		"group.id":               config.Group,
// 		"session.timeout.ms":     6000,
// 		"auto.offset.reset":      "earliest",
// 		"statistics.interval.ms": 5000,
// 	})
// 	if err != nil {
// 		logger.Error("unable to Dial", err)
// 		return &KafkaController{}
// 	}
//
// 	kafkaController := KafkaController{
// 		Client: c,
// 		Config: config,
// 	}
// 	return &kafkaController
// }
//
// // GetCount count the number of message in a topic
// // It returns the number of Messages in a tube
// func (b *KafkaController) GetCount() int32 {
// 	return 45
// }
//
// // Close will close  kafka connection
// // It returns the error
// func (r *KafkaController) Close() error {
// 	return nil
// }
