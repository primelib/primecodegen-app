package preset

import (
	"path/filepath"
	"strings"

	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type GoLibraryGenerator struct {
	APISpec     string                   `json:"-" yaml:"-"`
	Repository  config.Repository        `json:"-" yaml:"-"`
	Maintainers []config.Maintainer      `json:"-" yaml:"-"`
	Opts        config.GoLanguageOptions `json:"-" yaml:"-"`
}

func (n *GoLibraryGenerator) Name() string {
	return "go-httpclient"
}

func (n *GoLibraryGenerator) GetOutputName() string {
	return "go"
}

func (n *GoLibraryGenerator) Generate(opts generator.GenerateOptions) error {
	moduleName := n.Opts.ModuleName
	if moduleName == "" {
		moduleName = strings.TrimPrefix(n.Repository.URL, "https://")

		relPath, err := filepath.Rel(opts.ProjectDirectory, opts.OutputDirectory)
		if err == nil && relPath != "." {
			moduleName = filepath.Join(moduleName, relPath)
		}
	}

	log.Info().Str("dir", opts.OutputDirectory).Str("spec", n.APISpec).Msg("generating go library")
	gen := generator.PrimeCodeGenGenerator{
		OutputName: n.GetOutputName(),
		APISpec:    n.APISpec,
		Args: []string{
			"--md-artifact-id", moduleName,
		},
		Config: generator.PrimeCodeGenGeneratorConfig{
			TemplateLanguage: "go",
			TemplateType:     "httpclient",
			Patches:          []string{},
		},
	}

	return gen.Generate(opts)
}
