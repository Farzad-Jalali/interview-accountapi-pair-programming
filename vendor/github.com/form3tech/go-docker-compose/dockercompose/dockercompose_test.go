package dockercompose

import "testing"

func Test_DockerComposeWithAWSAuth(t *testing.T) {

	given, when, then := DockerComposeTest(t)

	given.
		a_docker_compose_file_that_uses_a_sqs_container_located_in_a_private_ecr_repo()

	when.
		the_containers_are_started()

	then.
		there_is_no_error_when_starting_the_container().and().
		the_sqs_container_is_running()

}

func Test_DynamicPortsAreMappedCorrectlyForRunningContainers(t *testing.T) {

	given, when, then := DockerComposeTest(t)

	given.
		a_docker_compose_file_that_uses_a_sqs_container_located_in_a_private_ecr_repo()

	when.
		the_containers_are_started().and().
		the_dynamic_ports_are_stored().and().
		the_dynamic_ports_are_cleared_from_the_environment().and().
		the_docker_compose_configuration_is_reloaded()

	then.
		new_ports_match_stored_ports()
}
