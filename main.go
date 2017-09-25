package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app       = kingpin.New("psst", "A command-line client for storing secrets in AWS Parameter Store.")
	region    = app.Flag("region", "AWS region").Envar("AWS_DEFAULT_REGION").Default("us-east-1").String()
	kmsAlias  = app.Flag("kms", "KMS key alias").Envar("AWS_KMS_ALIAS").Default("alias/aws/ssm").String()
	mfaSerial = app.Flag("mfa", "IAM MFA device ARN").Envar("AWS_MFA_ID").String()
	roleArn   = app.Flag("role", "IAM role ARN to assume").String()
)

func main() {
	app.HelpFlag.Short('h')
	app.Version("0.1.0")
	app.Author("Xavier Stevens <xavier.stevens@gmail.com>")

	configureReadCommand(app)
	configureWriteCommand(app)
	configureDeleteCommand(app)
	configureExecCommand(app)
	configureListCommand(app)
	configureHistoryCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}
