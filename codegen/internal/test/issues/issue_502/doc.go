// Package issue_502 tests anyOf with single $ref should not generate interface{}.
// https://github.com/oapi-codegen/oapi-codegen/issues/502
package issue_502

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
