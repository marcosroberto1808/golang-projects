package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var (
	//mySession = session.Must(session.NewSession())
	ec2Session *ec2.EC2
	vpcId = "vpc-f7138f91"
)

func init()	{
	ec2Session = ec2.New(
		session.Must(session.NewSession(&aws.Config {
			Region: aws.String("us-east-1"),
		})))
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

// Gets a map of subnet ID and CIDR address block from the specified VPC
func getSubnetNamesInVPC(vpcid string, ec2client *ec2.EC2) (map[string]string, error) {
	subnets, err := getSubnetsInVPC(vpcid, ec2client)
	if err != nil {
		return map[string]string{}, err
	}

	tmpMap := make(map[string]string)

	for _, subnet := range subnets {
		tmpMap[*subnet.SubnetId] = *subnet.CidrBlock
	}

	return tmpMap, nil
}


func main()	{
	//fmt.Println(listInstances())
	fmt.Println(getSubnetsInVPC(vpcId, ec2Session))
	//fmt.Println(getSubnetNamesInVPC(vpcId, ec2Session))
}
