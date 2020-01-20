package support

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type QueueReceiver struct {
	sqs       *sqs.SQS
	stackName string
}

func newSqs() *sqs.SQS {
	session := session.Must(session.NewSession())
	endpoint := viper.GetString("aws.default.sqs.endpoint")

	if endpoint == "" {
		return sqs.New(session)
	}

	config := &aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String("eu-west-1"),
	}

	return sqs.New(session, config)
}

func (r QueueReceiver) ListAll() []*string {
	result, err := r.sqs.ListQueues(&sqs.ListQueuesInput{})
	if err != nil {
		log.Warnf("fail to list all SQS Queues with error: %s", err)
		return []*string{}
	}
	return result.QueueUrls
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

func (r QueueReceiver) PurgeQueue(queueName string) error {
	queueUrl, err := buildQueueUrl(r.sqs, r.stackName, queueName)
	if err != nil {
		return err
	}

	_, err = r.sqs.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: queueUrl,
	})
	return err
}

func NewQueueReceiver() QueueReceiver {
	return QueueReceiver{
		sqs:       newSqs(),
		stackName: viper.GetString("stack_name"),
	}
}

func PurgeAllSQS() {
	queueReceiver := NewQueueReceiver()
	for _, queueUrl := range queueReceiver.ListAll() {
		_, err := queueReceiver.sqs.PurgeQueue(&sqs.PurgeQueueInput{
			QueueUrl: queueUrl,
		})
		if err != nil {
			log.Warnf("fail to purge SQS %s with error %s", *queueUrl, err)
		}
	}
}
