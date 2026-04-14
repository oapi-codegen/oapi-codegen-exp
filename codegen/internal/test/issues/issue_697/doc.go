// Package issue_697 tests properties alongside allOf should not be ignored.
// https://github.com/oapi-codegen/oapi-codegen/issues/697
package issue_697

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
