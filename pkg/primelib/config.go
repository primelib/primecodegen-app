package primelib

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Configuration struct {
	Modules []Module `yaml:"modules"`
}

type Module struct {
	// Name is the name of the module
	Name string `yaml:"name" required:"true"`
	// Dir is the relative path to the module
	Dir string `yaml:"dir" required:"true"`
	// SpecURL is the URL to the openapi spec
	SpecURL string `yaml:"spec_url"`
	// SpecFile is the relative path to the openapi spec
	SpecFile string `yaml:"spec_file" required:"true"`
	// SpecScript accepts a script that can be used to fix issues in the openapi spec
	SpecScript string `yaml:"spec_script"`
	// GenerateScript is the relative path to the script that generates the code
	SpecSources []SpecSource `yaml:"spec_sources"`
}

type SpecSource struct {
	Name string `yaml:"name" required:"true"`
	URL  string `yaml:"url" required:"true"`
}

func LoadProjectConfig(file string) (*Configuration, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", file, err)
	}

	var config Configuration
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", file, err)
	}

	return &config, nil
}

func ConfigFromString(content string) (*Configuration, error) {
	var config Configuration
	err := yaml.Unmarshal([]byte(content), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
