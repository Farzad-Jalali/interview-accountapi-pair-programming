package settings

import (
	"os"
	"sync"
)

const (
	ApiName     = "interview-accountapi"
	UserID      = "9ef0183d-600f-415b-975e-2b722afc74f2"
	ServiceName = "interview_accountapi"
)

var (
	ServerPort              = 8080
	ApplicationClientId     string
	ApplicationClientSecret string
	StackName               string
	LogFormat               string
	LogLevel                string
)

var settingsOnce sync.Once

func Configure() {
	settingsOnce.Do(func() {
		StackName = GetStringOrDefault("STACK_NAME", "local")
		LogFormat = os.Getenv("LOG_FORMAT")
		LogLevel = os.Getenv("LOG_LEVEL")
	})
}

func GetStringOrDefault(envName, defaultVal string) string {
	if os.Getenv(envName) != "" {
		return os.Getenv(envName)
	}
	return defaultVal
}


