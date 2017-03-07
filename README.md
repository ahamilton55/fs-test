# Full Screen Homework Challenge Response

I've built a tool in Go that will package and deploy basic static sites. It's built to be extensible with different plugins for packaging and deploying. To learn move see the README.md in the `pad` directory.

## Requirements
1. AWS credentials setup with `aws configure`.
1. AWS S3 bucket that can be used to upload packages.
1. AWS VPC public and private subnets for application deployment.
1. Go 1.8 installed and setup.

## Quickstart.

Use the following steps to get started.

1. Clone the repo using `go get`:
```bash
$ go get -u github.com/ahamilton55/fs-test
```

1. Go to the example directory in the newly cloned repo:
```bash
$ cd ${GOPATH}/src/github.com/ahamilton55/fs-test/example
```

1. Edit the configuration file example in `.deploy/config.toml` with the following information:

  * S3 Bucket (Updated in both places with the same value)
  * Public VPC Subnet IDs (comman separated list)
  * Private VPC Subnet IDs (comman separated list)
  * VPC ID
  * Public Key (if you want to login to the instance)

1. Package the example and upload it to the provided S3 bucket:
```bash
$ ${GOPATH}/bin/pad pack
```

1. Deploy the example package using AWS CloudFormation
```bash
$ ${GOPATH}/bin/pad deploy
```
