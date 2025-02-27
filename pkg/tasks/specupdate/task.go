package codegeneration

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/cidverse/go-vcsapp/pkg/task/simpletask"
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	cp "github.com/otiai10/copy"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/primelib"
	"github.com/primelib/primelib-app/pkg/specutil"
	"github.com/rs/zerolog/log"
)

const branchName = "feat/primelib-spec"

//go:embed templates/description.gohtml
var descriptionTemplate []byte

type SpecUpdateTask struct{}

// Name returns the name of the task
func (n SpecUpdateTask) Name() string {
	return "generate"
}

// Execute runs the task
func (n SpecUpdateTask) Execute(ctx taskcommon.TaskContext) error {
	content, err := ctx.Platform.FileContent(ctx.Repository, ctx.Repository.DefaultBranch, config.ConfigFileName)
	if err != nil {
		return fmt.Errorf("failed to get %s content: %w", config.ConfigFileName, err)
	}

	// load config
	conf, err := config.FromString(content)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", config.ConfigFileName, err)
	}

	// create helper
	helper := simpletask.New(ctx)

	// clone repository
	err = helper.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	// create and checkout new branch
	branch := branchName
	err = helper.CreateBranch(branch)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	// store original spec file
	specFile := path.Join(ctx.Directory, conf.Spec.File)
	originalSpecFile, err := os.CreateTemp("", "primelib-openapi-*"+filepath.Ext(conf.Spec.File))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	err = cp.Copy(specFile, originalSpecFile.Name())
	defer os.Remove(originalSpecFile.Name())

	// update spec
	err = primelib.Update(ctx.Directory, conf, ctx.Repository)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// store updated spec file
	diff, err := specutil.DiffSpec("openapi", originalSpecFile.Name(), specFile)
	if err != nil {
		log.Warn().Err(err).Msg("failed to diff spec file")
	}
	if len(diff.OpenAPI) > 15 {
		diff.OpenAPI = diff.OpenAPI[:15] // limit to the first n changes, sorted by level
	}

	// commit message and description
	changes, err := helper.VCSClient.UncommittedChanges()
	if err != nil {
		return fmt.Errorf("failed to get uncommitted changes: %w", err)
	}
	commitMessage := "feat: update openapi spec"
	description, err := vcsapp.Render(string(descriptionTemplate), map[string]interface{}{
		"PlatformName": ctx.Platform.Name(),
		"PlatformSlug": ctx.Platform.Slug(),
		"Name":         conf.Name,
		"SpecDiff":     diff,
		"Footer":       os.Getenv("PRIMEAPP_FOOTER_HIDE") != "true",
		"FooterCustom": os.Getenv("PRIMEAPP_FOOTER_CUSTOM"),
	})
	if err != nil {
		return fmt.Errorf("failed to render description template: %w", err)
	}

	// do not commit if only .openapi-generator/FILES changed
	if len(changes) == 0 {
		log.Info().Int("total-changes", len(changes)).Msg("no changes detected, skipping commit and merge request")
		return nil
	}

	// commit push and create or update merge request
	err = helper.CommitPushAndMergeRequest(commitMessage, commitMessage, string(description), "")
	if err != nil {
		return fmt.Errorf("failed to commit push and create or update merge request: %w", err)
	}

	return nil
}

func NewTask() SpecUpdateTask {
	return SpecUpdateTask{}
}
