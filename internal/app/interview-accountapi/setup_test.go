package interview_accountapi

import (
	"fmt"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
	"github.com/form3tech/go-security/security"
	"github.com/google/uuid"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

var testKeyPair *security.TestKeyPair
var ServerPort = settings.ServerPort
var AuthoriseAllActions = []security.AuthoriseAction{security.CREATE, security.READ, security.EDIT, security.DELETE}
var testUserId = uuid.MustParse("b1850ea8-be26-4664-9cda-39464c19f39f")

func getServerPort() int {
	var serverPort int
	if _, err := os.Stat(".serverport"); os.IsNotExist(err) {
		serverPort, _ = freeport.GetFreePort()
		_ = ioutil.WriteFile(".serverport", []byte(fmt.Sprintf("%d", serverPort)), 0644)
	} else {
		portStr, _ := ioutil.ReadFile(".serverport")
		serverPort, _ = strconv.Atoi(string(portStr))
	}
	return serverPort
}

func TestMain(m *testing.M) {
	_ = os.Setenv("STACK_NAME", "local")


	ServerPort = getServerPort()
	settings.ServerPort = ServerPort

	viper.Set(settings.ServiceName+"-address", fmt.Sprintf("http://localhost:%d", ServerPort))

	Configure()

	startedSignal := make(chan bool)
	stopServer := make(chan bool, 1)

	go func() { Start(stopServer, startedSignal) }()
	<-startedSignal

	result := m.Run()

	stopServer <- true

	os.Exit(result)
}