package queuex

import (
	"fmt"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/sqs"
	oldmonkv1 "github.com/evalsocket/oldmonk/pkg/apis/oldmonk/v1"
	"k8s.io/klog"
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
		return -1
	}
	return int32(count)
}

func (r *SqsController) getApproximateNumberOfMessages() (int32, error) {
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
		return 0, err
	}
	count, err := strconv.ParseInt(*resp.Attributes["ApproximateNumberOfMessages"], 10, 64)
	if err != nil {
		return 0, err
	}
	return int32(count), nil
}

func (r *SqsController) getAverageNumberOfMessagesSent() (float64, error) {
	period := int64(60)
	duration, err := time.ParseDuration("-5m")
	if err != nil {
		return 0.0, err
	}
	endTime := time.Now().Add(duration)
	startTime := endTime.Add(duration)

	query := &cloudwatch.MetricDataQuery{
		Id: aws.String("id1"),
		MetricStat: &cloudwatch.MetricStat{
			Metric: &cloudwatch.Metric{
				Namespace:  aws.String("AWS/SQS"),
				MetricName: aws.String("NumberOfMessagesSent"),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("QueueName"),
						Value: aws.String(path.Base(r.Config.Uri)),
					},
				},
			},
			Period: &period,
			Stat:   aws.String("Sum"),
		},
	}

	result, err := r.CloudwatchClient.GetMetricData(&cloudwatch.GetMetricDataInput{
		EndTime:           &endTime,
		StartTime:         &startTime,
		MetricDataQueries: []*cloudwatch.MetricDataQuery{query},
	})

	if err != nil {
		return 0.0, err
	}

	if len(result.MetricDataResults) > 1 {
		return 0.0, fmt.Errorf("Expecting cloudwatch metric to return single data point")
	}

	if result.MetricDataResults[0].Values != nil && len(result.MetricDataResults[0].Values) > 0 {
		var sum float64
		for i := 0; i < len(result.MetricDataResults[0].Values); i++ {
			sum += *result.MetricDataResults[0].Values[i]
		}
		return sum / float64(len(result.MetricDataResults[0].Values)), nil
	}

	klog.Errorf("NumberOfMessagesSent Cloudwatch API returned empty result for uri: %q", r.Config)

	return 0.0, nil
}

// Close will close  sqs connection
// It returns the error
func (r *SqsController) Close() error {
	return nil
}
