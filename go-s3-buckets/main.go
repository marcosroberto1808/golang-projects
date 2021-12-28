package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	s3session *s3.S3
)

func init()	{
	s3session = s3.New(
		session.Must(session.NewSession(&aws.Config {
		Region: aws.String("us-east-1"),
	})))
}

func listBuckets() ( resp *s3.ListBucketsOutput) {
	resp, err := s3session.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		panic(err)
	}
	return resp
}

func main()	{
	fmt.Println(listBuckets())
}
