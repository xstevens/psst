package main

import (
	"fmt"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"path"
	"strings"

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

func (ec *execCommand) nameToEnv(name *string) string {
	// paths are expected to be in the form:
	// `/env/service/component/secret_name`

	// get secret_name
	envName := fmt.Sprintf("%s_%s", path.Base(path.Dir(*name)), path.Base(*name))
	// upper
	envName = strings.ToUpper(envName)
	// replace dots with underscores if there are any
	envName = strings.Replace(envName, ".", "_", -1)
	// replace hyphen with underscores if there are any
	envName = strings.Replace(envName, "-", "_", -1)

	return envName
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
		env = append(env, fmt.Sprintf("%s=%s", ec.nameToEnv(param.Name), *param.Value))
	}
	env = append(env, os.Environ()...)

	return execCommandWithEnv(ec.Command, env)
}
