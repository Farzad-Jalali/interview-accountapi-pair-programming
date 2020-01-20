package messaging

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/viper"
)

func newSts() *sts.STS {
	awsSession := session.Must(session.NewSession())
	endpoint := viper.GetString("aws.default.sts.endpoint")

	if endpoint == "" {
		return sts.New(awsSession)
	}

	config := &aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String("eu-west-1"),
	}

	return sts.New(awsSession, config)
}
