package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type writeCommand struct {
	Name      string
	Value     string
	Overwrite bool
}

func configureWriteCommand(app *kingpin.Application) {
	wc := &writeCommand{}
	write := app.Command("write", "Write secret to parameter store").Action(wc.runWrite)
	write.Arg("name", "Secret name").StringVar(&wc.Name)
	write.Arg("value", "Secret value").StringVar(&wc.Value)
	write.Flag("overwrite", "Overwrite the existing secret").Default("false").Bool()
}

func (wc *writeCommand) runWrite(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open session: %v\n", err.Error())
		return err
	}
	ssmClient := ssm.New(sess, config)

	// write the secret to the parameter store
	ppinput := &ssm.PutParameterInput{
		KeyId:     kmsAlias,
		Name:      aws.String(wc.Name),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Value:     aws.String(wc.Value),
		Overwrite: aws.Bool(wc.Overwrite),
	}
	ppoutput, err := ssmClient.PutParameter(ppinput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write secret: %v\n", err.Error())
		return err
	}

	fmt.Println(ppoutput.GoString())

	return nil
}
