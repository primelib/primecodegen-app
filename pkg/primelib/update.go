package primelib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primecodegen-app/pkg/config"
	"github.com/primelib/primecodegen-app/pkg/openapi"
	"github.com/primelib/primecodegen-app/pkg/specutil"
	"github.com/primelib/primecodegen-app/pkg/util"
	"github.com/rs/zerolog/log"
)

// Update will update the openapi spec and apply patches
func Update(dir string, conf config.Configuration, repository api.Repository) error {
	spec := conf.Spec
	specFile := filepath.Join(dir, conf.Spec.File)
	log.Debug().Strs("spec-urls", spec.UrlSlice()).Str("spec-format", string(spec.Type)).Str("spec-file", specFile).Msg("processing module")

	// download spec sources
	targetSpecDir := spec.GetSourcesDir(dir)
	var specFiles []string
	var specFilesType []config.SpecType
	var tempFiles []string
	defer func() {
		for _, f := range tempFiles {
			_ = os.Remove(f)
		}
	}()

	// download spec sources
	for _, s := range spec.Sources {
		log.Debug().Str("url", s.URL).Str("type", string(s.Type)).Msg("fetching spec")
		var targetFile string
		var bytes []byte
		var err error

		// fetch spec
		if s.File != "" && s.URL == "" {
			bytes, err = os.ReadFile(filepath.Join(targetSpecDir, s.File))
		} else if s.URL != "" {
			bytes, err = fetchSpec(s)
		}
		if err != nil {
			return fmt.Errorf("failed to fetch spec: %w", err)
		}

		if s.File != "" && s.URL != "" {
			targetFile = filepath.Join(targetSpecDir, s.File)
		} else {
			tempFile, err := os.CreateTemp("", "api-spec-*.yaml")
			if err != nil {
				return fmt.Errorf("failed to create temp file: %w", err)
			}
			tempFiles = append(tempFiles, tempFile.Name())
			targetFile = tempFile.Name()
		}

		// write to file
		err = os.WriteFile(targetFile, bytes, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write api spec to file: %w", err)
		}
		specFiles = append(specFiles, targetFile)
		specFilesType = append(specFilesType, s.Type)
	}

	// spec type conversions
	for i, f := range specFiles {
		// convert from swagger to openapi
		if spec.Type == config.SpecTypeOpenAPI3 && specFilesType[i] == config.SpecTypeSwagger2 {
			log.Debug().Str("file", f).Msg("converting from swagger to openapi")
			err := specutil.ConvertSwaggerToOpenAPI(f)
			if err != nil {
				return fmt.Errorf("failed to convert swagger to openapi: %w", err)
			}
		}
	}

	// openapi processing
	if spec.Type == config.SpecTypeOpenAPI3 {
		// merge and patch
		log.Debug().Strs("files", specFiles).Str("output", specFile).Msg("merging and patching openapi spec")
		_ = os.Remove(specFile)
		err := specutil.MergeAndPatchOpenAPI(specFiles, spec.InputPatches, spec.Patches, specFile)
		if err != nil {
			return err
		}

		// apply customizations
		log.Debug().Str("file", specFile).Msg("applying customizations")
		bytes, err := os.ReadFile(specFile)
		if err != nil {
			return err
		}

		doc, err := openapi.OpenDocument(bytes)
		if err != nil {
			return fmt.Errorf("failed to open document: %w", err)
		}
		specInfo := doc.GetSpecInfo()
		doc = openapi.PatchDocument(doc, specInfo.SpecType, specInfo.SpecFormat, specInfo.VersionNumeric, conf.Spec.Customization)
		output, err := doc.Render()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to render document")
		}

		err = os.WriteFile(specFile, output, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write api spec to file: %w", err)
		}
	}

	return nil
}

// fetchSpec will download the spec from the source and merge it into the output
func fetchSpec(source config.SpecSource) ([]byte, error) {
	if source.Format == "" || source.Format == config.SourceTypeSpec {
		content, err := util.DownloadString(source.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to download spec source: %w", err)
		}
		return content, nil
	} else if source.Format == config.SourceTypeSwaggerUI {
		swaggerJsUrl := source.URL + "/swagger-ui-init.js"
		content, err := util.DownloadString(swaggerJsUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to download spec source: %w", err)
		}

		// extract spec
		re := regexp.MustCompile(`"swaggerDoc":([\S\s]*),[\n\s]*"customOptions"`)
		match := re.FindStringSubmatch(string(content))
		if len(match) < 2 {
			return nil, fmt.Errorf("failed to extract spec from swagger-ui-init.js")
		}

		return []byte(match[1]), nil
	}

	return nil, fmt.Errorf("unsupported source type: %s", source.Format)
}
