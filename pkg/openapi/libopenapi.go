package openapi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/primelib/primelib-app/pkg/config"
	"github.com/rs/zerolog/log"
)

func OpenDocumentFile(file string) (libopenapi.Document, error) {
	// read the file
	input, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	// config
	conf := datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
		BasePath:              filepath.Dir(file),
		//BaseURL:               baseURL,
	}

	// create a new document from specification bytes
	document, err := libopenapi.NewDocumentWithConfiguration(input, &conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create document from spec")
	}

	return document, nil
}

func OpenDocument(input []byte) (libopenapi.Document, error) {
	// config
	conf := datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
	}

	// create a new document from specification bytes
	document, err := libopenapi.NewDocumentWithConfiguration(input, &conf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create document from spec")
	}

	return document, nil
}

func PatchDocument(document libopenapi.Document, specType string, specFormat string, specVersion float32, customizations config.Customization) libopenapi.Document {
	// detect type
	if specType == "openapi" && specFormat == "oas3" {
		v3doc, errors := document.BuildV3Model()
		if len(errors) > 0 {
			log.Fatal().Errs("errors", errors).Msg("cannot create v3 model from document")
		}
		model := &v3doc.Model

		// info
		if customizations.Title != "" {
			model.Info.Title = customizations.Title
		}
		if customizations.Summary != "" && specVersion >= 3.1 {
			model.Info.Summary = customizations.Summary
		}
		if customizations.Description != "" {
			model.Info.Description = customizations.Description
		}

		if customizations.Version != "" {
			model.Info.Version = customizations.Version
		}
		if customizations.Contact.Name != "" || customizations.Contact.URL != "" || customizations.Contact.Email != "" {
			model.Info.Contact = &base.Contact{
				Name:  customizations.Contact.Name,
				URL:   customizations.Contact.URL,
				Email: customizations.Contact.Email,
			}
		}
		if customizations.License.Name != "" || customizations.License.URL != "" || customizations.License.Identifier != "" {
			model.Info.License = &base.License{
				Name:       customizations.License.Name,
				URL:        customizations.License.URL,
				Identifier: customizations.License.Identifier,
			}
		}

		// servers
		if len(customizations.Servers) > 0 {
			model.Servers = []*v3.Server{}
			for _, server := range customizations.Servers {
				model.Servers = append(model.Servers, &v3.Server{
					URL:         server.URL,
					Description: server.Description,
					Variables:   nil,
				})
			}
		}
	}

	return document
}
