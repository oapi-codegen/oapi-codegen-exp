// Package same_level tests that properties defined at the same level as allOf
// are included in the generated type, not ignored.
package same_level

//go:generate go run ../../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
