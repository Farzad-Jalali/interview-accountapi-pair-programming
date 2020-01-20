package messaging

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/spf13/viper"
)

func newSqs() *sqs.SQS {
	awsSession := session.Must(session.NewSession())
	endpoint := viper.GetString("aws.default.sqs.endpoint")

	if endpoint == "" {
		return sqs.New(awsSession)
	}

	config := &aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String("eu-west-1"),
	}

	return sqs.New(awsSession, config)
}

func buildQueueUrl(s *sqs.SQS, stackName, queueName string) (*string, error) {
	queue := fmt.Sprintf("%s-%s", stackName, queueName)

	result, err := s.GetQueueUrl(&sqs.GetQueueUrlInput{QueueName: &queue})

	if err != nil {
		return nil, fmt.Errorf("could not get queue url for %s in region %s, error: %v", queue,
			aws.StringValue(s.Config.Region), err)
	}

	return result.QueueUrl, nil
}
