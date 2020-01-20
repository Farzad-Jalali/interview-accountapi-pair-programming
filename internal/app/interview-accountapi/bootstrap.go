package interview_accountapi

import (
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
)

// this file is needed as you are not allowed a package with only test files

func Configure() {
	settings.Configure()
	api.Configure()
}

func Start(ch <-chan bool, startedSignal chan bool) {
	api.StartServer(ch, startedSignal)
}
