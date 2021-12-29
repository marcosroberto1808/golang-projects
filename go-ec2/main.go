package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	objSession = session.Must(session.NewSession(&aws.Config {
		Region: aws.String("us-east-1"),
	}))
	ec2Session *ec2.EC2
	cloudwatchSession *cloudwatch.CloudWatch
	vpcId = "vpc-f7138f91"
	namespace = "VPC/Subnets"
	resultmap map[string]int64
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

// Gets a map of AvailableIpAddressCount from subnetID
func getSubnetAvailableIpAddressCountInVPC(vpcid string, ec2client *ec2.EC2) (map[string]int64, error) {
	subnets, err := getSubnetsInVPC(vpcid, ec2client)
	if err != nil {
		return map[string]int64{}, err
	}

	tmpMap := make(map[string]int64)

	for _, subnet := range subnets {
		tmpMap[*subnet.SubnetId] = *subnet.AvailableIpAddressCount
	}

	// Generate the metrics
	putCloudWatchMetrics(tmpMap)

	return tmpMap, nil
}

// Generate the metric for each subnet in cloudwatch
func putCloudWatchMetrics(inputMap map[string]int64) {

	for key, element := range inputMap {
		_, err := cloudwatchSession.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
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


func main()	{
	//fmt.Println(listInstances())
	//fmt.Println(getSubnetsInVPC(vpcId, ec2Session))
	fmt.Println(getSubnetAvailableIpAddressCountInVPC(vpcId, ec2Session))

}
