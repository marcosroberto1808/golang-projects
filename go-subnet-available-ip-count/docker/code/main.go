package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Environment Variables
	REGION = os.Getenv("REGION")
	VPC_ID = os.Getenv("VPC_ID")

	// Session Object for AWS Connection
	objSession = session.Must(session.NewSession(&aws.Config {
		Region: aws.String(REGION),
	}))

	// Session Type variable
	ec2Session *ec2.EC2

)

// Collector Type Struct
type newSubnetCollector struct {
	subnetAvailable *prometheus.Desc
}

// Collector Constructor Method
func newSubnetAvailableCollector() *newSubnetCollector {
	return &newSubnetCollector{
		subnetAvailable: prometheus.NewDesc("subnet_available_ip_count",
			"Shows the amount of ip addresses available in the subnet.",
			[]string{"subnet_id"}, nil,
		),
	}
}

// Collector Constructor Method Describe
func (collector *newSubnetCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.subnetAvailable
}

// Collector Constructor Method Collect
func (collector *newSubnetCollector) Collect(ch chan<- prometheus.Metric) {

	// Call AWS API and save in resultMAP
	resultMap, err := getSubnetAvailableIpAddressCountInVPC(VPC_ID, ec2Session)
	fmt.Println(resultMap, err)
	fmt.Println("Last Execution Time:", time.Now().Format("2022-01-02 15:04:05"))

	// Write values from resultMap to the metric channel
	for subnet, value := range resultMap {
		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		//ch <- prometheus.MustNewConstMetric(collector.subnetAvailable, prometheus.CounterValue, metricValue)
		ch <- prometheus.MustNewConstMetric(collector.subnetAvailable, prometheus.GaugeValue, float64(value), subnet)
	}
}

// Init variables
func init()	{
	ec2Session = ec2.New(objSession)

}

// Gets a slice of subnets from the AWS API for the specified VPC
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

func main()	{

	// Start the Custom Collector
	getNewData := newSubnetAvailableCollector()
	// Register the Collector
	prometheus.MustRegister(getNewData)
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9000", nil))

}
