// Package issue_1189 tests anyOf, allOf, oneOf with duplicate enum variants.
// https://github.com/oapi-codegen/oapi-codegen/issues/1189
package issue_1189

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
