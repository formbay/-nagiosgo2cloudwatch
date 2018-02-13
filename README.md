# Nagios go to cloudwatch - nagiosgo2cloudwatch

## Description

Go wrapper for nagios style checks which sends the check result to cloudwatch metrics

Inspired by [n2cw](https://github.com/slank/n2cw).

## Usage

```
NAME:
   Nagios go 2 cloudwatch - Execute a nagios-style check and push check status and perfdata to cloudwatch metrics

USAGE:
   nagiosgo2cloudwatch [global options] command [command options] [arguments...]

VERSION:
   1.0.0

AUTHOR:
   Aaron Cossey

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --namespace value, -n value                               CloudWatch namespace
   --metric-name value, -m value                             Metric name of check
   --timeout value, -t value                                 Timeout for command
   --dimensions key=value,key=value, -d key=value,key=value  CloudWatch metric dimensions key=value,key=value
   --region value, -r value                                  The AWS region [$AWS_REGION]
   --help, -h                                                show help
   --version, -v                                             print the version
```

### Example

```
./nagiosgo2cloudwatch -n "NagiosChecks" -m check_disk -d "instance-id=i-hurdur,env=lolcat,instance-name=beavis" -r ap-southeast-2 /usr/lib64/nagios/plugins/check_disk -w 10% -c 5%
```

## Building

Make sure your $GOPATH is set: https://github.com/golang/go/wiki/GOPATH

Grab the external dependencies:
```
go get gopkg.in/urfave/cli.v1
go get -u github.com/aws/aws-sdk-go/...
```
Build
```
go build
```
## IAM Authentication

Your user or role needs to have at least `cloudwatch:PutMetricData`
