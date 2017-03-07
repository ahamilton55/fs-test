package cmd

import (
	"log"
	"os"

	"github.com/ahamilton55/fs-test/pad/utils"
	"github.com/spf13/cobra"
)

var (
	CmdConfig  utils.CommandConfig
	configFile string // Config file location

	RootCmd = &cobra.Command{
		Use:   "bad",
		Short: "Build and deploy (bad) a web service",
		Run: func(cmd *cobra.Command, args []string) {
			//Do stuff here
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			setup(configFile)
		},
	}
)

func init() {
	RootCmd.PersistentFlags().StringVar(&CmdConfig.Region, "region", "us-west-2", "AWS region for deploy (default: us-west-2)")
	RootCmd.PersistentFlags().StringVar(&CmdConfig.Env, "env", "stage", "Service environment deployed (default: stage)")
	RootCmd.PersistentFlags().StringVar(&CmdConfig.Service, "service", "", "Service to be deployed")
	RootCmd.PersistentFlags().StringVar(&CmdConfig.Profile, "profile", "", "AWS profile")
	RootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file location")
}

func setup(configFile string) {
	var err error

	CmdConfig.Params = make(map[string]string)

	CmdConfig.Config, err = utils.ReadConfig(configFile)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	if CmdConfig.Service == "" {
		CmdConfig.Service = CmdConfig.Config.Service
	}
}
