// Package nested_aggregate tests complex nesting of allOf, anyOf, and oneOf:
// arrays of anyOf, objects with anyOf/oneOf properties, allOf containing oneOf,
// oneOf with nested allOf and field preservation, and composition with enums.
package nested_aggregate

//go:generate go run ../../../../../cmd/oapi-codegen -config config.yaml spec.yaml
