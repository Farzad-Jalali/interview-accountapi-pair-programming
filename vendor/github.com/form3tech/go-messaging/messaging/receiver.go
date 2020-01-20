package messaging

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DefaultVisibilityTimeout   = 1800
	DefaultWaitTimeSeconds     = 5
	DefaultMaxNumberOfMessages = 1
)

type ReceiveMessageResult struct {
	Output *sqs.ReceiveMessageOutput
	Error  error
}

type (
	Receiver interface {
		Close()
		Subscribe(queueName string, handler interface{}) error
		ReceiveOne(queueName string, out interface{}) error
		SubscribeMessage(queueName string, handler func(ctx context.Context, m Message) error) error
		ReceiveOneMessage(queueName string) (*Message, error)
		Purge(queueName string) error
	}

	sqsReceiver struct {
		log                 *log.Logger
		stackName           string
		maxNumberOfMessages int64
		visibilityTimeout   int64
		waitTimeSeconds     int64
		sqs                 *sqs.SQS
		quitChan            chan struct{}
		quitWaitGroup       *sync.WaitGroup
		receiveCtx          context.Context
		receiveCancel       context.CancelFunc
	}

	sqsReceiverBuilder struct {
		log                 *log.Logger
		stackName           string
		maxNumberOfMessages int64
		visibilityTimeout   int64
		waitTimeSeconds     int64
	}
)

func NewSqsReceiverBuilder() *sqsReceiverBuilder {
	defaultStackName := viper.GetString("stack_name")
	return &sqsReceiverBuilder{
		log:                 log.StandardLogger(),
		stackName:           defaultStackName,
		maxNumberOfMessages: DefaultMaxNumberOfMessages,
		visibilityTimeout:   DefaultVisibilityTimeout,
		waitTimeSeconds:     DefaultWaitTimeSeconds,
	}
}

func (s *sqsReceiverBuilder) WithLogger(log *log.Logger) *sqsReceiverBuilder {
	s.log = log
	return s
}

func (s *sqsReceiverBuilder) WithStackName(stackName string) *sqsReceiverBuilder {
	s.stackName = stackName
	return s
}

func (s *sqsReceiverBuilder) WithMaxNumberOfMessages(maxNumMessages int64) *sqsReceiverBuilder {
	s.maxNumberOfMessages = maxNumMessages
	return s
}

func (s *sqsReceiverBuilder) WithVisibilityTimeout(visibilityTimeout int64) *sqsReceiverBuilder {
	s.visibilityTimeout = visibilityTimeout
	return s
}

func (s *sqsReceiverBuilder) WithWaitTimeSeconds(waitTimeSeconds int64) *sqsReceiverBuilder {
	s.waitTimeSeconds = waitTimeSeconds
	return s
}

func (s *sqsReceiverBuilder) Build() Receiver {
	ctx, cancel := context.WithCancel(context.Background())
	return sqsReceiver{
		log:                 s.log,
		stackName:           s.stackName,
		maxNumberOfMessages: s.maxNumberOfMessages,
		visibilityTimeout:   s.visibilityTimeout,
		waitTimeSeconds:     s.waitTimeSeconds,
		sqs:                 newSqs(),
		quitChan:            make(chan struct{}),
		quitWaitGroup:       &sync.WaitGroup{},
		receiveCtx:          ctx,
		receiveCancel:       cancel,
	}
}

func NewSqsReceiver() Receiver {
	return NewSqsReceiverBuilder().Build()
}

func (r sqsReceiver) Close() {
	r.log.Debugf("closing receiver...")
	close(r.quitChan)
	r.receiveCancel()
	r.quitWaitGroup.Wait()
	r.log.Debugf("closed receiver...")
}

func (r sqsReceiver) Subscribe(queueName string, handler interface{}) error {
	h := func(ctx context.Context, m Message) error {
		var handlerArgs []reflect.Value
		messageInputIdx := 0

		handlerType := reflect.TypeOf(handler)

		if handlerType.NumIn() == 2 {
			handlerArgs = append(handlerArgs, reflect.ValueOf(&ctx))
			messageInputIdx = 1
		}

		inputValue := reflect.New(handlerType.In(messageInputIdx))
		input := inputValue.Interface()
		if err := json.Unmarshal([]byte(m.Body.(string)), &input); err != nil {
			log.Errorf("error unmarshalling json from sqs message, error: %s", err)
			return err
		}
		inputElement := inputValue.Elem()
		log.Debugf("received message with type '%s' from queue '%s' with json payload %s", inputElement.Type().String(), queueName, m.Body.(string))

		handlerArgs = append(handlerArgs, inputElement)
		handlerResult := reflect.ValueOf(handler).Call(handlerArgs)
		if handlerResult[0].Interface() != nil {
			return handlerResult[0].Interface().(error)
		}
		return nil
	}
	return r.SubscribeMessage(queueName, h)
}

