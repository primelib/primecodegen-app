package generatetask

import (
	_ "embed"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/cidverse/vcs-app/pkg/task/simpletask"
	"github.com/cidverse/vcs-app/pkg/task/taskcommon"
	"github.com/cidverse/vcs-app/pkg/vcsapp"
	"github.com/primelib/primelib-app/pkg/primelib"
	"github.com/rs/zerolog/log"
)

const branchName = "feat/primelib-generate"

//go:embed templates/description.gohtml
var descriptionTemplate []byte

type PrimeLibGenerateTask struct {
}

// Name returns the name of the task
func (n PrimeLibGenerateTask) Name() string {
	return "primelib-generate"
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

	// generate
	err = primelib.Execute(ctx.Directory, module)
	if err != nil {
		return fmt.Errorf("failed to generate: %w", err)
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
		"Module":      moduleName,
		"SpecUpdated": slices.Contains(changes, module.SpecFile),
		"CodeUpdated": len(changes) > 1,
	})
	if err != nil {
		return fmt.Errorf("failed to render description template: %w", err)
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
