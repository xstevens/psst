package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func StdinTokenProvider() (string, error) {
	var v string
	fmt.Printf("MFA token code: ")
	_, err := fmt.Scanln(&v)

	return v, err
}

func newSession(config *aws.Config, serial *string, role *string) (*session.Session, error) {
	sess := session.Must(session.NewSession(config))
	if role != nil && len(*role) > 0 && serial != nil && len(*serial) > 0 {
		creds := stscreds.NewCredentials(sess, *role, func(p *stscreds.AssumeRoleProvider) {
			p.SerialNumber = aws.String(*serial)
			p.TokenProvider = StdinTokenProvider
		})
		config.WithCredentials(creds)
	} else if serial != nil && len(*serial) > 0 {
		token, err := StdinTokenProvider()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get token from stdin: %v\n", err.Error())
			return nil, err
		}
		svc := sts.New(sess)
		sessTokenInput := &sts.GetSessionTokenInput{
			SerialNumber: aws.String(*serial),
			TokenCode:    aws.String(token),
		}
		sessTokenOutput, err := svc.GetSessionToken(sessTokenInput)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get session token: %v\n", err.Error())
			return nil, err
		}
		credsVal := &credentials.Value{
			AccessKeyID:     *sessTokenOutput.Credentials.AccessKeyId,
			SecretAccessKey: *sessTokenOutput.Credentials.SecretAccessKey,
			SessionToken:    *sessTokenOutput.Credentials.SessionToken,
		}
		creds := credentials.NewStaticCredentialsFromCreds(*credsVal)
		config.WithCredentials(creds)
	}

	return sess, nil
}

// The run function runs a command in an environment.
// Stdout and stderr are preserved.
func execCommandWithEnv(command []string, env []string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	return cmd.Run()
}
