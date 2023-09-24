package cmd

import (
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/primelib/primelib-app/pkg/tasks/codegeneration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Run: func(cmd *cobra.Command, args []string) {
			// tasks
			tasks := []taskcommon.Task{codegeneration.NewTask()}

			// platform
			platform, err := vcsapp.GetPlatformFromEnvironment()
			if err != nil {
				log.Fatal().Err(err).Msg("failed to configure platform from environment")
			}

			// execute
			err = vcsapp.ExecuteTasks(platform, tasks)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to execute generate task")
			}
		},
	}
	cmd.Flags().Bool("dry-run", false, "Perform a dry run without making any changes")

	return cmd
}
