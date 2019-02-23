package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
  // h "github.com/rodaine/hclencoder"
  // "log"
  "github.com/Rizbe/terraforming/src/gen"
)

type config struct {
	s3Info *[]s3Info `hcl:"resource aws_s3_bucket"`
}
type s3Info struct {
	name    string `hcl:",key"`
  bucket  string `hcl:"bucket"`
	policy  string `hcl:"policy"`
	version string `hcl:"version"`
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

//GetBucketACL get all bucket ACL
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

	buckets := s3Info{}
	// test := "building-price-ranges-export"
	allBuckets, _ := a.ListBuckets()

  policy, _ := a.GetBucketPolicy(&allBuckets[0])
  version, _ := a.GetBucketVersioning(&allBuckets[0])

  if policy == "" {
    buckets = s3Info{name: allBuckets[0], bucket: allBuckets[0]}

  } else {
    buckets = s3Info{name: allBuckets[0], bucket: allBuckets[0], policy: policy, version: version}

  }
  var allBucketsList  []s3Info
  allBucketsList = append(allBucketsList, buckets)


  input := config{s3Info: &allBucketsList}
  hcl, err := gen.GenHCL(input)
  if err != nil {
    fmt.Println(err)
  }

  fmt.Println(string(hcl))
  // hcl, err := h.Encode(input)
	// if err != nil {
	// 	log.Fatal("unable to encode: ", err)
	// }
  // fmt.Println(string(hcl))
  // fmt.Println(buckets)
  // hcl, _ := gen.GenHCL(buckets)
  // fmt.Println(hcl)

}
