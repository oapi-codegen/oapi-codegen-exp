// Package recursive_all_of tests recursive allOf self-references in schema
// definitions without causing a stack overflow.
package recursive_all_of

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
