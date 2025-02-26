package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type JavaLibraryGenerator struct {
	APISpec     string                     `json:"-" yaml:"-"`
	Repository  config.Repository          `json:"-" yaml:"-"`
	Maintainers []config.Maintainer        `json:"-" yaml:"-"`
	Opts        config.JavaLanguageOptions `json:"-" yaml:"-"`
}

func (n *JavaLibraryGenerator) Name() string {
	return "java-httpclient"
}

func (n *JavaLibraryGenerator) GetOutputName() string {
	return "java"
}

func (n *JavaLibraryGenerator) Generate(opts generator.GenerateOptions) error {
	log.Info().Str("dir", opts.OutputDirectory).Str("spec", n.APISpec).Msg("generating java library")

	gen := generator.PrimeCodeGenGenerator{
		OutputName: n.GetOutputName(),
		APISpec:    n.APISpec,
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

	return gen.Generate(opts)
}
