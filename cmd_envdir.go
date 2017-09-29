package main

import (
	"os"
	"path"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	"gopkg.in/alecthomas/kingpin.v2"
)

type envDirCommand struct {
	Path      string
	OutputDir string
}

func configureEnvDirCommand(app *kingpin.Application) {
	edc := &envDirCommand{}
	envDir := app.Command("envdir", "Write secrets into environment variable files (e.g. chpst -e)").Action(edc.runEnvDir)
	envDir.Flag("with-path", "Path to fetch secrets from").Required().StringVar(&edc.Path)
	envDir.Arg("output-dir", "The output directory").StringVar(&edc.OutputDir)
}

func (ec *envDirCommand) runEnvDir(ctx *kingpin.ParseContext) error {
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

	// write secrets to environment var files in output directory
	for _, param := range gpOutput.Parameters {
		f, err := os.OpenFile(path.Join(ec.OutputDir, nameToEnv(param.Name)), os.O_RDWR|os.O_CREATE, 0640)
		if err != nil {
			return err
		}

		if _, err = f.WriteString(*param.Value); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}
