package main

import (
	"fmt"
	"text/tabwriter"
	"time"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type listCommand struct {
	Path string
}

func configureListCommand(app *kingpin.Application) {
	lc := &listCommand{}
	list := app.Command("list", "List all secrets under a path in parameter store").Action(lc.runList)
	list.Arg("path", "Path prefix to fetch secrets from").Required().StringVar(&lc.Path)
}

func writeMetadataTable(metadata []*ssm.ParameterMetadata) {
	writer := tabwriter.NewWriter(os.Stdout, 4, 0, 4, ' ', 0)
	fmt.Fprintln(writer, "Name\tKey Alias\tLast Modified\tUser\t")
	for _, paramMeta := range metadata {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n",
			*paramMeta.Name,
			*paramMeta.KeyId,
			paramMeta.LastModifiedDate.Format(time.RFC3339),
			*paramMeta.LastModifiedUser)
	}
	writer.Flush()
}

func (lc *listCommand) runList(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region).WithCredentialsChainVerboseErrors(true)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		return err
	}
	ssmClient := ssm.New(sess, config)

	// get parameter metadata; recursive path search
	paths := []*string{aws.String(lc.Path)}
	pathFilter := &ssm.ParameterStringFilter{
		Key:    aws.String("Path"),
		Option: aws.String("Recursive"),
		Values: paths,
	}
	filters := []*ssm.ParameterStringFilter{pathFilter}
	dpInput := &ssm.DescribeParametersInput{
		ParameterFilters: filters,
	}
	dpOutput, err := ssmClient.DescribeParameters(dpInput)

	writeMetadataTable(dpOutput.Parameters)

	return nil
}
