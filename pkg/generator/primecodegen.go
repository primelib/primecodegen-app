package generator

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

type PrimeCodeGenGenerator struct {
	OutputName string   `json:"-" yaml:"-"`
	APISpec    string   `json:"-" yaml:"-"`
	Args       []string `json:"-" yaml:"-"`
	Config     PrimeCodeGenGeneratorConfig
}

type PrimeCodeGenGeneratorConfig struct {
	TemplateLanguage string   `json:"templateLanguage" yaml:"templateLanguage"`
	TemplateType     string   `json:"templateType" yaml:"templateType"`
	Patches          []string `json:"patches" yaml:"patches"`
	GroupId          string   `json:"groupId" yaml:"groupId"`
	ArtifactId       string   `json:"artifactId" yaml:"artifactId"`
}

// Name returns the name of the task
func (n *PrimeCodeGenGenerator) Name() string {
	return "primecodegen"
}

func (n *PrimeCodeGenGenerator) GetOutputName() string {
	return n.OutputName
}

func (n *PrimeCodeGenGenerator) Generate(opts GenerateOptions) error {
	// create dir
	_ = os.MkdirAll(opts.OutputDirectory, os.ModePerm)

	// generate
	err := n.generateCode(opts)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

func (n *PrimeCodeGenGenerator) generateCode(opts GenerateOptions) error {
	// primecodegen bin and args
	executable := "primecodegen"
	args := []string{
		"--log-level", "trace",
		"openapi-generate",
		"-i", n.APISpec,
		"-g", n.Config.TemplateLanguage,
		"-t", n.Config.TemplateType,
		"-o", opts.OutputDirectory,
		"--md-group-id", n.Config.GroupId,
		"--md-artifact-id", n.Config.ArtifactId,
	}
	for _, p := range n.Config.Patches {
		args = append(args, "--patches", p)
	}

	allArgs := append(args, n.Args...)
	cmd := exec.Command(executable, allArgs...)
	cmd.Dir = opts.ProjectDirectory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Trace().Str("command", cmd.String()).Msg("executing code generation")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute code generation: %w", err)
	}

	return nil
}
