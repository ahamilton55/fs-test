package deployer

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ahamilton55/fs-test/pad/lib/packager"
	"github.com/ahamilton55/fs-test/pad/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CfnDeployer struct {
	Template     string
	Parameters   map[string][]string
	Capabilities []string
	Config       utils.CommandConfig
	P            packager.Packager
}

const DefaultWebSiteUrlKey = "WebsiteURL"

func (cd CfnDeployer) Deploy() (DeployOutput, error) {
	var newStack bool
	stackName := fmt.Sprintf("%s-%s", cd.Config.Env, cd.Config.Service)

	sess, err := utils.GetAWSSession(cd.Config.Region, cd.Config.Profile)
	if err != nil {
		utils.ErrorAndQuit("Error getting AWS Session", err, 3)
	}

	cf := cloudformation.New(sess)
	descStackInput := &cloudformation.DescribeStacksInput{
		StackName: aws.String(stackName),
	}
	_, err = cf.DescribeStacks(descStackInput)
	if err != nil && strings.Contains(err.Error(), "does not exist") {
		newStack = true
	} else if err != nil {
		utils.ErrorAndQuit("Error checking the stack's status", err, 5)
	}

	if cd.Template == "" {
		utils.ErrorAndQuit("No CF template set", nil, 3)
	}

	var parameters []*cloudformation.Parameter
	parameters, err = setupParameters(cd, cd.Parameters)
	if err != nil {
		utils.ErrorAndQuit("Error building parameters", err, 5)
	}

	err = launchStack(newStack, stackName, cd, parameters, cf)
	if err != nil && strings.Contains(err.Error(), "No updates are to be performed") {
		log.Println("Nothing to update")
		return DeployOutput{}, nil
	} else if err != nil {
		utils.ErrorAndQuit("Unable to setup stack", err, 6)
	} else {
		if newStack {
			if err = cf.WaitUntilStackCreateComplete(descStackInput); err != nil {
				utils.ErrorAndQuit("Stack creation was not successful", err, 7)
			}
		} else {
			if err = cf.WaitUntilStackUpdateComplete(descStackInput); err != nil {
				utils.ErrorAndQuit("Stack update was not successful", err, 7)
			}
		}
	}

	stackInfo, err := cf.DescribeStacks(descStackInput)
	if err != nil {
		utils.ErrorAndQuit("Error looking up stack", err, 8)
	}

	var websiteUrl string
	for _, output := range stackInfo.Stacks[0].Outputs {
		if aws.StringValue(output.OutputKey) == DefaultWebSiteUrlKey {
			websiteUrl = aws.StringValue(output.OutputValue)
		}
	}

	return DeployOutput{WebsiteUrl: websiteUrl}, nil
}

func setupParameters(cd CfnDeployer, params map[string][]string) ([]*cloudformation.Parameter, error) {
	var parameters []*cloudformation.Parameter

	for _, val := range params[cd.Config.Env] {
		parameter := strings.Split(val, "=")
		parameters = append(parameters, &cloudformation.Parameter{
			ParameterKey:   aws.String(parameter[0]),
			ParameterValue: aws.String(parameter[1]),
		})
	}

	s3Package, err := cd.P.FindPackage(cd.Config.Params["pkgVer"])
	if err != nil {
		return parameters, err
	}

	parameters = append(parameters, &cloudformation.Parameter{
		ParameterKey:   aws.String("Package"),
		ParameterValue: aws.String(s3Package),
	})

	return parameters, nil
}

func launchStack(newStack bool, stackName string, cd CfnDeployer, parameters []*cloudformation.Parameter, cf *cloudformation.CloudFormation) error {
	if newStack {
		stackInput := cloudformation.CreateStackInput{}

		stackInput.Parameters = parameters
		stackInput.StackName = aws.String(stackName)
		if strings.Contains(cd.Template, "s3://") {
			stackInput.TemplateURL = aws.String(cd.Template)
		} else {
			contents, err := ioutil.ReadFile(cd.Template)
			if err != nil {
				return err
			}
			stackInput.TemplateBody = aws.String(string(contents))
		}

		for _, cap := range cd.Capabilities {
			stackInput.Capabilities = append(stackInput.Capabilities, aws.String(cap))
		}

		_, err := cf.CreateStack(&stackInput)
		if err != nil {
			return err
		}
	} else {
		stackInput := cloudformation.UpdateStackInput{}

		stackInput.Parameters = parameters
		stackInput.StackName = aws.String(stackName)
		if strings.Contains(cd.Template, "s3://") {
			stackInput.TemplateURL = aws.String(cd.Template)
		} else {
			contents, err := ioutil.ReadFile(cd.Template)
			if err != nil {
				return err
			}
			stackInput.TemplateBody = aws.String(string(contents))
		}

		for _, cap := range cd.Capabilities {
			stackInput.Capabilities = append(stackInput.Capabilities, aws.String(cap))
		}

		_, err := cf.UpdateStack(&stackInput)
		if err != nil {
			return err
		}
	}

	return nil
}
