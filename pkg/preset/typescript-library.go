package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type TypeScriptLibraryGenerator struct {
	APISpec     string                           `json:"-" yaml:"-"`
	Repository  config.Repository                `json:"-" yaml:"-"`
	Maintainers []config.Maintainer              `json:"-" yaml:"-"`
	Opts        config.TypescriptLanguageOptions `json:"-" yaml:"-"`
}

func (n *TypeScriptLibraryGenerator) Name() string {
	return "typescript-httpclient"
}

func (n *TypeScriptLibraryGenerator) GetOutputName() string {
	return "typescript"
}

func (n *TypeScriptLibraryGenerator) Generate(opts generator.GenerateOptions) error {
	log.Info().Str("dir", opts.OutputDirectory).Str("spec", n.APISpec).Msg("generating python library")

	gen := generator.OpenAPIGenerator{
		OutputName: n.GetOutputName(),
		APISpec:    n.APISpec,
		Config: generator.OpenAPIGeneratorConfig{
			GeneratorName:         "typescript-axios",
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
		},
	}

	return gen.Generate(opts)
}
