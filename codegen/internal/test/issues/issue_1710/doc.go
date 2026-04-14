// Package issue_1710 tests nested allOf containing oneOf with field preservation.
// https://github.com/oapi-codegen/oapi-codegen/issues/1710
package issue_1710

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
