package cmd

import (
	"github.com/cidverse/vcs-app/pkg/task/taskcommon"
	"github.com/cidverse/vcs-app/pkg/vcsapp"
	"github.com/primelib/primelib-app/pkg/generatetask"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:     "run",
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		// tasks
		var tasks = []taskcommon.Task{
			generatetask.NewTask(),
		}

		// platform
		platform, err := vcsapp.GetPlatformFromEnvironment()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to configure platform from environment")
		}

		// execute
		err = vcsapp.ExecuteTasks(platform, tasks)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to execute task(s)")
		}
	},
}
