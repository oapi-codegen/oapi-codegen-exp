// Package runtime provides shared helper types and functions for code generated
// by oapi-codegen. The contents of this package are produced by the code
// generator itself from the same embeddable templates that are normally inlined
// into every generated file. By generating the helpers once into standalone
// sub-packages (types, params, helpers), multiple generated packages can import
// them instead of each duplicating their own copy.
//
// Sub-packages:
//   - types/   — custom Go types for OpenAPI format mappings (Date, Email, UUID, File, Nullable)
//   - params/  — parameter serialization/deserialization functions
//   - helpers/ — utility functions for request body encoding (MarshalForm)
//
// Regenerate after changing any codegen template:
//
//	go generate ./runtime/...
package runtime

//go:generate go run ../cmd/oapi-codegen --generate-runtime github.com/oapi-codegen/oapi-codegen-exp/experimental/runtime
