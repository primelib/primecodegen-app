package primelib

import (
	"fmt"
	"path/filepath"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/primelib/primelib-app/pkg/preset"
	"github.com/primelib/primelib-app/pkg/util"
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

	// generate code
	var generators []generator.Generator
	addGenerator := func(enabled bool, langDir string, gen generator.Generator) {
		if enabled {
			targetDir := filepath.Join(dir, conf.Output)
			if conf.MultiLanguage() {
				targetDir = filepath.Join(targetDir, langDir)
			}
			gen.SetOutputDirectory(targetDir)
			generators = append(generators, gen)
		}
	}

	// presets
	addGenerator(conf.Presets.Java.Enabled, "java", &preset.JavaLibraryGenerator{
		APISpec: specFile,
		Opts:    conf.Presets.Java,
	})
	addGenerator(conf.Presets.Go.Enabled, "go", &preset.GoLibraryGenerator{
		APISpec: specFile,
		Opts:    conf.Presets.Go,
	})

	// custom generators
	for _, g := range conf.Generators {
		var gen generator.Generator
		switch g.Type {
		case config.GeneratorTypeOpenApiGenerator:
			gen = &generator.OpenAPIGenerator{
				APISpec: specFile,
				Args:    g.Arguments,
				Config: generator.OpenAPIGeneratorConfig{
					GeneratorName:         util.GetMapString(g.Config, "generatorName", ""),
					InvokerPackage:        util.GetMapString(g.Config, "invokerPackage", ""),
					ApiPackage:            util.GetMapString(g.Config, "apiPackage", ""),
					ModelPackage:          util.GetMapString(g.Config, "modelPackage", ""),
					EnablePostProcessFile: util.GetMapBool(g.Config, "enablePostProcessFile", false),
					GlobalProperty:        util.GetMapMap(g.Config, "globalProperty"),
					AdditionalProperties:  util.GetMapMap(g.Config, "additionalProperties"),
				},
			}
		case config.GeneratorTypePrimeCodeGen:
			gen = &generator.PrimeCodeGenGenerator{
				APISpec: specFile,
				Args:    g.Arguments,
				Config: generator.PrimeCodeGenGeneratorConfig{
					TemplateLanguage: util.GetMapString(g.Config, "templateLanguage", ""),
					TemplateType:     util.GetMapString(g.Config, "templateType", ""),
					Patches:          util.GetMapSliceString(g.Config, "patches", []string{}),
					GroupId:          util.GetMapString(g.Config, "groupId", ""),
					ArtifactId:       util.GetMapString(g.Config, "artifactId", ""),
				},
			}
		}

		addGenerator(g.Enabled, g.Name, gen)
	}

	// execute generators
	for _, gen := range generators {
		log.Info().Str("generator", gen.Name()).Msg("running code generator")
		err := gen.Generate()
		if err != nil {
			return fmt.Errorf("failed to generate code: %w", err)
		}
		log.Info().Str("generator", gen.Name()).Msg("code generation completed")
	}

	return nil
}
