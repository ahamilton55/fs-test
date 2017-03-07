package cmd

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ahamilton55/fs-test/pad/lib/packager"
	"github.com/ahamilton55/fs-test/pad/utils"
	"github.com/spf13/cobra"
)

var (
	packCli = &cobra.Command{
		Use:   "pack",
		Short: "Pack a service and push the package to a remote storage location",
		Long:  `Pack a service and push the package to a remote storage location`,
		Run: func(cmd *cobra.Command, args []string) {
			pack()
		},
	}
)

func init() {
	RootCmd.AddCommand(packCli)
}

func pack() {
	st := getPackager()

	tgz, err := st.Build()
	if err != nil {
		utils.ErrorAndQuit("", err, 1)
	}

	err = st.Push(tgz)
	if err != nil {
		utils.ErrorAndQuit("", err, 2)
	}

	err = st.Cleanup(tgz)
	if err != nil {
		utils.ErrorAndQuit("", err, 3)
	}
}

func getPackager() packager.Packager {
	var p packager.Packager

	switch CmdConfig.Config.PackagerType {
	case "s3tarball":
		var st packager.S3Tarball
		if _, err := toml.Decode(CmdConfig.Config.PackagerArgs, &st); err != nil {
			log.Printf("Error decoding deployer config: %s", err.Error())
			os.Exit(1)
		}

		st.Config = CmdConfig
		p = st
	}

	return p

}
