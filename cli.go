package main

import (
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"strings"
)

// CliArgs holds data passed in via CLI args
type CliArgs struct {
	Namespace  string
	MetricName string
	TimeOut    int
	Dimensions []*cloudwatch.Dimension
	Region     string
	Command    string
}

// SetDimensions appends a string of dimensions to the Dimensions slice in the CliArgs
func (ca *CliArgs) SetDimensions(d string) *CliArgs {
	for _, pair := range strings.Split(d, ",") {
		keyVal := strings.Split(pair, "=")
		if len(keyVal) == 1 {
			continue
		}
		dim := cloudwatch.Dimension{
			Name:  &keyVal[0],
			Value: &keyVal[1],
		}
		ca.Dimensions = append(ca.Dimensions, &dim)
	}
	return ca
}
