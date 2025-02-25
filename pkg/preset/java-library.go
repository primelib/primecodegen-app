package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type JavaLibraryGenerator struct {
	Directory   string                     `json:"-" yaml:"-"`
	APISpec     string                     `json:"-" yaml:"-"`
	Repository  config.Repository          `json:"-" yaml:"-"`
	Maintainers []config.Maintainer        `json:"-" yaml:"-"`
	Opts        config.JavaLanguageOptions `json:"-" yaml:"-"`
}

// Name returns the name of the task
func (n *JavaLibraryGenerator) Name() string {
	return "java-httpclient"
}

func (n *JavaLibraryGenerator) SetOutputDirectory(dir string) {
	n.Directory = dir
}

func (n *JavaLibraryGenerator) GetOutputDirectory() string {
	return n.Directory
}

func (n *JavaLibraryGenerator) Generate() error {
	log.Info().Str("dir", n.Directory).Str("spec", n.APISpec).Msg("generating java library")

	gen := generator.PrimeCodeGenGenerator{
		Directory: n.Directory,
		APISpec:   n.APISpec,
		Args: []string{
			"--md-group-id", n.Opts.GroupId,
			"--md-artifact-id", n.Opts.ArtifactId,
		},
		Config: generator.PrimeCodeGenGeneratorConfig{
			TemplateLanguage: "java",
			TemplateType:     "httpclient",
			Patches:          []string{},
		},
	}

	return gen.Generate()
}
