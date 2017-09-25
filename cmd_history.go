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

type historyCommand struct {
	Name    string
	Decrypt bool
}

func configureHistoryCommand(app *kingpin.Application) {
	hc := &historyCommand{}
	history := app.Command("history", "Get secret history from parameter store").Action(hc.runHistory)
	history.Arg("name", "Secret name").StringVar(&hc.Name)
	history.Flag("decrypt", "Return decrypted value").BoolVar(&hc.Decrypt)
}

func writeHistoryTable(history []*ssm.ParameterHistory) {
	writer := tabwriter.NewWriter(os.Stdout, 4, 0, 4, ' ', 0)
	fmt.Fprintln(writer, "Value\tKey Alias\tLast Modified\tUser\t")
	for _, entry := range history {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t\n",
			*entry.Value,
			*entry.KeyId,
			entry.LastModifiedDate.Format(time.RFC3339),
			*entry.LastModifiedUser)
	}
	writer.Flush()
}

func (hc *historyCommand) runHistory(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open session")
		return err
	}
	ssmClient := ssm.New(sess, config)

	// get the history of the secret from parameter store
	gphInput := &ssm.GetParameterHistoryInput{
		Name:           aws.String(hc.Name),
		WithDecryption: aws.Bool(hc.Decrypt),
	}
	gphOutput, err := ssmClient.GetParameterHistory(gphInput)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read secret history")
		return err
	}

	writeHistoryTable(gphOutput.Parameters)

	return nil
}
