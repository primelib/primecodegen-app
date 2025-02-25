package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type GoLibraryGenerator struct {
	Directory   string                   `json:"-" yaml:"-"`
	APISpec     string                   `json:"-" yaml:"-"`
	Repository  config.Repository        `json:"-" yaml:"-"`
	Maintainers []config.Maintainer      `json:"-" yaml:"-"`
	Opts        config.GoLanguageOptions `json:"-" yaml:"-"`
}

// Name returns the name of the task
func (n *GoLibraryGenerator) Name() string {
	return "go-httpclient"
}

func (n *GoLibraryGenerator) SetOutputDirectory(dir string) {
	n.Directory = dir
}

func (n *GoLibraryGenerator) GetOutputDirectory() string {
	return n.Directory
}

func (n *GoLibraryGenerator) Generate() error {
	log.Info().Str("dir", n.Directory).Str("spec", n.APISpec).Msg("generating go library")

	gen := generator.PrimeCodeGenGenerator{
		Directory: n.Directory,
		APISpec:   n.APISpec,
		Args: []string{
			"--md-artifact-id", n.Opts.ModuleName,
		},
		Config: generator.PrimeCodeGenGeneratorConfig{
			TemplateLanguage: "go",
			TemplateType:     "httpclient",
			Patches:          []string{},
		},
	}

	return gen.Generate()
}
