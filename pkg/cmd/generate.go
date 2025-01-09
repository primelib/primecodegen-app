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

func generateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"g"},
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := cmd.Flags().GetString("dir")

			if dir == "" {
				generateApp()
			} else {
				generateLocal(dir)
			}
		},
	}
	cmd.Flags().Bool("dry-run", false, "Perform a dry run without making any changes")
	cmd.Flags().String("dir", "", "Directory of the project for local code generation")

	return cmd
}

func generateApp() {
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

func generateLocal(dir string) {
	configPath := path.Join(dir, "primelib.yaml")
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal().Err(err).Str("config-path", configPath).Msg("failed to read primelib.yaml")
	}

	// load config
	conf, err := config.FromString(string(bytes))
	if err != nil {
		log.Fatal().Err(err).Str("config-path", configPath).Msg("failed to parse primelib.yaml")
	}

	// for each module
	log.Info().Str("dir", dir).Str("config", configPath).Msg("running local generation")
	genErr := primelib.Generate(dir, conf, api.Repository{})
	if genErr != nil {
		log.Warn().Err(genErr).Msg("failed to generate code")
	}
}
