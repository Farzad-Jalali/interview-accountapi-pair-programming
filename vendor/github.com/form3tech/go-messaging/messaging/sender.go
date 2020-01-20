package messaging

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Sender interface {
	Send(destination string, message Message) error
}

type sqsSender struct {
	stackName string
	sqs       *sqs.SQS
}

type snsSender struct {
	stackName string
	sns       *sns.SNS
}

func CreateMessage(body interface{}) Message {
	messageAttributes := make(map[string]interface{})
	if attributes, ok := body.(MessageAttributes); ok {
		messageAttributes = attributes.GetMessageAttributes()
	}
	messageAttributes[TypeInfoAttribute] = reflect.TypeOf(body).String()
	return Message{Body: body, MessageAttributes: messageAttributes}
}

func NewSqsSender() Sender {
	return newSqsSender(viper.GetString("stack_name"))
}

func newSqsSender(stackName string) Sender {
	return sqsSender{
		stackName: stackName,
		sqs:       newSqs(),
	}
}

func (s sqsSender) Send(queueName string, message Message) error {

	queueUrl, err := buildQueueUrl(s.sqs, s.stackName, queueName)

	if err != nil {
		return err
	}

	jsonBytes, err := json.Marshal(message.Body)

	if err != nil {
		return err
	}

	messageBody := string(jsonBytes)
	if log.GetLevel() == log.DebugLevel {
		messageType := reflect.TypeOf(message.Body).String()
		log.Debugf("sending message type '%s' to queue '%s' with json payload '%s'", messageType, *queueUrl, messageBody)
	}
	_, err = s.sqs.SendMessage(&sqs.SendMessageInput{
		QueueUrl:          queueUrl,
		MessageBody:       &messageBody,
		MessageAttributes: s.getMessageAttributes(message),
	})

	return err
}

func NewSnsSender() Sender {
	return newSnsSender(viper.GetString("stack_name"))
}

func newSnsSender(stackName string) Sender {
	return snsSender{
		stackName: stackName,
		sns:       newSns(),
	}
}

func (s snsSender) Send(topicName string, message Message) error {

	topicArn, err := buildTopicArn(topicName, viper.GetString("stack_name"))

	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(message.Body)

	if err != nil {
		return err
	}

	messageBody := string(jsonBytes)

	log.Debugf("sending to %v json: %v", topicArn, messageBody)

	_, err = s.sns.Publish(&sns.PublishInput{
		Message:           aws.String(messageBody),
		TopicArn:          aws.String(topicArn),
		MessageAttributes: s.getMessageAttributes(message),
	})

	return err
}

func (snsSender) getMessageAttributes(message Message) map[string]*sns.MessageAttributeValue {
	stringDataType := "String"

	result := make(map[string]*sns.MessageAttributeValue)
	for k, v := range message.MessageAttributes {
		stringValue := fmt.Sprint(v)
		result[k] = &sns.MessageAttributeValue{
			DataType:    &stringDataType,
			StringValue: &stringValue,
		}
	}
	return result
}

func (sqsSender) getMessageAttributes(message Message) map[string]*sqs.MessageAttributeValue {
	stringDataType := "String"

	result := make(map[string]*sqs.MessageAttributeValue)
	for k, v := range message.MessageAttributes {
		stringValue := fmt.Sprint(v)
		result[k] = &sqs.MessageAttributeValue{
			DataType:    &stringDataType,
			StringValue: &stringValue,
		}
	}
	return result
}
