package primelib

import (
	"fmt"
	"path/filepath"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primecodegen-app/pkg/config"
	"github.com/primelib/primecodegen-app/pkg/generator"
	"github.com/primelib/primecodegen-app/pkg/preset"
	"github.com/rs/zerolog/log"
)

func Generate(dir string, conf config.Configuration, repository api.Repository) error {
	spec := conf.Spec
	specFile := filepath.Join(dir, conf.Spec.File)
	log.Debug().Strs("spec-urls", spec.UrlSlice()).Str("spec-file", specFile).Msg("processing module")

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

	// prepare generators
	generators := preset.Generators(specFile, conf)

	// execute generators
	for _, gen := range generators {
		outputDir := filepath.Join(dir, conf.Output)
		if conf.MultiLanguage() {
			outputDir = filepath.Join(outputDir, gen.GetOutputName())
		}

		log.Info().Str("generator", gen.Name()).Str("projectDir", dir).Str("outputDir", outputDir).Msg("running code generator")
		err := gen.Generate(generator.GenerateOptions{
			ProjectDirectory: dir,
			OutputDirectory:  outputDir,
		})
		if err != nil {
			return fmt.Errorf("failed to generate code: %w", err)
		}
		log.Info().Str("generator", gen.Name()).Msg("code generation completed")
	}

	return nil
}
