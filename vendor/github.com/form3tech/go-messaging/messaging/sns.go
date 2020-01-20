package messaging

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/viper"
)

func newSns() *sns.SNS {
	awsSession := session.Must(session.NewSession())
	endpoint := viper.GetString("aws.default.sns.endpoint")

	if endpoint == "" {
		return sns.New(awsSession)
	}

	config := &aws.Config{
		Endpoint: &endpoint,
		Region:   aws.String("eu-west-1"),
	}

	return sns.New(awsSession, config)
}

func buildTopicArn(topicName string, stackName string) (string, error) {
	stsClient := newSts()
	input := &sts.GetCallerIdentityInput{}

	result, err := stsClient.GetCallerIdentity(input)
	if err != nil {
		return "", err
	}

	region := "eu-west-1"
	if stsClient.Config.Region != nil && *stsClient.Config.Region != "" {
		region = *stsClient.Config.Region
	}

	return fmt.Sprintf("arn:aws:sns:%v:%v:%v-%v", region, *result.Account, stackName, topicName), nil
}
