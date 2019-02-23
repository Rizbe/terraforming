package main

import (
	"github.com/Rizbe/terraforming/src/aws"
)

func main() {
	s3session := aws.ClientS3{}
	s3session.Initialize("us-east-1")
	// allBuckets, _ := s3session.ListBuckets()
	// s3session.GetAllInfo()
	s3session.TestAllInfo()
	// fmt.Println(GetBucketACL)
}
