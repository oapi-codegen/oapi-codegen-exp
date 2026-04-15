// Package any_of_single_ref tests anyOf with a single $ref — should generate
// a typed property, not interface{}.
package any_of_single_ref

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
