package codegen

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
)

// constOneOfItem is one branch of an OpenAPI 3.1 enum-via-oneOf schema.
// It captures the per-value name (from `title`), the raw enum value
// (stringified from `const`), and the doc comment (from `description`).
type constOneOfItem struct {
	Title string
	Value string
	Doc   string
}

// isConstOneOfEnum reports whether a schema matches the OpenAPI 3.1
// enum-via-oneOf idiom:
//
//	type: integer|string
//	oneOf:
//	  - { title: NAME, const: VALUE, description?: TEXT }
//	  - ...
//
// All members must carry both `title` and `const`, and no member may itself
// be a composition (oneOf/allOf/anyOf) or declare properties. The outer
// schema's primary type must be a scalar (string or integer).
//
// When the idiom matches, the per-branch values are returned in declaration
// order. Otherwise returns (nil, false).
func isConstOneOfEnum(schema *base.Schema) ([]constOneOfItem, bool) {
	if schema == nil || len(schema.OneOf) == 0 {
		return nil, false
	}

	primary := getPrimaryType(schema)
	if primary != "string" && primary != "integer" {
		return nil, false
	}

	items := make([]constOneOfItem, 0, len(schema.OneOf))
	for _, proxy := range schema.OneOf {
		if proxy == nil {
			return nil, false
		}
		m := proxy.Schema()
		if m == nil {
			return nil, false
		}
		if m.Title == "" || m.Const == nil {
			return nil, false
		}
		// Members must be simple scalar-const schemas, not nested composition.
		if len(m.OneOf) > 0 || len(m.AllOf) > 0 || len(m.AnyOf) > 0 {
			return nil, false
		}
		if m.Properties != nil && m.Properties.Len() > 0 {
			return nil, false
		}
		items = append(items, constOneOfItem{
			Title: m.Title,
			Value: m.Const.Value,
			Doc:   m.Description,
		})
	}
	return items, true
}
