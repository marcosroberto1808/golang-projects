//robertom-aws-subnet-available-ip-count-check
package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-lambda-go/lambda"

)

const (
	REGION = "us-east-1"
	VPC_ID = "vpc-f7138f91"
	NAMESPACE = "VPC/Subnets"
)

var (
	objSession = session.Must(session.NewSession(&aws.Config {
		Region: aws.String(REGION),
	}))
	ec2Session *ec2.EC2
	cloudwatchSession *cloudwatch.CloudWatch
)

func init()	{
	ec2Session = ec2.New(objSession)
	cloudwatchSession = cloudwatch.New(objSession)
}

// Gets a slice of subnets from the API for the specified VPC
func getSubnetsInVPC(vpcid string, ec2client *ec2.EC2) ([]*ec2.Subnet, error) {
	subnetReq := ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("vpc-id"),
				Values: []*string{
					aws.String(vpcid),
				},
			},
		},
	}
	subnetResp, err := ec2client.DescribeSubnets(&subnetReq)
	if err != nil {
		return []*ec2.Subnet{}, err
	}
	return subnetResp.Subnets, nil
}

// Gets a map of subnets with AvailableIpAddressCount for the specified VPC ID
func getSubnetAvailableIpAddressCountInVPC(vpcid string, ec2client *ec2.EC2) (map[string]int64, error) {
	subnets, err := getSubnetsInVPC(vpcid, ec2client)
	if err != nil {
		return map[string]int64{}, err
	}

	tmpMap := make(map[string]int64)

	for _, subnet := range subnets {
		tmpMap[*subnet.SubnetId] = *subnet.AvailableIpAddressCount
	}

	return tmpMap, nil
}

// Generate the metric for each subnet in cloudwatch
func putCloudWatchMetrics(inputMap map[string]int64) {

	for key, element := range inputMap {
		_, err := cloudwatchSession.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(NAMESPACE),
			MetricData: []*cloudwatch.MetricDatum{
				&cloudwatch.MetricDatum{
					MetricName: aws.String("AvailableIpAddressCount"),
					Unit:       aws.String("Count"),
					Value:      aws.Float64(float64(element)),
					Dimensions: []*cloudwatch.Dimension{
						&cloudwatch.Dimension{
							Name:  aws.String("Subnet ID"),
							Value: aws.String(key),
						},
					},
				},
			},
		})
		if err != nil {
			fmt.Println("Error adding metrics:", err.Error())
			return
		}
	}
}

func Handler() {
	resultMap, err := getSubnetAvailableIpAddressCountInVPC(VPC_ID, ec2Session)
	fmt.Println(resultMap, err)
	putCloudWatchMetrics(resultMap)
	fmt.Println("Function invoked!")
}

func main()	{
	lambda.Start(Handler)
}
