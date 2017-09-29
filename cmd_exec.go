package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type execCommand struct {
	Path    string
	Command []string
}

func configureExecCommand(app *kingpin.Application) {
	ec := &execCommand{}
	exec := app.Command("exec", "Execute command with secrets populated in the environment").Action(ec.runExec)
	exec.Flag("with-path", "Path to fetch secrets from").Required().StringVar(&ec.Path)
	exec.Arg("command", "The command to execute").StringsVar(&ec.Command)
}

func (ec *execCommand) runExec(ctx *kingpin.ParseContext) error {
	config := aws.NewConfig().WithRegion(*region).WithCredentialsChainVerboseErrors(true)
	sess, err := newSession(config, mfaSerial, roleArn)
	if err != nil {
		return err
	}
	ssmClient := ssm.New(sess, config)

	// get parameters; recursive path search and decrypted
	gpInput := &ssm.GetParametersByPathInput{
		Path:           aws.String(ec.Path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}
	gpOutput, err := ssmClient.GetParametersByPath(gpInput)
	if err != nil {
		return err
	}

	// construct new environment where secrets override OS environment presets
	env := []string{}
	for _, param := range gpOutput.Parameters {
		env = append(env, fmt.Sprintf("%s=%s", nameToEnv(param.Name), *param.Value))
	}
	env = append(env, os.Environ()...)

	return execCommandWithEnv(ec.Command, env)
}
