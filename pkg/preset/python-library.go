package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

type PythonLibraryGenerator struct {
	Directory   string                       `json:"-" yaml:"-"`
	APISpec     string                       `json:"-" yaml:"-"`
	Repository  config.Repository            `json:"-" yaml:"-"`
	Maintainers []config.Maintainer          `json:"-" yaml:"-"`
	Opts        config.PythonLanguageOptions `json:"-" yaml:"-"`
}

// Name returns the name of the task
func (n *PythonLibraryGenerator) Name() string {
	return "python-httpclient"
}

func (n *PythonLibraryGenerator) SetOutputDirectory(dir string) {
	n.Directory = dir
}

func (n *PythonLibraryGenerator) GetOutputDirectory() string {
	return n.Directory
}

func (n *PythonLibraryGenerator) Generate() error {
	log.Info().Str("dir", n.Directory).Str("spec", n.APISpec).Msg("generating python library")

	gen := generator.OpenAPIGenerator{
		Directory: n.Directory,
		APISpec:   n.APISpec,
		Config: generator.OpenAPIGeneratorConfig{
			GeneratorName:         "python",
			EnablePostProcessFile: false,
			GlobalProperty:        nil,
			AdditionalProperties: map[string]interface{}{
				"library":        "urllib3",
				"projectName":    n.Repository.Name,
				"packageName":    n.Opts.PypiPackageName,
				"packageUrl":     n.Repository.URL,
				"packageVersion": "",
			},
			IgnoreFiles: []string{
				"README.md",
				"tox.ini",
				".travis.yml",
				".gitlab-ci.yml",
				".gitignore",
				"git_push.sh",
				".github/*",
				"docs/*",
			},
		},
	}

	return gen.Generate()
}
