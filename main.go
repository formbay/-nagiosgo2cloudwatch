package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	var region string
	var dimensions string
	var command string
	var timeout int

	app := cli.NewApp()
	app.Name = "Nagios go 2 cloudwatch"
	app.Version = "1.0.0"
	app.Usage = "Execute a nagios-style check and push check status and perfdata to cloudwatch metrics"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "namespace, n", Usage: "CloudWatch namespace"},
		cli.StringFlag{Name: "base, b", Usage: "Base name for checks"},
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
			region = string(body[:])
			defer resp.Body.Close()
		} else {
			region = c.String("region")
		}

		//Dimensions
		if !c.IsSet("dimensions") {
			dimensions = ""
		} else {
			dimensions = c.String("dimensions")
		}

		// Timeout
		if !c.IsSet("timeout") {
			timeout = 0
		} else {
			timeout = c.Int("timeout")
		}

		//Command
		comm := c.Args()
		if len(comm) == 0 {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: No command was given to run", -1)
		} else {
			command = strings.Join(comm, " ")
		}

		//Create session
		mySession := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            aws.Config{Region: aws.String(region)},
		}))
		svc := cloudwatch.New(mySession)
		fmt.Println(svc)

		//Create CW object
		cw := NewCW(c.String("namespace"), c.String("base"), dimensions)

		//Run check command
		status, output := RunCommand(c.Args().First(), c.Args().Tail(), timeout)

		//Add status to CW object
		fmt.Println(cw.AddData("Status", float64(status)))
		ProcessOutput(output)
		for k, v := range ProcessOutput(output) {
			fmt.Println(cw.AddData(k, v))
		}
		return nil
	}

	app.Run(os.Args)
}
