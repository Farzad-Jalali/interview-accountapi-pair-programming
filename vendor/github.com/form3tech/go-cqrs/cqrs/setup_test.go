package cqrs

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/giantswarm/retry-go"

	"github.com/form3tech/go-cqrs/cqrs/support"
	"github.com/spf13/viper"

	"github.com/form3tech/go-docker-compose/dockercompose"
)

var stageRetryOptions = []retry.RetryOption{retry.MaxTries(400), retry.Timeout(10 * time.Second), retry.Sleep(100 * time.Millisecond)}

func TestMain(m *testing.M) {

	Must(os.Setenv("STACK_NAME", "local"))
	viper.Set("stack_name", "local")

	dynamicPorts, err := dockercompose.NewDynamicPorts("CONSUL_PORT:consul:8500", "SQS_PORT:sqs:1212")

	if err != nil {
		panic(err)
	}

	dir, _ := os.Getwd()

	dc, err := dockercompose.NewDockerCompose(
		dockercompose.NewAwsEcrAuth("288840537196", "eu-west-1"),
		filepath.Join(dir, "dockercompose/docker-compose.yml"),
		"cqrstesting",
		dynamicPorts,
		"wait_for")

	if err != nil {
		panic(err)
	}

	viper.Set("aws.default.sqs.endpoint", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("SQS_PORT")))

	containerWaiter := dockercompose.WaitForContainersToStart().
		ContainerLogLine("consul_1", "Consul agent running!").
		ContainerLogLine("sqs_1", "ElasticMQ server (0.13.8) started")

	err = dc.Start(containerWaiter)

	if err != nil {
		panic(err)
	}

	Must(os.Setenv("CONSUL_HTTP_ADDR", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("CONSUL_PORT"))))

	support.PurgeAllSQS()

	result := m.Run()

	if os.Getenv("STOP_DOCKER") != "" {
		dc.Stop()
	}

	os.Exit(result)

}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
