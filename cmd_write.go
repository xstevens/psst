package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"gopkg.in/alecthomas/kingpin.v2"
)

type writeCommand struct {
	Name      string
	Overwrite bool
}

func configureWriteCommand(app *kingpin.Application) {
	wc := &writeCommand{}
	write := app.Command("write", "Write secret to parameter store").Action(wc.runWrite)
	write.Arg("name", "Secret name").StringVar(&wc.Name)
	write.Flag("overwrite", "Overwrite the existing secret").Default("false").BoolVar(&wc.Overwrite)
}

func (wc *writeCommand) runWrite(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		return err
	}
	ssmClient := ssm.New(sess, config)

	// read secret value from stdin and convert to a string
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	val := strings.TrimSpace(string(data))

	// write the secret to the parameter store
	ppInput := &ssm.PutParameterInput{
		KeyId:     kmsAlias,
		Name:      aws.String(wc.Name),
		Type:      aws.String(ssm.ParameterTypeSecureString),
		Value:     aws.String(val),
		Overwrite: aws.Bool(wc.Overwrite),
	}
	_, err = ssmClient.PutParameter(ppInput)

	return err
}
