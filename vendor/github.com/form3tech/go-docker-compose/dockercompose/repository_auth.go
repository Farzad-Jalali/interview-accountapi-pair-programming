package dockercompose

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/google/uuid"
)

type RepositoryAuth interface {
	GenerateAuthConfig() (string, error)
}

type awsEcrAuth struct {
	accountId string
	region    string
}

func NewAwsEcrAuth(accountId, region string) *awsEcrAuth {
	return &awsEcrAuth{
		accountId: accountId,
		region:    region,
	}
}

type authTemplateData struct {
	AccountId string
	Region    string
	Token     string
}

func (a *awsEcrAuth) GenerateAuthConfig() (string, error) {

	session := session.Must(session.NewSession())

	params := &ecr.GetAuthorizationTokenInput{
		RegistryIds: []*string{aws.String(a.accountId)},
	}

	svc := ecr.New(session, aws.NewConfig().WithMaxRetries(10).WithRegion(a.region))

	resp, err := svc.GetAuthorizationToken(params)
	if err != nil {
		return "", fmt.Errorf("error authorizing: %s\n", err.Error())
	}

	if len(resp.AuthorizationData) != 1 {
		return "", fmt.Errorf("too many logins returned")
	}

	t, err := template.New("auth").Parse(configTmpl)

	if err != nil {
		return "", err
	}

	var config bytes.Buffer
	t.Execute(&config, authTemplateData{Token: *resp.AuthorizationData[0].AuthorizationToken, Region: a.region, AccountId: a.accountId})

	configDir, err := ioutil.TempDir("", uuid.New().String())

	if err != nil {
		panic(fmt.Errorf("could not create tmp dir for docker login, error: %v", err))
	}
	s := string(config.Bytes())
	fmt.Printf(s)

	err = ioutil.WriteFile(filepath.Join(configDir, "config.json"), config.Bytes(), os.ModePerm)

	if err != nil {
		return "", fmt.Errorf("could not write config.json file for docker compose, error: %v", err)
	}

	return configDir, nil

}

const configTmpl = `
{
  "auths": {
    "{{.AccountId}}.dkr.ecr.{{.Region}}.amazonaws.com": {
      "auth": "{{.Token}}",
      "email": "none"
    }
  }
}
`
