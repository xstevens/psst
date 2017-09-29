# psst
A command-line client for storing secrets in AWS Parameter Store.

## Install
```
go get github.com/xstevens/psst
```

### Build

To run subsequent builds, use `go build`:

```
# Ensure you're in the `psst` source directory.
cd $GOPATH/src/github.com/xstevens/psst

# Run the build.
go build
```

### Cross-compiling
With Go 1.5 or above, cross-compilation support is built in.
See [Dave Cheney's blog post](http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5)
for a tutorial and the [golang.org docs](https://golang.org/doc/install/source#environment)
for details on `GOOS` and `GOARCH` values for various target operating systems.

A typical build for Linux would be:
```
# Ensure you're in the `psst` source directory.
cd $GOPATH/src/github.com/xstevens/psst

# Run the build.
GOOS=linux GOARCH=amd64 go build
```

## Usage
```
$ ./psst -h
usage: psst [<flags>] <command> [<args> ...]

A command-line client for storing secrets in AWS Parameter Store.

Flags:
  -h, --help                 Show context-sensitive help (also try --help-long and --help-man).
      --region="us-east-1"   AWS region
      --kms="alias/aws/ssm"  KMS key alias
      --mfa=MFA              IAM MFA device ARN
      --role=ROLE            IAM role ARN to assume
      --version              Show application version.

Commands:
  help [<command>...]
    Show help.

  read [<flags>] [<name>]
    Read secret from parameter store

  write [<flags>] [<name>]
    Write secret to parameter store

  delete [<name>]
    Delete secret from parameter store

  exec --with-path=WITH-PATH [<command>...]
    Execute command with secrets populated in the environment

  envdir --with-path=WITH-PATH [<output-dir>]
    Write secrets into environment variable files (e.g. chpst -e)

  list <path>
    List all secrets under a path in parameter store

  history [<flags>] [<name>]
    Get secret history from parameter store
```

## License
All aspects of this software are distributed under the MIT License. See LICENSE file for full license text.

## Inspirations and similar work
- [Chamber](https://github.com/segmentio/chamber)
- [Confidant](https://lyft.github.io/confidant)
- [Keywhiz](https://square.github.io/keywhiz)
- [Sneaker](https://github.com/codahale/sneaker)
- [Vault](https://www.vaultproject.io)
