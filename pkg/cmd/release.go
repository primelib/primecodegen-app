package cmd

import (
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/primelib/primecodegen-app/pkg/tasks/createtag"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func releaseCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "release",
		Aliases: []string{"r"},
		Run: func(cmd *cobra.Command, args []string) {
			// tasks
			tasks := []taskcommon.Task{createtag.NewTask()}

			// platform
			platform, err := vcsapp.GetPlatformFromEnvironment()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to configure platform from environment")
			}

			// execute
			err = vcsapp.ExecuteTasks(platform, tasks)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to execute release task")
			}
		},
	}
	cmd.Flags().Bool("dry-run", false, "Perform a dry run without making any changes")

	return cmd
}
