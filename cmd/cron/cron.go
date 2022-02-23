package cron

import (
	"sync"

	"github.com/eyecuelab/kit/cmd"
	"github.com/spf13/cobra"
)

// CronCommand ...
var CronCommand = &cobra.Command{
	Use:   "cron",
	Short: "cron job",
	Run:   cronCmd,
}

var wg sync.WaitGroup

// Init ...
func Init() {
	cmd.Add(CronCommand)
}

func cronCmd(cmd *cobra.Command, args []string) {
	processSomethingExample()
	wg.Wait()
}
