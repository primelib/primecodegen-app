package primelib

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primelib-app/pkg/util"
	"github.com/rs/zerolog/log"
)

func Execute(dir string, module Module, repository api.Repository) error {
	specFile := filepath.Join(dir, module.Dir, module.SpecFile)
	log.Debug().Str("module", module.Name).Str("spec_url", module.SpecURL).Str("spec-file", specFile).Msg("processing module")

	// download spec file
	if module.SpecURL != "" {
		err := util.DownloadFile(module.SpecURL, specFile)
		if err != nil {
			return fmt.Errorf("failed to download spec file: %w", err)
		}
		log.Info().Str("file", specFile).Msg("downloaded spec file")
	}

	// download spec sources
	if len(module.SpecSources) > 0 {
		mergedData := make(map[string]interface{})

		for _, source := range module.SpecSources {
			content, err := util.DownloadString(source.URL)
			if err != nil {
				return fmt.Errorf("failed to download spec source: %w", err)
			}

			var jsonData map[string]interface{}
			err = json.Unmarshal([]byte(content), &jsonData)
			if err != nil {
				return fmt.Errorf("failed to parse spec source: %w", err)
			}

			log.Info().Str("source", source.Name).Str("url", source.URL).Msg("downloaded spec from source")
			util.MergeJSON(mergedData, jsonData)
		}

		// marshal merged data
		bytes, err := json.MarshalIndent(mergedData, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal merged data: %w", err)
		}

		// write to file
		err = os.WriteFile(specFile, bytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to write merged data to file: %w", err)
		}
	}

	// patch spec file
	if module.SpecScript != "" {
		err := patchSpecFile(filepath.Join(dir, module.Dir, module.SpecFile), module.SpecScript)
		if err != nil {
			return fmt.Errorf("failed to patch spec file: %w", err)
		}
		log.Debug().Str("file", specFile).Msg("patched openapi spec")
	}

	// delete generated files
	err := deleteGeneratedFiles(filepath.Join(dir, module.Dir))
	if err != nil {
		return fmt.Errorf("failed to delete generated files: %w", err)
	}

	// regenerate code
	err = generateCode(specFile, filepath.Join(dir, module.Dir), module.Config, repository)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	return nil
}

func patchSpecFile(specFile string, patch string) error {
	// create bash script in temp dir (patch is the script content)
	tempFile, err := os.CreateTemp("", "primelib-script")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// write content to file
	_, err = tempFile.WriteString(patch)
	if err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// execute script
	cmd := exec.Command("bash", tempFile.Name(), specFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return nil
}

// deleteGeneratedFiles removes generated files before running the next code generation to ensure that no old files are left
func deleteGeneratedFiles(moduleDirectory string) error {
	// check if .openapi-generator/FILES exists
	filesDir := filepath.Join(moduleDirectory, ".openapi-generator", "FILES")
	if _, err := os.Stat(filesDir); os.IsNotExist(err) {
		return nil
	}

	// get file list
	bytes, err := os.ReadFile(filesDir)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", filesDir, err)
	}

	// iterate over files
	files := strings.Split(string(bytes), "\n")
	log.Info().Int("count", len(files)).Msg("deleting generated files")
	for _, file := range files {
		// skip empty lines
		if file == "" {
			continue
		}

		log.Trace().Str("path", filepath.Join(moduleDirectory, file)).Msg("deleting file")
		err = os.Remove(filepath.Join(moduleDirectory, file))
		if err != nil {
			return fmt.Errorf("failed to delete file %s: %w", file, err)
		}
	}

	return nil
}

func generateCode(specFile string, moduleDirectory string, config GeneratorConfig, repository api.Repository) error {
	// auto generate config
	tempConfigFile, tmpErr := os.CreateTemp("", "openapi-generator.json")
	if tmpErr != nil {
		return fmt.Errorf("failed to create temp file: %w", tmpErr)
	}
	defer tempConfigFile.Close()

	// config
	configFile := path.Join(moduleDirectory, "openapi-generator.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// set defaults and missing properties
		config.EnablePostProcessFile = true
		if _, ok := config.AdditionalProperties["projectName"]; !ok {
			config.AdditionalProperties["projectName"] = repository.Name
		}
		if _, ok := config.AdditionalProperties["projectDescription"]; !ok {
			config.AdditionalProperties["projectDescription"] = repository.Description
		}
		if _, ok := config.AdditionalProperties["projectRepository"]; !ok {
			config.AdditionalProperties["projectRepository"] = repository.URL
		}
		if _, ok := config.AdditionalProperties["projectInceptionYear"]; !ok {
			config.AdditionalProperties["projectInceptionYear"] = repository.CreatedAt.Year()
		}
		if _, ok := config.AdditionalProperties["projectLicenseName"]; !ok {
			config.AdditionalProperties["projectLicenseName"] = repository.LicenseName
		}
		if _, ok := config.AdditionalProperties["projectLicenseUrl"]; !ok {
			config.AdditionalProperties["projectLicenseUrl"] = repository.LicenseURL
		}

		// marshal config
		bytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}

		// write to temp file
		err = os.WriteFile(tempConfigFile.Name(), bytes, 0644)
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}

		configFile = tempConfigFile.Name()
	}

	// generate code
	args := []string{
		"primecodegen",
		"generate",
		"-e", "auto",
		"-i", specFile,
		"-o", moduleDirectory,
		"-c", configFile,
		"--openapi-normalizer", "SIMPLIFY_ONEOF_ANYOF=true",
		"--openapi-normalizer", "SIMPLIFY_BOOLEAN_ENUM=true",
		"--openapi-normalizer", "REMOVE_ANYOF_ONEOF_AND_KEEP_PROPERTIES_ONLY=true",
		"--openapi-normalizer", "REFACTOR_ALLOF_WITH_PROPERTIES_ONLY=true",
		"--skip-validate-spec",
	}

	cmd := exec.Command("bash", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute code generation: %w", err)
	}

	return nil
}
