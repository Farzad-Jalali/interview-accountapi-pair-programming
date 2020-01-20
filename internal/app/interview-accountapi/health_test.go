package interview_accountapi

import (
	"testing"
)

func TestAcc_GetHealth(t *testing.T) {
	given, when, then := HealthTest(t)

	given.
		the_service_is_healthy()
	when.
		a_health_check_request_is_made()
	then.
		the_response_is_200_ok()
}
