package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type deleteCommand struct {
	Name      string
	Value     string
	Overwrite bool
}

func configureDeleteCommand(app *kingpin.Application) {
	dc := &deleteCommand{}
	delete := app.Command("delete", "Delete secret from parameter store").Action(dc.runDelete)
	delete.Arg("name", "Secret name").StringVar(&dc.Name)
}

func (dc *deleteCommand) runDelete(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open session: %v\n", err.Error())
		return err
	}
	ssmClient := ssm.New(sess, config)

	// delete the secret in parameter store
	dpinput := &ssm.DeleteParameterInput{
		Name: aws.String(dc.Name),
	}
	_, err = ssmClient.DeleteParameter(dpinput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to delete secret: %v\n", err.Error())
	}

	return err
}
