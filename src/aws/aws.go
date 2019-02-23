package aws

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//NewSession create new aws session
func NewSession(region string) (*session.Session, error) {

	err := RegionCheck(region)
	if err != nil {
		return nil, err

	}

	config := &aws.Config{
		Credentials: credentials.NewEnvCredentials(),
		Region:      &region,
	}

	newSession := session.New(config)
	return newSession, nil
}

//RegionCheck Check to see if valid aws region
func RegionCheck(region string) error {
	regions := map[string]string{
		"us-east-2":      "US East (Ohio)",
		"us-east-1":      "US East (N. Virginia)",
		"us-west-1":      "US West (N. California)",
		"us-west-2":      "US West (Oregon)",
		"ap-northeast-1": "Asia Pacific (Tokyo)",
		"ap-northeast-2": "Asia Pacific (Seoul)",
		"ap-northeast-3": "Asia Pacific (Osaka-Local)",
		"ap-south-1":     "Asia Pacific (Mumbai)",
		"ap-southeast-1": "Asia Pacific (Singapore)",
		"ap-southeast-2": "Asia Pacific (Sydney)",
		"ca-central-1":   "Canada (Central)",
		"cn-north-1":     "China (Beijing)",
		"cn-northwest-1": "China (Ningxia)",
		"eu-central-1":   "EU (Frankfurt)",
		"eu-west-1":      "EU (Ireland)",
		"eu-west-2":      "EU (London)",
		"eu-west-3":      "EU (Paris)",
		"sa-east-1":      "South America (São Paulo)",
	}
	if val, ok := regions[region]; ok {
		fmt.Println(val)
		return nil
	}

	return errors.New("Invalid region")

}

//ClientS3 s3 client struct
type ClientS3 struct {
	Auth *s3.S3
}

//Initialize the client
func (a *ClientS3) Initialize(region string) {

	session, err := NewSession(region)
	if err != nil {
		log.Fatal(err)
	}
	a.Auth = s3.New(session)

}
