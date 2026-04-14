// Package issue_1530 tests oneOf with discriminator having multiple mappings
// pointing to the same schema.
// https://github.com/oapi-codegen/oapi-codegen/issues/1530
package issue_1530

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
