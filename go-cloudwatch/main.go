package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var (
	objSession = session.Must(session.NewSession(&aws.Config {
		Region: aws.String("us-east-1"),
	}))
	cloudwatchSession *cloudwatch.CloudWatch
)

func init()	{
	cloudwatchSession = cloudwatch.New(objSession)
}

func main() {

	_, err := cloudwatchSession.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("VPC/Subnets"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("AvailableIpAddressCount"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(15),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("Subnet ID"),
						Value: aws.String("subnet-9dd392c6"),
					},
				},
			},
			//&cloudwatch.MetricDatum{
			//	MetricName: aws.String("UniqueVisits"),
			//	Unit:       aws.String("Count"),
			//	Value:      aws.Float64(8628.0),
			//	Dimensions: []*cloudwatch.Dimension{
			//		&cloudwatch.Dimension{
			//			Name:  aws.String("SiteName"),
			//			Value: aws.String("example.com"),
			//		},
			//	},
			//},
			//&cloudwatch.MetricDatum{
			//	MetricName: aws.String("PageViews"),
			//	Unit:       aws.String("Count"),
			//	Value:      aws.Float64(18057.0),
			//	Dimensions: []*cloudwatch.Dimension{
			//		&cloudwatch.Dimension{
			//			Name:  aws.String("PageURL"),
			//			Value: aws.String("my-page.html"),
			//		},
			//	},
			//},
		},
	})
	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
		return
	}

	// Get information about metrics
	result, err := cloudwatchSession.ListMetrics(&cloudwatch.ListMetricsInput{
		Namespace: aws.String("VPC/Subnets"),
	})
	if err != nil {
		fmt.Println("Error getting metrics:", err.Error())
		return
	}

	for _, metric := range result.Metrics {
		fmt.Println(*metric.MetricName)

		for _, dim := range metric.Dimensions {
			fmt.Println(*dim.Name + ":", *dim.Value)
			fmt.Println()
		}
	}
}