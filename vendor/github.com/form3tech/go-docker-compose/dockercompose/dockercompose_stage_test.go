package dockercompose

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
)

type dockerComposeStage struct {
	t             *testing.T
	dockerCompose DockerCompose
	error         error
	sqsPort       string
	sqsPort2      string
}

func DockerComposeTest(t *testing.T) (*dockerComposeStage, *dockerComposeStage, *dockerComposeStage) {

	stage := &dockerComposeStage{
		t: t,
	}

	return stage, stage, stage

}

func (s *dockerComposeStage) a_docker_compose_file_that_uses_a_sqs_container_located_in_a_private_ecr_repo() *dockerComposeStage {

	dynamicPorts, _ := NewDynamicPorts("SQS_PORT:sqs:1212", "SQS_PORT2:sqs2:1212")

	dir, _ := os.Getwd()
	var compose DockerCompose
	var err error
	compose, err = NewDockerCompose(NewAwsEcrAuth("288840537196", "eu-west-1"), filepath.Join(dir, "testassets/docker-compose.yml"), "dockercomposetesting", dynamicPorts)

	s.dockerCompose = compose
	s.error = err
	return s
}

func (s *dockerComposeStage) the_docker_compose_configuration_is_reloaded() *dockerComposeStage {

	return s.a_docker_compose_file_that_uses_a_sqs_container_located_in_a_private_ecr_repo()
}

func (s *dockerComposeStage) the_containers_are_started() *dockerComposeStage {

	containerWaiter := WaitForContainersToStart().
		ContainerLogLine("sqs_1", "ElasticMQ server (0.13.8) started")

	s.error = s.dockerCompose.Start(containerWaiter)

	return s
}

func (s *dockerComposeStage) there_is_no_error_when_starting_the_container() *dockerComposeStage {

	assert.Nil(s.t, s.error)

	return s
}

func (s *dockerComposeStage) and() *dockerComposeStage {

	return s
}

func (s *dockerComposeStage) the_sqs_container_is_running() *dockerComposeStage {
	endpoint := fmt.Sprintf("http://localhost:%d", s.dockerCompose.GetDynamicContainerPort("SQS_PORT"))

	session := session.Must(session.NewSession())

	config := &aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String("eu-west-1"),
	}

	sqs := sqs.New(session, config)

	queues, err := sqs.ListQueues(nil)

	assert.Nil(s.t, err)

	assert.Equal(s.t, len(queues.QueueUrls), 1)

	s.dockerCompose.Stop()

	return s
}

func (s *dockerComposeStage) the_dynamic_ports_are_cleared_from_the_environment() *dockerComposeStage {
	os.Unsetenv("SQS_PORT")
	os.Unsetenv("SQS_PORT2")
	return s
}
func (s *dockerComposeStage) the_dynamic_ports_are_stored() *dockerComposeStage {
	s.sqsPort = os.Getenv("SQS_PORT")
	s.sqsPort2 = os.Getenv("SQS_PORT2")
	return s
}

func (s *dockerComposeStage) new_ports_match_stored_ports() *dockerComposeStage {
	assert.Equal(s.t, s.sqsPort, os.Getenv("SQS_PORT"))
	assert.Equal(s.t, s.sqsPort2, os.Getenv("SQS_PORT2"))
	return s
}
