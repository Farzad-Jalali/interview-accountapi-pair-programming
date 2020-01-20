package messaging

import (
	"errors"
	"testing"
	"time"

	"fmt"
	"os"
	"path/filepath"

	"context"

	"runtime"

	"github.com/form3tech/go-docker-compose/dockercompose"
	retry "github.com/giantswarm/retry-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type SendAndReceiveStage struct {
	t               *testing.T
	expectedMessage TestMessage
	actualMessage   TestMessage
	context         *context.Context
	error           error
	receiver        Receiver
}

type TestMessage struct {
	Payload string
}

const (
	QueueName  = "build-completed-queue"
	QueueName2 = "build-completed-queue2"
	TopicName  = "form3-events"
)

var dc dockercompose.DockerCompose

func SetUp() {
	viper.Set("stack_name", "local")
	dir, _ := os.Getwd()

	logrus.SetLevel(logrus.DebugLevel)

	dynamicPorts, err := dockercompose.NewDynamicPorts("SQS_PORT:sqs:1212", "SNS_PORT:sns:9911")

	if err != nil {
		panic(err)
	}

	var localAddress string
	if runtime.GOOS == "darwin" {
		localAddress = "docker.for.mac.host.internal"
	} else if os.Getenv("HOST_IP") != "" {
		localAddress = os.Getenv("HOST_IP")
	} else {
		localAddress = "localhost"
	}
	Must(os.Setenv("SQS_HOST", localAddress))

	dc, err = dockercompose.NewDockerCompose(
		dockercompose.NewAwsEcrAuth("288840537196", "eu-west-1"),
		filepath.Join(dir, "docker/docker-compose.yml"),
		"messagingtesting",
		dynamicPorts,
		"wait_for")
	if err != nil {
		panic(err)
	}

	containerWaiter := dockercompose.WaitForContainersToStartWithTimeout(5*time.Minute).
		ContainerLogLine("sqs_1", "ElasticMQ server (0.13.8) started").
		ContainerLogLine("sns_1", "Apache Camel 0.1.3 (CamelContext: sns) started")
	Must(dc.Start(containerWaiter))

	viper.Set("aws.default.sqs.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("SQS_PORT")))
	viper.Set("aws.default.sns.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("SNS_PORT")))
}

func TearDown() {
	dc.Stop()
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func SendAndReceiveTest(t *testing.T) (*SendAndReceiveStage, *SendAndReceiveStage, *SendAndReceiveStage) {
	stage := SendAndReceiveStage{
		t:        t,
		receiver: NewSqsReceiver(),
	}

	return &stage, &stage, &stage
}

func (s *SendAndReceiveStage) Teardown() *SendAndReceiveStage {
	s.receiver.Close()
	return s
}

func (s *SendAndReceiveStage) And() *SendAndReceiveStage {
	return s
}

func (s *SendAndReceiveStage) A_message_with_the_following_payload(message string) *SendAndReceiveStage {
	s.expectedMessage = TestMessage{
		Payload: message,
	}
	return s
}

func (s *SendAndReceiveStage) A_listener_for_queue_(queueName string) *SendAndReceiveStage {
	if err := s.receiver.Subscribe(queueName, s.messageReceived); err != nil {
		s.t.Fatalf("could not subscribe to queue: error: %v", err)
	}
	return s
}

func (s *SendAndReceiveStage) The_message_is_sent_to_queue_(queueName string) *SendAndReceiveStage {
	sender := NewSqsSender()
	err := sender.Send(queueName, CreateMessage(s.expectedMessage))

	assert.Nil(s.t, err)
	return s
}

func (s *SendAndReceiveStage) The_message_is_received_by_the_listener() *SendAndReceiveStage {
	err := retry.Do(func() error {
		if s.actualMessage == s.expectedMessage {
			return nil
		} else {
			return errors.New("the message was not received")
		}
	}, retry.Sleep(500*time.Millisecond), retry.Timeout(60*time.Second), retry.MaxTries(20000))

	assert.Nil(s.t, err)

	return s
}

func (s *SendAndReceiveStage) messageReceived(message TestMessage) error {
	s.actualMessage = message
	return nil
}

func (s *SendAndReceiveStage) messageReceivedWithContext(ctx *context.Context, message TestMessage) error {
	s.actualMessage = message
	s.context = ctx
	return nil
}

func (s *SendAndReceiveStage) A_listener_with_context_for_queue_(queueName string) *SendAndReceiveStage {
	receiver := NewSqsReceiver()
	if err := receiver.Subscribe(queueName, s.messageReceivedWithContext); err != nil {
		s.t.Fatalf("could not subscribe to queue, error: %v", err)
	}
	return s
}

func (s *SendAndReceiveStage) The_context_must_contains_correlationId() *SendAndReceiveStage {
	if s.context == nil {
		s.t.Fatalf("expected context not to be nil")
	} else {
		assert.NotNil(s.t, (*s.context).Value("correlation-id"))
		assert.NotEqual(s.t, "", (*s.context).Value("correlation-id"))
	}
	return s
}

func (s *SendAndReceiveStage) A_notification_is_sent() *SendAndReceiveStage {
	sender := NewSnsSender()
	s.error = sender.Send(TopicName, CreateMessage(s.expectedMessage))
	return s
}

func (s *SendAndReceiveStage) The_notification_is_sent_succesfully() *SendAndReceiveStage {
	assert.Nil(s.t, s.error)
	return s
}

func (s *SendAndReceiveStage) The_queue_is_purged(queueName string) *SendAndReceiveStage {
	s.receiver.Purge(queueName)
	return s
}

func (s *SendAndReceiveStage) The_queue_is_empty(queueName string) *SendAndReceiveStage {
	_, err := s.receiver.ReceiveOneMessage(queueName)
	assert.NotNil(s.t, err)
	assert.Equal(s.t, "no message received", string(err.Error()))
	return s
}
