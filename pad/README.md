# Package and Deploy Tool

The Package and Deploy (pad) tool does just what it says. It allows for the packaging to a service and the the deployment of that service.

The packers and deployers are built to allow for extension and multiple types to satisfy multiple service. Additional packagers and deployers may be created for different scenarios without requiring a massive change to the code. The new packager or deployer would simply need to follow the current Go Interfaces and are basically plug and play. For example, one service may use the S3Tarball packager to create a tarball and upload it to S3 while another service may want a Docker or container packager which could be easily built.

The Package and Deploy tool is controlled using a configuration file with a default location of `.deploy/config.toml`. Different configuration files may be specified at runtime using a flag.

## Building

To build and install the tool run the following in this directory:

```bash
$ go install
```

This will build the tool and place it in `${GOPATH}/bin/pad`.

## Usage

The `pad` command has two subcommands: `pack` and `deploy`.

* `pack`: Use this sub-command to package up your service for deployment using on of the available packages.
* `deploy`: Use this sub-command to deploy your service using one of the available deployers.

### Usage Examples
Package and deploy a service using a different configuration file and then deploy to a stage environment.
```bash
$ ${GOPATH}/bin/pad --config .deploy/another.toml pack
$ ${GOPATH}/bin/pad --config .deploy/another.toml --env stage deploy
```

## configuration
Configuration for the service is placed in a configuration file in TOML format (similar to INI). This configuration file will hold information that is standard to the packaging and deployment of the service. Multiple configuration files may be used but multiple environments can live in a single configuration. This tool should make it easy to build adhoc development and test environments using the same CloudFormation templates as staging and production just by adding a new set of parameters to the configuration file.

Information such as AWS region, AWS profile, or environment are passed using command line flags to the pad tool.

## Default Configuration Values
The default configuration currently provided is for the service name.

```toml
Service = "example"
```

## Packagers
The following packagers are currently available:

* S3Tarball

### S3Tarball
The S3Tarball packager will create a gzipped tarball of the service directory, excluding any "dot" directories, and then upload that package to S3. The bucket used for upload is provided in the configuration file and sub-directory of the service name is created for storing packages.

Packages are created in the following format `YYMMDDHHMM-${SERVICENAME}` and a file called `latest-build` is updated to container the name of the latest, uploaded filename. The later is done to make it easier to find the latest build when doing a deploy.

#### Configuration Options

* `Bucket`: Name of the bucket where packages are stored.

### Configuration Example
```toml
PackagerType = "s3tarball"
PackagerArgs = """
Bucket = "fs-test-packages"
"""
```

* `PackagerType`: Designates what packager to use.
* `PackagerArgs`: A string in TOML format that is later decoded to provide the configuration for the specified packager.

## Deployers
The following deployers are currently available:

* Cloudformation

### Cloudformation
The Cloudformation deployer will utilize AWS CloudFormation to deploy the application. A CloudFormation template and a set of Parameters are provided by the service to be used for deployment.

#### Configuration Options
* `Template`: Location of a CloudFormation template to be used for deployment of the service.
* `Capabilities`: A list of capabilities that are required to run the cloudformation template.
* `Parameters`: A set of parameters for different environments (passed to the command using the `--env` flag.)

### Configuration Example
```toml
DeployerType = "cloudformation"
DeployerArgs = """
Template = ".deploy/site.template"
Capabilities = [ "CAPABILITY_IAM" ]

[Parameters]
stage = [
  "BucketName=fs-test-packages",
  "Environment=stage",
  "KeyName=mykey",
  "PrivateSubnetIDs=subnet-aaaaaa,subnet-bbbbbb,subnet-cccccc",
  "PublicSubnetIDs=subnet-111111,subnet-222222,subnet-333333",
  "SiteName=example",
  "VpcId=vpc-abcdef01",
]
"""
```

* `DeployerType`: Desginates what deployer to use.
* `DeployerArgs`: A string in TOML format that is later decoded to provide the configuration for the specified packager.
