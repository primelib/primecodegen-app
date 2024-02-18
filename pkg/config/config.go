package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Name        string `yaml:"name"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	Output      string `yaml:"output" jsonschema_description:"output directory for the generated code"`

	Repository  Repository   `yaml:"repository"`
	Maintainers []Maintainer `yaml:"maintainers"`
	Generators  Generators   `yaml:"generator"`

	Spec Spec `yaml:"spec"`
}

type Repository struct {
	Name          string `yaml:"name"`
	Description   string `yaml:"description"`
	URL           string `yaml:"url"`
	InceptionYear int    `yaml:"inceptionYear"`
	LicenseName   string `yaml:"licenseName"`
	LicenseURL    string `yaml:"licenseURL"`
}

type Maintainer struct {
	ID    string `yaml:"id"`
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
	URL   string `yaml:"url"`
}

type Generators struct {
	Go         GoLanguageOptions         `yaml:"go"`
	Java       JavaLanguageOptions       `yaml:"java"`
	Python     PythonLanguageOptions     `yaml:"python"`
	Typescript TypescriptLanguageOptions `yaml:"typescript"`
}

func (c Generators) EnabledCount() int {
	enabledCount := 0

	if c.Go.Enabled {
		enabledCount++
	}
	if c.Java.Enabled {
		enabledCount++
	}
	if c.Python.Enabled {
		enabledCount++
	}
	if c.Typescript.Enabled {
		enabledCount++
	}

	return enabledCount
}

func (c Generators) MultiLanguage() bool {
	return c.EnabledCount() >= 2
}

type GoLanguageOptions struct {
	Enabled     bool     `yaml:"enabled"`
	IgnoreFiles []string `yaml:"ignoreFiles"`
}

type JavaLanguageOptions struct {
	Enabled     bool     `yaml:"enabled"`
	IgnoreFiles []string `yaml:"ignoreFiles"`

	GroupId    string `yaml:"groupId"`
	ArtifactId string `yaml:"artifactId"`
}

type PythonLanguageOptions struct {
	Enabled     bool     `yaml:"enabled"`
	IgnoreFiles []string `yaml:"ignoreFiles"`

	PypiPackageName string `yaml:"pypiPackageName"`
}

type TypescriptLanguageOptions struct {
	Enabled     bool     `yaml:"enabled"`
	IgnoreFiles []string `yaml:"ignoreFiles"`

	NpmOrg  string `yaml:"npmOrg"`
	NpmName string `yaml:"npmName"`
}

type Spec struct {
	// File is the path to the openapi specification file
	File string `yaml:"file" default:"openapi.yaml" required:"true"`
	// Urls contains one or multiple urls to the openapi specifications, all documents will be merged
	Urls []string `yaml:"urls" required:"true"`
	// Format is the format of the api specification
	Format SpecType `yaml:"format" required:"true"`
	// Patches are the patches that are applied to the openapi specification
	Customization Customization `yaml:"customization"`
}

type SpecSource struct {
	Name string `yaml:"name" required:"true"`
	URL  string `yaml:"url" required:"true"`
}

type SpecType string

const (
	OpenAPI3 SpecType = "openapi3"
	Swagger2 SpecType = "swagger2"
)

type Customization struct {
	Title       string                `yaml:"title"`
	Summary     string                `yaml:"summary"`
	Description string                `yaml:"description"`
	Version     string                `yaml:"version"`
	Contact     CustomizationContact  `yaml:"contact"`
	License     CustomizationLicense  `yaml:"license"`
	Servers     []CustomizationServer `yaml:"servers"`

	// Prune operations, tags and schemas
	PruneOperations []string `yaml:"pruneOperations"`
	PruneTags       []string `yaml:"pruneTags"`
	PruneSchemas    []string `yaml:"pruneSchemas"`
}

type CustomizationContact struct {
	Name  string `yaml:"name"`
	URL   string `yaml:"url"`
	Email string `yaml:"email"`
}

type CustomizationLicense struct {
	Name       string `yaml:"name"`
	URL        string `yaml:"url"`
	Identifier string `yaml:"identifier"`
}

type CustomizationServer struct {
	URL         string `yaml:"url"`
	Description string `yaml:"description"`
}

type GeneratorConfig struct {
	GeneratorName         string                 `json:"generatorName" yaml:"generatorName"`
	InvokerPackage        string                 `json:"invokerPackage" yaml:"invokerPackage"`
	ApiPackage            string                 `json:"apiPackage" yaml:"apiPackage"`
	ModelPackage          string                 `json:"modelPackage" yaml:"modelPackage"`
	EnablePostProcessFile bool                   `json:"enablePostProcessFile" yaml:"enablePostProcessFile"`
	GlobalProperty        map[string]interface{} `json:"globalProperty" yaml:"globalProperty"`
	AdditionalProperties  map[string]interface{} `json:"additionalProperties" yaml:"additionalProperties"`
}

type GeneratorArgs struct {
	// UserArgs are the arguments that are passed to the generator
	OpenAPIGeneratorArgs []string `yaml:"openapi_generator"`
}

func LoadProjectConfig(file string) (Configuration, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to read %s: %w", file, err)
	}

	return ConfigFromString(string(bytes))
}

func ConfigFromString(content string) (Configuration, error) {
	var config Configuration
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return Configuration{}, fmt.Errorf("failed to parse config: %w", err)
	}

	// repository defaults
	if config.Repository.Name == "" {
		config.Repository.Name = config.Name
	}
	if config.Repository.Description == "" {
		config.Repository.Description = config.Summary
	}

	// spec defaults
	if config.Spec.Customization.Title == "" {
		config.Spec.Customization.Title = config.Name
	}
	if config.Spec.File == "" {
		config.Spec.File = "openapi.yaml"
	}

	return config, nil
}