func (r sqsReceiver) SubscribeMessage(queueName string, handler func(context context.Context, m Message) error) error {
	queue, err := buildQueueUrl(r.sqs, r.stackName, queueName)

	if err != nil {
		return err
	}

	r.quitWaitGroup.Add(1)

	subscription := func(queueUrl *string) {
		defer r.quitWaitGroup.Done()
		r.log.Infof("subscribing to %s", *queueUrl)
		for messageResult := range r.receiveMessages(queueUrl) {
			if messageResult.Error != nil {
				log.Errorf("failed to receive messages from SQS. Error: %s", messageResult.Error)
				continue
			}
			if messageResult.Output == nil || len(messageResult.Output.Messages) == 0 {
				continue
			}

			for _, m := range messageResult.Output.Messages {
				message := Message{
					Body:              *m.Body,
					MessageAttributes: r.readMessageAttributes(m.MessageAttributes),
				}
				err = handler(buildContextFromMessage(message), message)
				if err != nil {
					log.Error(err)
				} else {
					r.deleteMessage(queueUrl, m)
				}
			}
		}
	}

	go subscription(queue)

	return nil
}

func (sqsReceiver) readMessageAttributes(a map[string]*sqs.MessageAttributeValue) map[string]interface{} {
	stringDataType := "String"
	result := make(map[string]interface{})
	for k, v := range a {
		if *v.DataType == stringDataType {
			result[k] = *v.StringValue
		}
	}
	return result
}

func (r sqsReceiver) ReceiveOne(queueName string, out interface{}) error {
	message, err := r.ReceiveOneMessage(queueName)
	if err != nil {
		return err
	}

	if message == nil {
		return errors.New("no message received")
	}

	if err := json.Unmarshal([]byte(message.Body.(string)), out); err != nil {
		return err
	}

	return nil
}

func (r sqsReceiver) ReceiveOneMessage(queueName string) (*Message, error) {
	queueUrl, err := buildQueueUrl(r.sqs, r.stackName, queueName)
	if err != nil {
		return nil, err
	}

	messageOutput := <-r.receiveMessages(queueUrl)
	if messageOutput.Error != nil {
		return nil, err
	}
	if len(messageOutput.Output.Messages) > 0 {
		for _, message := range messageOutput.Output.Messages {
			r.deleteMessage(queueUrl, message)
			return &Message{
				Body:              *message.Body,
				MessageAttributes: r.readMessageAttributes(message.MessageAttributes),
			}, nil
		}
	}
	return nil, errors.New("no message received")
}

func (r sqsReceiver) receiveMessages(queueUrl *string) <-chan ReceiveMessageResult {

	receiveChan := make(chan ReceiveMessageResult, 1)

	go func() {
		defer close(receiveChan)
		for {
			output, err := r.sqs.ReceiveMessageWithContext(r.receiveCtx, &sqs.ReceiveMessageInput{
				AttributeNames: []*string{
					aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
				},
				MessageAttributeNames: []*string{
					aws.String(sqs.QueueAttributeNameAll),
				},
				QueueUrl:            queueUrl,
				MaxNumberOfMessages: aws.Int64(r.maxNumberOfMessages),
				VisibilityTimeout:   aws.Int64(r.visibilityTimeout),
				WaitTimeSeconds:     aws.Int64(r.waitTimeSeconds),
			})

			receiveChan <- ReceiveMessageResult{
				Output: output,
				Error:  err,
			}

			select {
			case <-r.quitChan:
				return
			default:
			}
		}
	}()

	return receiveChan
}

func (r sqsReceiver) deleteMessage(queueUrl *string, message *sqs.Message) {
	_, err := r.sqs.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueUrl,
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		log.Errorf("Could not delete message from sqs, error: %s", err)
	}
}

func buildContextFromMessage(message Message) context.Context {
	var correlationId string
	if message.MessageAttributes["correlation-id"] == nil || (message.MessageAttributes["correlation-id"]).(string) == "" {
		correlationId = strings.Replace(uuid.New().String(), "-", "", -1)
	} else {
		correlationId = (message.MessageAttributes["correlation-id"]).(string)
	}
	ctx := context.WithValue(context.Background(), "correlation-id", correlationId)
	return ctx
}

func (r sqsReceiver) Purge(queueName string) error {
	queueUrl, err := buildQueueUrl(r.sqs, r.stackName, queueName)
	if err != nil {
		return err
	}

	_, err = r.sqs.PurgeQueue(&sqs.PurgeQueueInput{
		QueueUrl: queueUrl,
	})
	if err != nil {
		return err
	}

	return nil
}
