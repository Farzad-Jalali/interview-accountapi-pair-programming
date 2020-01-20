package interview_accountapi

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type healthStage struct {
	t                   *testing.T
	error               error
	healthCheckResponse *http.Response
}

func HealthTest(t *testing.T) (*healthStage, *healthStage, *healthStage) {
	stage := &healthStage{
		t: t,
	}
	return stage, stage, stage
}

func (s *healthStage) the_service_is_healthy() *healthStage {
	return s
}

func (s *healthStage) a_health_check_request_is_made() *healthStage {
	s.healthCheckResponse, s.error = http.Get(fmt.Sprintf("%s/v1/health", viper.GetString(settings.ServiceName+"-address")))
	return s
}

func (s *healthStage) the_response_is_200_ok() *healthStage {
	assert.Nil(s.t, s.error)
	assert.Equal(s.t, 200, s.healthCheckResponse.StatusCode)
	return s
}
