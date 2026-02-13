// Package codegen provides the public API for oapi-codegen's experimental code generator.
//
// This package re-exports the core types and functions from the internal
// implementation, providing a stable public interface for external consumers.
package codegen

import (
	"github.com/pb33f/libopenapi"

	impl "github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal"
)

// Configuration is the top-level configuration for code generation.
type Configuration = impl.Configuration

// GenerationOptions controls which parts of the code are generated.
type GenerationOptions = impl.GenerationOptions

// OutputOptions controls filtering of operations and schemas.
type OutputOptions = impl.OutputOptions

// ModelsPackage specifies an external package containing the model types.
type ModelsPackage = impl.ModelsPackage

// TypeMapping allows customizing OpenAPI type/format to Go type mappings.
type TypeMapping = impl.TypeMapping

// NameMangling configures how OpenAPI names are converted to Go identifiers.
type NameMangling = impl.NameMangling

// NameSubstitutions allows direct overrides of generated names.
type NameSubstitutions = impl.NameSubstitutions

// StructTagsConfig configures how struct tags are generated for fields.
type StructTagsConfig = impl.StructTagsConfig

// RuntimePackageConfig specifies an external package containing runtime helpers.
type RuntimePackageConfig = impl.RuntimePackageConfig

// RuntimeOutput holds the generated code for each runtime sub-package.
type RuntimeOutput = impl.RuntimeOutput

// Generate produces Go code from the parsed OpenAPI document.
// specData is the raw spec bytes used to embed the spec in the generated code.
func Generate(doc libopenapi.Document, specData []byte, cfg Configuration) (string, error) {
	return impl.Generate(doc, specData, cfg)
}

// GenerateRuntime produces standalone Go source files for each of the three
// runtime sub-packages (types, params, helpers). baseImportPath is the base
// import path for the runtime module (e.g., "github.com/org/project/runtime").
func GenerateRuntime(baseImportPath string) (*RuntimeOutput, error) {
	return impl.GenerateRuntime(baseImportPath)
}
