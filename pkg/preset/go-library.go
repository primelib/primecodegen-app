package preset

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/primelib/primecodegen-app/pkg/config"
	"github.com/primelib/primecodegen-app/pkg/generator"
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
	moduleName := suggestGoModuleName(n.Opts.ModuleName, n.Repository, opts.ProjectDirectory, opts.OutputDirectory)

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

func suggestGoModuleName(moduleName string, repository config.Repository, projectDirectory string, outputDirectory string) string {
	if moduleName != "" {
		return moduleName
	}

	// trim protocol prefix
	parsedURL, err := url.Parse(repository.URL)
	if err != nil {
		return "example.com/unknown-module"
	}
	moduleName = parsedURL.Host + parsedURL.Path

	// append relative path in case output directory is not the project directory
	relPath, err := filepath.Rel(projectDirectory, outputDirectory)
	if err == nil && relPath != "." {
		moduleName = filepath.Join(moduleName, relPath)
	}

	// replace all backslashes with forward slashes
	moduleName = strings.ReplaceAll(moduleName, "\\", "/")

	return moduleName
}
