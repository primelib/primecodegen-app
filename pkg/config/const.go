package config

// ConfigFileName is the default name of the configuration file
const ConfigFileName = "primelib.yaml"

type GeneratorType string

const (
	GeneratorTypeOpenApiGenerator GeneratorType = "openapi-generator"
	GeneratorTypePrimeCodeGen     GeneratorType = "primecodegen"
)

type SourceType string

const (
	SourceTypeSpec      SourceType = "spec"
	SourceTypeSwaggerUI SourceType = "swagger-ui"
)

type SpecType string

const (
	SpecTypeOpenAPI3 SpecType = "openapi3"
	SpecTypeSwagger2 SpecType = "swagger2"
)
