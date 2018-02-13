package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	var mdlist []*cloudwatch.MetricDatum

	cliargs := &CliArgs{}

	app := cli.NewApp()
	app.Name = "Nagios go 2 cloudwatch"
	app.Version = "1.0.0"
	app.Usage = "Execute a nagios-style check and push check status and perfdata to cloudwatch metrics"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "namespace, n", Usage: "CloudWatch namespace"},
		cli.StringFlag{Name: "metric-name, m", Usage: "Metric name of check"},
		cli.StringFlag{Name: "timeout, t", Usage: "Timeout for command"},
		cli.StringFlag{Name: "dimensions, d", Usage: "CloudWatch metric dimensions `key=value,key=value`"},
		cli.StringFlag{Name: "region, r", Usage: "The AWS region", EnvVar: "AWS_REGION"},
	}
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Aaron Cossey",
		},
	}
	app.Action = func(c *cli.Context) error {
		//Region
		if !c.IsSet("region") {
			resp, err := http.Get("http://169.254.169.254/latest/meta-data/placement/availability-zone")
			if err != nil {
				log.Fatal("Unable to determine AWS region")
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Unable to determine AWS region")
			}
			cliargs.Region = string(body[:])
			defer resp.Body.Close()
		} else {
			cliargs.Region = c.String("region")
		}

		//Dimensions
		if !c.IsSet("dimensions") {
			cliargs.SetDimensions(c.String("dimensions"))
		}

		// Timeout
		if !c.IsSet("timeout") {
			cliargs.TimeOut = 0
		} else {
			cliargs.TimeOut = c.Int("timeout")
		}

		// Namespace
		if !c.IsSet("namespace") {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: Must provide namespace", -1)
		} else {
			cliargs.Namespace = c.String("namespace")
		}

		// Metric Name
		if !c.IsSet("metric-name") {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: Must provide metric name", -1)
		} else {
			cliargs.MetricName = c.String("metric-name")
		}

		//Command
		comm := c.Args()
		if len(comm) == 0 {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: No command was given to run", -1)
		} else {
			cliargs.Command = strings.Join(comm, " ")
		}

		//Create session
		mySession := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String(cliargs.Region)},
		}))
		svc := cloudwatch.New(mySession)

		//Run check command
		timestamp := time.Now()
		status, output := RunCommand(c.Args().First(), c.Args().Tail(), cliargs.TimeOut)

		// Process Status
		metricNameStatus := cliargs.MetricName + "-status"
		statusf := float64(status)

		md := cloudwatch.MetricDatum{
			MetricName: &metricNameStatus,
			Dimensions: cliargs.Dimensions,
			Timestamp:  &timestamp,
			Value:      &statusf,
		}
		mdlist = append(mdlist, &md)

		// Process Output
		processedOutput := ProcessOutput(output)
		for k := range processedOutput {
			metricNameOutput := cliargs.MetricName + "-" + k
			v := processedOutput[k]
			md := cloudwatch.MetricDatum{
				MetricName: &metricNameOutput,
				Dimensions: cliargs.Dimensions,
				Timestamp:  &timestamp,
				Value:      &v,
			}
			mdlist = append(mdlist, &md)
		}

		// PutMetricDataInput object
		pmdi := &cloudwatch.PutMetricDataInput{}
		pmdi.SetNamespace(cliargs.Namespace)
		pmdi.SetMetricData(mdlist)
		err := pmdi.Validate()
		if err != nil {
			log.Fatal(err)
		}

		// Call PutMetricData()
		_, err = svc.PutMetricData(pmdi)
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}

	app.Run(os.Args)
}
