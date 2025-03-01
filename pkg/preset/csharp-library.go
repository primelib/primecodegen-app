package preset

import (
	"github.com/primelib/primecodegen-app/pkg/config"
	"github.com/primelib/primecodegen-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type CSharpLibraryGenerator struct {
	APISpec     string                       `json:"-" yaml:"-"`
	Repository  config.Repository            `json:"-" yaml:"-"`
	Maintainers []config.Maintainer          `json:"-" yaml:"-"`
	Opts        config.CSharpLanguageOptions `json:"-" yaml:"-"`
}

func (n *CSharpLibraryGenerator) Name() string {
	return "csharp-httpclient"
}

func (n *CSharpLibraryGenerator) GetOutputName() string {
	return "csharp"
}

func (n *CSharpLibraryGenerator) Generate(opts generator.GenerateOptions) error {
	log.Info().Str("dir", opts.OutputDirectory).Str("spec", n.APISpec).Msg("generating csharp library")

	gen := generator.OpenAPIGenerator{
		OutputName: n.GetOutputName(),
		APISpec:    n.APISpec,
		Config: generator.OpenAPIGeneratorConfig{
			GeneratorName:         "csharp",
			EnablePostProcessFile: false,
			GlobalProperty:        nil,
			AdditionalProperties: map[string]interface{}{
				"projectName": n.Repository.Name,
			},
			IgnoreFiles: []string{
				"README.md",
				".travis.yml",
				"appveyor.yml",
				".gitlab-ci.yml",
				".gitignore",
				"git_push.sh",
				".github/*",
				"docs/*",
			},
			Repository:  n.Repository,
			Maintainers: n.Maintainers,
		},
	}

	return gen.Generate(opts)
}
