package queuex

import (
	"testing"
)

func TestRabbitmq(t *testing.T) {
	// config := &oldmonkv1.ListOptions{
	// 	Uri:      "amqp://zfwltqwc:1xuKkeEoO5nOE7F-K2Bh1nzCLcWH-C7q@shark.rmq.cloudamqp.com/zfwltqwc",
	// 	Queue:    "testing",
	// 	Type:     "direct",
	// 	Exchange: "testing",
	// 	Key:      "testing",
	// }
	// c := NewQueueConnection("RABBITMQ", config)
	// size := c.GetCount()
	// fmt.Println(size)
	// if size == 0 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", size, 10)
	// }
}

func TestSqs(t *testing.T) {
	// config := &oldmonkv1.ListOptions{
	// 	Uri:    "sqs://",
	// 	Region: "south-1",
	// }
	// c := NewQueueConnection("SQS", config)
	// size := c.GetCount()
	// fmt.Println(size)
	// if size == 0 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", size, 10)
	// }
}

func TestBeanstalk(t *testing.T) {
	// config := &oldmonkv1.ListOptions{
	// 	Uri:  "127.0.0.1:11300",
	// 	Tube: "default",
	// }
	// c := NewQueueConnection("BEANSTALK", config)
	// size := c.GetCount()
	// if true {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", size, 10)
	// }
}

func TestNats(t *testing.T) {
	// config := &oldmonkv1.ListOptions{
	// 	Uri:   "127.0.0.1:11300",
	// 	Queue: "default",
	// }
	// c := NewQueueConnection("NATS", config)
	// size := c.GetCount()
	// if true {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", size, 10)
	// }
}

func TestKafka(t *testing.T) {
	//  config := &oldmonkv1.ListOptions{
	//    Broker : "test",
	//    Group : "",
	//    Topic : "",
	//  }
	//  c := NewQueueConnection("KAFKA",config)
	//  size := c.GetCount()
	//  fmt.Println(size)
	// if size == 0 {
	//    t.Errorf("Sum was incorrect, got: %d, want: %d.", size, 10)
	// }
}
