package primelib

// OpenApiGeneratorArgumentAllowList is a list of arguments that are allowed to be passed to the openapi generator
var openApiGeneratorArgumentAllowList = []string{
	// spec validation
	"--skip-validate-spec",
	// normalizer - see https://openapi-generator.tech/docs/customization/#openapi-normalizer
	"--openapi-normalizer",
	"SIMPLIFY_ANYOF_STRING_AND_ENUM_STRING=true",
	"SIMPLIFY_ANYOF_STRING_AND_ENUM_STRING=false",
	"SIMPLIFY_BOOLEAN_ENUM=true",
	"SIMPLIFY_BOOLEAN_ENUM=false",
	"SIMPLIFY_ONEOF_ANYOF=true",
	"SIMPLIFY_ONEOF_ANYOF=false",
	"ADD_UNSIGNED_TO_INTEGER_WITH_INVALID_MAX_VALUE=true",
	"ADD_UNSIGNED_TO_INTEGER_WITH_INVALID_MAX_VALUE=false",
	"REFACTOR_ALLOF_WITH_PROPERTIES_ONLY=true",
	"REFACTOR_ALLOF_WITH_PROPERTIES_ONLY=false",
	"REF_AS_PARENT_IN_ALLOF=true",
	"REF_AS_PARENT_IN_ALLOF=false",
	"REMOVE_ANYOF_ONEOF_AND_KEEP_PROPERTIES_ONLY=true",
	"REMOVE_ANYOF_ONEOF_AND_KEEP_PROPERTIES_ONLY=false",
	"KEEP_ONLY_FIRST_TAG_IN_OPERATION=true",
	"KEEP_ONLY_FIRST_TAG_IN_OPERATION=false",
	"SET_TAGS_FOR_ALL_OPERATIONS=true",
	"SET_TAGS_FOR_ALL_OPERATIONS=false",
	"DISABLE_ALL=true",
}
