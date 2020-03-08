package queuex

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"
)

// SqsController hold configration and client for sqs
type SqsController struct {
	Client *sqs.SQS
	Config *oldmonkv1.ListOptions
}

// NewSqsClient create a new SqsController object
// It returns the SqsController
func NewSqsClient(config *oldmonkv1.ListOptions) *SqsController {
	sess := session.New(&aws.Config{
		Region:     aws.String(config.Region),
		MaxRetries: aws.Int(5),
	})
	sqsController := SqsController{
		Config: config,
		Client: sqs.New(sess),
	}
	return &sqsController
}

// GetCount count the number of message in a queue
// It returns the number of Messages in a queue
func (r *SqsController) GetCount() int32 {
	attrib := "ApproximateNumberOfMessages"
	sendParams := &sqs.GetQueueAttributesInput{
		QueueUrl: aws.String(r.Config.Uri), // Required
		AttributeNames: []*string{
			&attrib, // Required
		},
	}
	resp, err := r.Client.GetQueueAttributes(sendParams)
	if err != nil {
		logger.Error("unable to get count", err)
		return -1
	}
	count, err := strconv.ParseInt(*resp.Attributes["ApproximateNumberOfMessages"], 10, 64)
	if err != nil {
		logger.Error("unable to parse count", err)
		return -1
	}
	return int32(count)
}

// Close will close  sqs connection
// It returns the error
func (r *SqsController) Close() error {
	return nil
}
