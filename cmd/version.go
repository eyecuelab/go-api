package cmd

import (
	"github.com/eyecuelab/kit/cmd"
	"github.com/eyecuelab/kit/log"

	"github.com/eyecuelab/go-api/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run:   run,
}

var foo string

func init() {
	cmd.Add(versionCmd)
}

func run(cmd *cobra.Command, args []string) {
	log.Infof("Version: %s", version.Version)
	log.Infof("Git Rev: %s", version.GitRev)
	log.Infof("Date: %s", version.Date)
}
