package generator

import (
	"github.com/primelib/primelib-app/pkg/config"
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
	return "openapi-generator"
}

func (n *GoLibraryGenerator) SetOutputDirectory(dir string) {
	n.Directory = dir
}

func (n *GoLibraryGenerator) GetOutputDirectory() string {
	return n.Directory
}

func (n *GoLibraryGenerator) Generate() error {
	log.Info().Str("dir", n.Directory).Str("spec", n.APISpec).Msg("generating go library")
	gen := OpenAPIGenerator{
		Directory: n.Directory,
		APISpec:   n.APISpec,
		Args:      []string{},
		Config: OpenAPIGeneratorConfig{
			GeneratorName: "go",
			//InvokerPackage:        n.Opts.GroupId,
			//ApiPackage:            n.Opts.GroupId + ".api",
			//ModelPackage:          n.Opts.GroupId + ".model",
			EnablePostProcessFile: true,
			GlobalProperty:        nil,
			AdditionalProperties: map[string]interface{}{
				//"projectArtifactGroupId": n.Opts.GroupId,
				//"projectArtifactId":      n.Opts.ArtifactId,
				"projectName":          n.Repository.Name,
				"projectDescription":   n.Repository.Description,
				"projectRepository":    n.Repository.URL,
				"projectInceptionYear": n.Repository.InceptionYear,
				"projectLicenseName":   n.Repository.LicenseName,
				"projectLicenseUrl":    n.Repository.LicenseURL,
			},
		},
	}

	return gen.Generate()
}
