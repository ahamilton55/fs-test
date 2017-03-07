package deployer

type Deployer interface {
	Deploy() (DeployOutput, error)
}

type DeployerConfig struct {
	Deployer   string
	Template   string
	Parameters map[string][]string
}

type DeployOutput struct {
	WebsiteUrl string
}
