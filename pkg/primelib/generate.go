package primelib

import (
	"fmt"
	"path/filepath"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/rs/zerolog/log"
)

func Generate(dir string, conf config.Configuration, repository api.Repository) error {
	spec := conf.Spec
	specFile := filepath.Join(dir, conf.Spec.File)
	log.Debug().Strs("spec_urls", spec.Urls).Str("spec-file", specFile).Msg("processing module")

	/*
		if _, fileErr := os.Stat(configFile); os.IsNotExist(fileErr) {
				// set defaults and missing properties
				n.EnablePostProcessFile = true
				if _, ok := n.AdditionalProperties["projectName"]; !ok {
					n.AdditionalProperties["projectName"] = repository.Name
				}
				if _, ok := n.AdditionalProperties["projectDescription"]; !ok {
					n.AdditionalProperties["projectDescription"] = repository.Description
				}
				if _, ok := n.AdditionalProperties["projectRepository"]; !ok {
					n.AdditionalProperties["projectRepository"] = repository.URL
				}
				if _, ok := n.AdditionalProperties["projectInceptionYear"]; !ok {
					if repository.CreatedAt != nil {
						n.AdditionalProperties["projectInceptionYear"] = repository.CreatedAt.Year()
					}
				}
				if _, ok := n.AdditionalProperties["projectLicenseName"]; !ok {
					n.AdditionalProperties["projectLicenseName"] = repository.LicenseName
				}
				if _, ok := n.AdditionalProperties["projectLicenseUrl"]; !ok {
					n.AdditionalProperties["projectLicenseUrl"] = repository.LicenseURL
				}
	*/

	// generate code
	var generators []generator.Generator
	addGenerator := func(enabled bool, langDir string, gen generator.Generator) {
		if enabled {
			targetDir := filepath.Join(dir, conf.Output)
			if conf.Generators.MultiLanguage() {
				targetDir = filepath.Join(targetDir, langDir)
			}
			gen.SetOutputDirectory(targetDir)
			generators = append(generators, gen)
		}
	}
	addGenerator(conf.Generators.Java.Enabled, "java", &generator.JavaLibraryGenerator{
		APISpec: specFile,
		Opts:    conf.Generators.Java,
	})
	addGenerator(conf.Generators.Go.Enabled, "go", &generator.GoLibraryGenerator{
		APISpec: specFile,
		Opts:    conf.Generators.Go,
	})
	for _, gen := range generators {
		log.Debug().Str("generator", gen.Name()).Msg("executing generator")
		err := gen.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate code: %w", err)
		}
	}

	return nil
}
