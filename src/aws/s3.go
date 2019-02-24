package aws

import (
	"fmt"
	"time"

	"github.com/Rizbe/terraforming/src/gen"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type config struct {
	s3Info *[]s3Info `hcl:"resource aws_s3_bucket"`
}
type s3Info struct {
	name     string `hcl:",key"`
	bucket   string `hcl:"bucket"`
	policy   string `hcl:"policy"`
	version  string `hcl:"version"`
	corsRule s3Cors `hcl:"cors_rule"`
	// lifecycle_rule s3Lifecycle `hcl:"lifecycle_rule"`
}

type s3Cors struct {
	AllowedHeaders []*string `hcl:"AllowedHeaders"`
	AllowedMethods []*string `hcl:"AllowedHeaders"`
	AllowedOrigins []*string `hcl:"AllowedOrigins"`
	ExposeHeaders  []*string `hcl:"ExposeHeaders"`
	MaxAgeSeconds  *int64    `hcl:"MaxAgeSeconds"`
}

type s3Lifecycle struct {
	AbortIncompleteMultipartUpload s3AbortIncompleteMultipartUpload
	Expiration                     s3LifecycleExpiration
	Filter                         s3LifecycleRuleFilter
	ID                             *string
	Status                         *string
}

type s3AbortIncompleteMultipartUpload struct {
	DaysAfterInitiation *int64
}

type s3LifecycleExpiration struct {
	Date                      *time.Time
	Days                      *int64
	ExpiredObjectDeleteMarker *bool
}

type s3LifecycleRuleFilter struct {
	Prefix *string
	Tag    *Tag
}

type Tag struct {
	Key   *string
	Value *string
}

//ListBuckets all S4 buckets
func (a *ClientS3) ListBuckets() ([]string, error) {
	bucketList, err := a.Auth.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
	}
	s := make([]string, len(bucketList.Buckets))

	// fmt.Println(*bucketList.Buckets[0].Name)
	for i, name := range bucketList.Buckets {
		s[i] = *name.Name

	}

	return s, nil
}

//GetBucketPolicy get all bucket ACL
func (a *ClientS3) GetBucketPolicy(bucketName *string) (string, error) {
	bucketACL, err := a.Auth.GetBucketPolicy(&s3.GetBucketPolicyInput{Bucket: bucketName})
	if err != nil {
		// fmt.Println(err)
		return "", nil
	}

	return *bucketACL.Policy, nil

}

//GetBucketVersioning get versioning
func (a *ClientS3) GetBucketVersioning(bucketName *string) (string, error) {
	bucketVersion, err := a.Auth.GetBucketVersioning(&s3.GetBucketVersioningInput{Bucket: bucketName})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return "", nil
	}
	version := bucketVersion.Status
	return *version, nil

}

//GetBucketCors get all bucket ACL
func (a *ClientS3) GetBucketCors(bucketName *string) (s3Cors, error) {
	bucketCors, err := a.Auth.GetBucketCors(&s3.GetBucketCorsInput{Bucket: bucketName})
	if err != nil {
		fmt.Println(err)
		// return nil, nil
		fmt.Println(err)
	}
	t := s3Cors{}
	t = s3Cors{AllowedHeaders: bucketCors.CORSRules[0].AllowedHeaders, AllowedMethods: bucketCors.CORSRules[0].AllowedMethods, AllowedOrigins: bucketCors.CORSRules[0].AllowedOrigins, ExposeHeaders: bucketCors.CORSRules[0].ExposeHeaders, MaxAgeSeconds: bucketCors.CORSRules[0].MaxAgeSeconds}

	return t, nil

}

//GetBucketLifecycle get all bucket ACL
func (a *ClientS3) GetBucketLifecycle(bucketName *string) ([]*s3.LifecycleRule, error) {
	bucketLifecycle, err := a.Auth.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{Bucket: bucketName})
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	fmt.Println(bucketLifecycle.Rules[0].Expiration.Date)

	return bucketLifecycle.Rules, nil

}

//GetAllInfo Test
func (a *ClientS3) GetAllInfo() {
	var version, bucketName string

	buckets := s3Info{}
	// test := "building-price-ranges-export"
	allBuckets, _ := a.ListBuckets()

	for _, i := range allBuckets {
		bucketName = i
		fmt.Println(i)
		policy, _ := a.GetBucketPolicy(&i)
		version, _ = a.GetBucketVersioning(&bucketName)

		if policy == "" {
			buckets = s3Info{name: i}

		} else {
			buckets = s3Info{name: i, policy: policy, version: version}

		}
		fmt.Println(buckets)

	}

}

//GetAllInfo Test
func (a *ClientS3) TestAllInfo() {
	// var version, bucketName string
	var allBucketsList []s3Info
	buckets := s3Info{}
	// test := "building-price-ranges-export"
	allBuckets, _ := a.ListBuckets()

	policy, _ := a.GetBucketPolicy(&allBuckets[0])
	version, _ := a.GetBucketVersioning(&allBuckets[0])
	cors, _ := a.GetBucketCors(&allBuckets[0])
	// lifecycle_rule, _ := a.GetBucketLifecycle(&allBuckets[0])
	// l = s3Lifecycle{AbortIncompleteMultipartUpload: s3AbortIncompleteMultipartUpload{DaysAfterInitiation: lifecycle_rule[0].AbortIncompleteMultipartUpload.DaysAfterInitiation}, Expiration: s3LifecycleExpiration{Date: lifecycle_rule[0].Expiration.Date, Days: lifecycle_rule[0].Expiration.Days, ExpiredObjectDeleteMarker: lifecycle_rule[0].Expiration.ExpiredObjectDeleteMarker}, Filter: s3LifecycleRuleFilter{Prefix: lifecycle_rule[0].Filter.Prefix, Tag: Tag{} }}

	if (policy == "") || (len(cors.AllowedMethods) == 0) {
		buckets = s3Info{name: allBuckets[0], bucket: allBuckets[0]}

	} else {
		buckets = s3Info{name: allBuckets[0], bucket: allBuckets[0], policy: policy, version: version, corsRule: cors}

	}
	allBucketsList = append(allBucketsList, buckets)

	input := config{s3Info: &allBucketsList}
	hcl, err := gen.GenHCL(input)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(hcl))

	fmt.Println(&cors)
}
