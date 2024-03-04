package cmd

import (
	"os"
	"path"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/primelib"
	"github.com/primelib/primelib-app/pkg/tasks/codegeneration"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func updateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Aliases: []string{"u"},
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("dir")

			if dir == "" {
				updateTaskApp()
			} else {
				updateLocal(dir)
			}
		},
	}
	cmd.Flags().Bool("dry-run", false, "Perform a dry run without making any changes")
	cmd.Flags().String("dir", "", "Directory of the project for local code generation")

	return cmd
}

func updateTaskApp() {
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
}

func updateLocal(dir string) {
	configPath := path.Join(dir, config.ConfigFileName)
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Str("config-path", configPath).Msg("failed to read primelib.yaml")
	}

	// load config
	conf, err := config.ConfigFromString(string(bytes))
	if err != nil {
		log.Fatal().Err(err).Str("config-path", configPath).Msg("failed to parse primelib.yaml")
	}

	// for each module
	log.Info().Str("dir", dir).Str("config", configPath).Msg("running local update")
	err = primelib.Update(dir, conf, api.Repository{})
	if err != nil {
		log.Warn().Err(err).Msg("failed to update spec")
	}
}
