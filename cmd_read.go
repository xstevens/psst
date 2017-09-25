package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type readCommand struct {
	Name    string
	Decrypt bool
}

func configureReadCommand(app *kingpin.Application) {
	rc := &readCommand{}
	read := app.Command("read", "Read secret from parameter store").Action(rc.runRead)
	read.Arg("name", "Secret name").StringVar(&rc.Name)
	read.Flag("decrypt", "Return decrypted value").BoolVar(&rc.Decrypt)
}

func (rc *readCommand) runRead(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		return err
	}
	ssmClient := ssm.New(sess, config)

	// read the secret to the parameter store
	gpInput := &ssm.GetParameterInput{
		Name:           &rc.Name,
		WithDecryption: &rc.Decrypt,
	}
	gpOutput, err := ssmClient.GetParameter(gpInput)
	if err != nil {
		return err
	}

	fmt.Println(*gpOutput.Parameter.Value)

	return nil
}
