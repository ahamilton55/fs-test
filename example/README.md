# Example Web SiteName

A basic "Hello World" site to show off the Package and Deploy tool.

Configurations are kept in the `.deploy/` directory.

## Requirements

* **S3 Bucket**: A S3 bucket is required for storing gzipped tarballs for deployment. This bucket needs to be added to the configuration file at `.deploy/config.toml`.
* **VPC Subnetting**: The setup requires at least one public subnet and one private subnet. Ideally, three public and three private subnets would be available and 3 instances would be run but, since this is a test application a single public subnet and single private subnet may be used.
* **Key Pair**: It is assumed that there is already a public key uploaded to AWS.

## Defaults
* **AMI**: Amazon Linux x86_64 HVM image
* **Instance Type**: t2.micro
* **ELB Health Endpoint**: HTTP:80/index.html
* **Desired Capacity**: 1
* **Max Capacity**: 2

## Infrastructure created
* **AutoScaling Group**
* **Launch Configuration**
* **CloudWatch Alarms**
  * **High CPU (Scale out) Alarm**
  * **Low CPU (Scale in) Alarm**
* **AutoScaling Actions**
  * **Scale Out Action**
  * **Scale In Action**
* **EC2 Instance Role**
* **EC2 Instance Profile**
* **Security Groups**
  * **Instance Security Group**
  * **ELBv2 Security Group**
* **ELBv2**
  * **ELBv2 Listener**
* **Target Group**

## Outputs
* **WebsiteURL**: URL for the ELBv2.

## Provisioning
Provisioning of the site is simple and is handled by a user-data script inside of the CloudFormation template. The user-data script simply installs nginx, downloads the S3 package, unzips and untars the package into the default nginx directory, restarts nginx to make sure that it is running.

The user-data script also signals to Cloudformation that the instance has completed setup so that it can continue with other tasks.
