package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ahamilton55/fs-test/pad/lib/deployer"
	"github.com/ahamilton55/fs-test/pad/utils"
	"github.com/spf13/cobra"
)

var (
	pkgVer string

	deployCli = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an ECR service with Cloudformation",
		Long:  `Deploy an ECR service with Cloudformation`,
		Run: func(cmd *cobra.Command, args []string) {
			deploy(CmdConfig)
		},
	}
)

func init() {
	deployCli.Flags().StringVar(&pkgVer, "pkg-ver", "latest", "Version of the package deploy (default: latest)")
	RootCmd.AddCommand(deployCli)
}

func deploy(config utils.CommandConfig) {
	CmdConfig.Params["pkgVer"] = pkgVer

	var d deployer.Deployer

	switch CmdConfig.Config.DeployerType {
	case "cloudformation":
		var cfnDeployer deployer.CfnDeployer
		if _, err := toml.Decode(config.Config.DeployerArgs, &cfnDeployer); err != nil {
			log.Printf("Error decoding deployer config: %s", err.Error())
			os.Exit(1)
		}

		cfnDeployer.Config = CmdConfig
		cfnDeployer.P = getPackager()

		d = cfnDeployer
	}

	out, err := d.Deploy()
	if err != nil {
		utils.ErrorAndQuit("Error running deploy", err, 10)
	}

	fmt.Printf("URL: %s\n", out.WebsiteUrl)
}
