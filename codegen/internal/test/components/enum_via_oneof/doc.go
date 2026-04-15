// Package enum_via_oneof tests the OpenAPI 3.1 enum-via-oneOf idiom:
// a scalar `type` with `oneOf` branches that each carry `const` + `title`
// is generated as a Go enum with named constants and per-value doc comments.
package enum_via_oneof

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
