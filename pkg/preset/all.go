package preset

import (
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/generator"
	"github.com/primelib/primelib-app/pkg/util"
)

func Generators(specFile string, conf config.Configuration) []generator.Generator {
	var generators []generator.Generator

	// presets
	generators = addGeneratorIfEnabled(generators, conf.Presets.Java.Enabled, &JavaLibraryGenerator{
		APISpec:     specFile,
		Repository:  conf.Repository,
		Maintainers: conf.Maintainers,
		Opts:        conf.Presets.Java,
	})
	generators = addGeneratorIfEnabled(generators, conf.Presets.Go.Enabled, &GoLibraryGenerator{
		APISpec:     specFile,
		Repository:  conf.Repository,
		Maintainers: conf.Maintainers,
		Opts:        conf.Presets.Go,
	})
	generators = addGeneratorIfEnabled(generators, conf.Presets.Python.Enabled, &PythonLibraryGenerator{
		APISpec:     specFile,
		Repository:  conf.Repository,
		Maintainers: conf.Maintainers,
		Opts:        conf.Presets.Python,
	})
	generators = addGeneratorIfEnabled(generators, conf.Presets.CSharp.Enabled, &CSharpLibraryGenerator{
		APISpec:     specFile,
		Repository:  conf.Repository,
		Maintainers: conf.Maintainers,
		Opts:        conf.Presets.CSharp,
	})
	generators = addGeneratorIfEnabled(generators, conf.Presets.Typescript.Enabled, &TypeScriptLibraryGenerator{
		APISpec:     specFile,
		Repository:  conf.Repository,
		Maintainers: conf.Maintainers,
		Opts:        conf.Presets.Typescript,
	})

	// custom generators
	for _, g := range conf.Generators {
		var gen generator.Generator
		switch g.Type {
		case config.GeneratorTypeOpenApiGenerator:
			gen = &generator.OpenAPIGenerator{
				OutputName: g.Name,
				APISpec:    specFile,
				Args:       g.Arguments,
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
				OutputName: g.Name,
				APISpec:    specFile,
				Args:       g.Arguments,
				Config: generator.PrimeCodeGenGeneratorConfig{
					TemplateLanguage: util.GetMapString(g.Config, "templateLanguage", ""),
					TemplateType:     util.GetMapString(g.Config, "templateType", ""),
					Patches:          util.GetMapSliceString(g.Config, "patches", []string{}),
					GroupId:          util.GetMapString(g.Config, "groupId", ""),
					ArtifactId:       util.GetMapString(g.Config, "artifactId", ""),
				},
			}
		default:
			continue
		}

		addGeneratorIfEnabled(generators, g.Enabled, gen)
	}

	return generators
}

func addGeneratorIfEnabled(generators []generator.Generator, enabled bool, gen generator.Generator) []generator.Generator {
	if enabled {
		return append(generators, gen)
	}
	return generators
}
