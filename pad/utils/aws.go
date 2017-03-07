package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func GetAWSSession(region, profile string) (*session.Session, error) {
	var awsConfig *aws.Config

	if region != "" {
		awsConfig = &aws.Config{Region: aws.String(region)}
	} else {
		awsConfig = &aws.Config{}
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	})

	return sess, err
}
