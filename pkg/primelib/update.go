package primelib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cidverse/go-vcsapp/pkg/platform/api"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/primelib/primelib-app/pkg/openapi"
	"github.com/primelib/primelib-app/pkg/util"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

// Update will update the openapi spec and apply patches
func Update(dir string, conf config.Configuration, repository api.Repository) error {
	spec := conf.Spec
	specFile := filepath.Join(dir, conf.Spec.File)
	log.Debug().Strs("spec-urls", spec.UrlSlice()).Str("spec-format", string(spec.Format)).Str("spec-file", specFile).Msg("processing module")

	// download spec sources
	if len(spec.Urls) > 0 {
		// fetch and merge specs
		mergedData := make(map[string]interface{})
		var fetchErr error
		for _, s := range spec.Urls {
			mergedData, fetchErr = fetchSpec(s, mergedData)
			if fetchErr != nil {
				return fetchErr
			}
		}

		// marshal merged data
		bytes, err := yaml.Marshal(mergedData)
		if err != nil {
			return fmt.Errorf("failed to marshal api spec: %w", err)
		}

		// open document and apply patches
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

		// write to file
		err = os.WriteFile(specFile, output, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write api spec to file: %w", err)
		}
	}

	return nil
}

// fetchSpec will download the spec from the source and merge it into the output
func fetchSpec(source config.SpecSource, output map[string]interface{}) (map[string]interface{}, error) {
	var sourceData map[string]interface{}
	if source.Format == "" || source.Format == config.SourceTypeSpec {
		content, err := util.DownloadString(source.URL)
		if err != nil {
			return output, fmt.Errorf("failed to download spec source: %w", err)
		}

		if strings.HasSuffix(source.URL, ".json") {
			err = json.Unmarshal([]byte(content), &sourceData)
			if err != nil {
				return output, fmt.Errorf("failed to parse json spec: %w", err)
			}
		} else if strings.HasSuffix(source.URL, ".yaml") {
			err = yaml.Unmarshal([]byte(content), &sourceData)
			if err != nil {
				return output, fmt.Errorf("failed to parse yaml spec: %w", err)
			}
		} else {
			err = yaml.Unmarshal([]byte(content), &sourceData)
			if err != nil {
				err = json.Unmarshal([]byte(content), &sourceData)
				if err != nil {
					return output, fmt.Errorf("failed to parse spec (attempts: yaml, json): %w", err)
				}
			}
		}
	} else if source.Format == config.SourceTypeSwaggerUI {
		swaggerJsUrl := source.URL + "/swagger-ui-init.js"
		content, err := util.DownloadString(swaggerJsUrl)
		if err != nil {
			return output, fmt.Errorf("failed to download spec source: %w", err)
		}

		// extract spec
		re := regexp.MustCompile(`"swaggerDoc":([\S\s]*),[\n\s]*"customOptions"`)
		match := re.FindStringSubmatch(content)
		if len(match) < 2 {
			return output, fmt.Errorf("failed to extract spec from swagger-ui-init.js")
		}

		// unmarshal
		err = json.Unmarshal([]byte(match[1]), &sourceData)
		if err != nil {
			return output, fmt.Errorf("failed to parse swagger-ui embedded json spec: %w", err)
		}
	}

	log.Info().Str("url", source.URL).Str("type", string(source.Format)).Msg("downloaded spec from source")
	util.MergeMaps(output, sourceData)
	return output, nil
}
