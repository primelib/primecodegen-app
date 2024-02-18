package primelib

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
	//projectDir := filepath.Join(dir)
	spec := conf.Spec
	specFile := filepath.Join(dir, conf.Spec.File)
	log.Debug().Strs("spec_urls", spec.Urls).Str("spec-file", specFile).Msg("processing module")

	// download spec sources
	if len(spec.Urls) > 0 {
		mergedData := make(map[string]interface{})

		for _, source := range spec.Urls {
			content, err := util.DownloadString(source)
			if err != nil {
				return fmt.Errorf("failed to download spec source: %w", err)
			}

			var sourceData map[string]interface{}
			if strings.HasSuffix(source, ".json") {
				err = json.Unmarshal([]byte(content), &sourceData)
				if err != nil {
					return fmt.Errorf("failed to parse json spec: %w", err)
				}
			} else if strings.HasSuffix(source, ".yaml") {
				err = yaml.Unmarshal([]byte(content), &sourceData)
				if err != nil {
					return fmt.Errorf("failed to parse yaml spec: %w", err)
				}
			}

			log.Info().Str("url", source).Msg("downloaded spec from source")
			util.MergeMaps(mergedData, sourceData)
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
