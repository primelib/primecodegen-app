package codegeneration

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/cidverse/go-vcsapp/pkg/task/simpletask"
	"github.com/cidverse/go-vcsapp/pkg/task/taskcommon"
	"github.com/cidverse/go-vcsapp/pkg/vcsapp"
	cp "github.com/otiai10/copy"
	"github.com/primelib/primelib-app/pkg/primelib"
	"github.com/primelib/primelib-app/pkg/spec"
	"github.com/rs/zerolog/log"
)

const branchName = "feat/primelib-generate"

//go:embed templates/description.gohtml
var descriptionTemplate []byte

type PrimeLibGenerateTask struct {
}

// Name returns the name of the task
func (n PrimeLibGenerateTask) Name() string {
	return "generate"
}

// Execute runs the task
func (n PrimeLibGenerateTask) Execute(ctx taskcommon.TaskContext) error {
	content, err := ctx.Platform.FileContent(ctx.Repository, ctx.Repository.DefaultBranch, "primelib.yaml")
	if err != nil {
		return fmt.Errorf("failed to get primelib.yaml content: %w", err)
	}

	// load config
	config, err := primelib.ConfigFromString(content)
	if err != nil {
		return fmt.Errorf("failed to load primelib.yaml: %w", err)
	}

	// for each module
	for _, module := range config.Modules {
		err = n.ExecuteModule(ctx, module)
		if err != nil {
			log.Warn().Err(err).Str("module", module.Name).Msg("failed to execute for module")
		}
	}

	return nil
}

func (n PrimeLibGenerateTask) ExecuteModule(ctx taskcommon.TaskContext, module primelib.Module) error {
	// create temp directory (override, so we can run the modules individually)
	tempDir, err := os.MkdirTemp("", "vcs-app-*")
	if err != nil {
		return fmt.Errorf("failed to prepare temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	ctx.Directory = tempDir

	// create helper
	helper := simpletask.New(ctx)

	// clone repository
	err = helper.Clone()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	moduleName := toModuleName(module.Name)
	branch := branchName
	commitSuffix := ""
	if moduleName != "" {
		branch = fmt.Sprintf("%s-%s", branchName, moduleName)
		commitSuffix = fmt.Sprintf(" of module %s", moduleName)
	}

	// create and checkout new branch
	err = helper.CreateBranch(branch)
	if err != nil {
		return fmt.Errorf("failed to create branch: %w", err)
	}

	// store original spec file
	originalSpecFile, err := os.CreateTemp("", "primelib-openapi-*"+filepath.Ext(module.SpecFile))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	err = cp.Copy(path.Join(ctx.Directory, module.Dir, module.SpecFile), originalSpecFile.Name())
	defer os.Remove(originalSpecFile.Name())
	if err != nil {
		os.Remove(originalSpecFile.Name())
	}

	// generate
	err = primelib.Execute(ctx.Directory, module, ctx.Repository)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
	}

	// store updated spec file
	updatedSpecFile := path.Join(ctx.Directory, module.Dir, module.SpecFile)
	diff, err := spec.DiffSpec(module.SpecFormat, originalSpecFile.Name(), updatedSpecFile)
	if err != nil {
		log.Warn().Err(err).Msg("failed to diff spec file")
	}
	if len(diff.OpenAPI) > 20 {
		diff.OpenAPI = diff.OpenAPI[:20] // limit to the first n changes, sorted by level
	}

	// commit message and description
	changes, err := helper.VCSClient.UncommittedChanges()
	if err != nil {
		return fmt.Errorf("failed to get uncommitted changes: %w", err)
	}
	commitMessage := fmt.Sprintf("feat: update generated code%s", commitSuffix)
	if slices.Contains(changes, module.SpecFile) {
		commitMessage = fmt.Sprintf("feat: update openapi spec%s", commitSuffix)
	}
	description, err := vcsapp.Render(string(descriptionTemplate), map[string]interface{}{
		"PlatformName": ctx.Platform.Name(),
		"PlatformSlug": ctx.Platform.Slug(),
		"Module":       moduleName,
		"SpecUpdated":  slices.Contains(changes, path.Join(module.Dir, module.SpecFile)),
		"CodeUpdated":  len(changes) > 1,
		"SpecDiff":     diff,
		"Footer":       os.Getenv("PRIMEAPP_FOOTER_HIDE") != "true",
		"FooterCustom": os.Getenv("PRIMEAPP_FOOTER_CUSTOM"),
	})
	if err != nil {
		return fmt.Errorf("failed to render description template: %w", err)
	}

	// do not commit if only .openapi-generator/FILES changed
	if len(changes) == 1 && strings.HasSuffix(changes[0], ".openapi-generator/FILES") {
		log.Info().Msg("no changes detected, skipping commit and merge request")
		return nil
	}

	// commit push and create or update merge request
	err = helper.CommitPushAndMergeRequest(commitMessage, commitMessage, string(description), "")
	if err != nil {
		return fmt.Errorf("failed to commit push and create or update merge request: %w", err)
	}

	return nil
}

func NewTask() PrimeLibGenerateTask {
	return PrimeLibGenerateTask{}
}

func toModuleName(input string) string {
	if input != "" && input != "root" {
		return strings.ToLower(input)
	}

	return ""
}
