// Package nested_aggregate tests complex nesting of allOf, anyOf, and oneOf
// including arrays of anyOf, objects with anyOf/oneOf properties, and allOf
// containing oneOf variants.
package nested_aggregate

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
