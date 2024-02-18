package generator

import (
	"github.com/primelib/primelib-app/pkg/config"
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
	return "openapi-generator"
}

func (n *JavaLibraryGenerator) SetOutputDirectory(dir string) {
	n.Directory = dir
}

func (n *JavaLibraryGenerator) GetOutputDirectory() string {
	return n.Directory
}

func (n *JavaLibraryGenerator) Generate() error {
	log.Info().Str("dir", n.Directory).Str("spec", n.APISpec).Msg("generating java library")
	gen := OpenAPIGenerator{
		Directory: n.Directory,
		APISpec:   n.APISpec,
		Args:      []string{},
		Config: OpenAPIGeneratorConfig{
			GeneratorName:         "primecodegen-java-feign",
			InvokerPackage:        n.Opts.GroupId,
			ApiPackage:            n.Opts.GroupId + ".api",
			ModelPackage:          n.Opts.GroupId + ".model",
			EnablePostProcessFile: true,
			GlobalProperty:        nil,
			AdditionalProperties: map[string]interface{}{
				"projectArtifactGroupId": n.Opts.GroupId,
				"projectArtifactId":      n.Opts.ArtifactId,
				"projectName":            n.Repository.Name,
				"projectDescription":     n.Repository.Description,
				"projectRepository":      n.Repository.URL,
				"projectInceptionYear":   n.Repository.InceptionYear,
				"projectLicenseName":     n.Repository.LicenseName,
				"projectLicenseUrl":      n.Repository.LicenseURL,
			},
		},
	}

	if len(n.Maintainers) > 0 {
		firstMaintainer := n.Maintainers[0]
		gen.Config.AdditionalProperties["projectMaintainerId"] = firstMaintainer.ID
		gen.Config.AdditionalProperties["projectMaintainerName"] = firstMaintainer.Name
		gen.Config.AdditionalProperties["projectMaintainerEMail"] = firstMaintainer.Email
	}

	return gen.Generate()
}
