// Package issue_775 tests allOf used to add format (uuid, date) to base type.
// https://github.com/oapi-codegen/oapi-codegen/issues/775
package issue_775

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
