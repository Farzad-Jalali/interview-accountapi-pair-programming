package interview_accountapi

import (
	"database/sql"
	"fmt"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/externalmodels"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/form3tech/go-docker-compose/dockercompose"
	"github.com/form3tech/go-security/security"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/api/settings"
	"github.com/form3tech-oss/interview-accountapi-pair-programming/internal/app/interview-accountapi/log"
	"github.com/giantswarm/retry-go"
	"github.com/hashicorp/vault/api"
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
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

	dir, _ := os.Getwd()

	dynamicPorts, err := dockercompose.NewDynamicPorts(
		"POSTGRES_PORT:postgresql:5432",
		"VAULT_PORT:vault:8200",
	)

	if err != nil {
		panic(err)
	}

	ServerPort = getServerPort()
	settings.ServerPort = ServerPort

	var localAddress string
	if runtime.GOOS == "darwin" {
		localAddress = "docker.for.mac.host.internal"
	} else if os.Getenv("HOST_IP") != "" {
		localAddress = os.Getenv("HOST_IP")
	} else {
		localAddress = "localhost"
	}

	_ = os.Setenv("SQS_HOST", localAddress)

	dc, err := dockercompose.NewDockerCompose(
		dockercompose.NewAwsEcrAuth("288840537196", "eu-west-1"),
		filepath.Join(dir, "dockercompose/docker-compose.yml"),
		"interviewaccountapitesting",
		dynamicPorts)

	if err != nil {
		panic(err)
	}

	_ = os.Setenv("VAULT_TOKEN", "8fb95528-57c6-422e-9722-d2147bcba8ed")
	_ = os.Setenv("VAULT_ADDR", fmt.Sprintf("http://localhost:%d", dc.GetDynamicContainerPort("VAULT_PORT")))
	_ = os.Setenv("AWS_REGION", "eu-west-1")
	_ = os.Setenv("LOG_LEVEL", "debug")

	databaseUrl := fmt.Sprintf("postgres://root:password@localhost:%d?sslmode=disable", dc.GetDynamicContainerPort("POSTGRES_PORT"))
	viper.Set("DATABASE_URL", databaseUrl)
	viper.Set("database-host", "localhost")
	viper.Set("database-port", dc.GetDynamicContainerPort("POSTGRES_PORT"))
	viper.Set("database-ssl-mode", "disable")

	viper.Set("MessageVisibilityTimeout", 5)

	containerWaiter := dockercompose.WaitForContainersToStart().
		ContainerLogLine("postgresql_1", "database system is ready to accept connections").
		ContainerLogLine("vault_1", "core: post-unseal setup complete").
		Container("postgresql_1", func(finish chan error) {
			err := retry.Do(func() error {

				c, err := sql.Open("postgres", databaseUrl)
				if err != nil {
					return err
				}

				_, err = c.Exec("SELECT 1;")
				return err
			}, retry.MaxTries(500), retry.Sleep(time.Duration(500*time.Millisecond)))
			finish <- err
		})

	err = dc.Start(containerWaiter)

	if err != nil {
		panic(err)
	}

	connection, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec("DROP DATABASE IF EXISTS " + settings.ServiceName + ";")
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec("DROP USER IF EXISTS " + settings.ServiceName + "_user;")
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec("CREATE USER " + settings.ServiceName + "_user WITH PASSWORD '123';")
	if err != nil {
		panic(err)
	}

	_, err = connection.Exec("CREATE DATABASE " + settings.ServiceName + " OWNER " + settings.ServiceName + "_user;")
	if err != nil {
		panic(err)
	}

	testKeyPair, err = security.GenerateTestKeyPair()
	if err != nil {
		panic(err)
	}

	viper.Set("jwt-rsa-private-key", testKeyPair.PrivateKeyPem)
	viper.Set("jwt-public-key", testKeyPair.PublicKeyPem)
	viper.Set("jwt-public-der", testKeyPair.PublicKeyDer)

	log.Infof("vault token: %s", os.Getenv("VAULT_TOKEN"))
	vaultClient, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		panic(err)
	}

	secret, err := vaultClient.Logical().Read("/secret/application")

	if err != nil {
		panic(err)
	}

	if secret != nil {
		log.Infof("%v", secret)
	}

	_, err = vaultClient.Logical().Write("/secret/application", map[string]interface{}{
		"jwt-public-der": testKeyPair.PublicKeyDer,
	})

	if err != nil {
		panic(err)
	}

	viper.Set(settings.ServiceName+"-address", fmt.Sprintf("http://localhost:%d", ServerPort))

	Configure()

	startedSignal := make(chan bool)
	stopServer := make(chan bool, 1)

	go func() { Start(stopServer, startedSignal) }()
	<-startedSignal

	result := m.Run()

	if os.Getenv("STOP_DOCKER") != "" {
		dc.Stop()
	}

	stopServer <- true

	os.Exit(result)
}

func buildDefaultToken(organisations ...uuid.UUID) string {
	return buildTokenFor(organisations, map[string][]security.AuthoriseAction{
		string(models.ResourceTypeAccounts): AuthoriseAllActions,
	})
}

func buildTokenFor(organisations []uuid.UUID, accessPermissions map[string][]security.AuthoriseAction) string {
	jwtBuilder := security.NewJwtToken(testKeyPair.RsaPrivateKey, testUserId).
		ForOrganisations(organisations...)

	for recordType, permissions := range accessPermissions {
		jwtBuilder = jwtBuilder.GivePermissions(permissions...).ToRecordType(recordType)
	}
	token, err := jwtBuilder.Build()
	if err != nil {
		panic(err)
	}
	return token
}