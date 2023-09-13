package cmd

import (
	"os"
	"slices"

	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/primelib/primelib-app/pkg/tasks/codegeneration"
	"github.com/primelib/primelib-app/pkg/tasks/createtag"
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
		var tasks []taskcommon.Task
		if slices.Contains(args, "generate") {
			tasks = append(tasks, codegeneration.NewTask())
		}
		if slices.Contains(args, "release") {
			tasks = append(tasks, createtag.NewTask())
		}

		if len(tasks) == 0 {
			log.Fatal().Msg("no tasks specified, supported tasks are: generate, release")
			os.Exit(1)
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
