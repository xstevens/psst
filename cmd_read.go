package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type ReadCommand struct {
	Name    string
	Value   string
	Decrypt bool
}

func configureReadCommand(app *kingpin.Application) {
	rc := &ReadCommand{}
	read := app.Command("read", "Read secret from parameter store").Action(rc.runRead)
	read.Arg("name", "Secret name").StringVar(&rc.Name)
	read.Arg("value", "Secret value").StringVar(&rc.Value)
	read.Flag("decrypt", "Return decrypted value").BoolVar(&rc.Decrypt)
}

func (rc *ReadCommand) runRead(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open session: %v\n", err.Error())
		return err
	}
	ssmClient := ssm.New(sess, config)

	// read the secret to the parameter store
	gpinput := &ssm.GetParameterInput{
		Name:           &rc.Name,
		WithDecryption: &rc.Decrypt,
	}
	gpoutput, err := ssmClient.GetParameter(gpinput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read secret: %v\n", err.Error())
		return err
	}

	fmt.Println(gpoutput.GoString())

	return nil
}
