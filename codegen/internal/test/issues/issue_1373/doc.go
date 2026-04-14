// Package issue_1373 tests recursive allOf self-references in schema definitions.
// https://github.com/oapi-codegen/oapi-codegen/issues/1373
package issue_1373

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
