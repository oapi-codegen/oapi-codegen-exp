package codegen

import (
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// targetSchema parses a minimal OpenAPI 3.1 spec with a single component
// schema named "Target" and returns the resolved high-level schema.
// The input is the YAML body of the Target schema (indented by 6 spaces to
// sit under `components.schemas.Target`).
func targetSchema(t *testing.T, targetYAML string) *base.Schema {
	t.Helper()
	const preamble = "openapi: 3.1.0\n" +
		"info:\n" +
		"  title: t\n" +
		"  version: \"1\"\n" +
		"paths: {}\n" +
		"components:\n" +
		"  schemas:\n" +
		"    Target:\n"
	doc, err := libopenapi.NewDocument([]byte(preamble + targetYAML))
	require.NoError(t, err)
	model, errs := doc.BuildV3Model()
	require.Empty(t, errs, "BuildV3Model errors")
	require.NotNil(t, model)
	proxy := model.Model.Components.Schemas.GetOrZero("Target")
	require.NotNil(t, proxy)
	sch := proxy.Schema()
	require.NotNil(t, sch)
	return sch
}

func TestIsConstOneOfEnum_Integer(t *testing.T) {
	sch := targetSchema(t, `      type: integer
      oneOf:
        - title: HIGH
          const: 2
          description: An urgent problem
        - title: MEDIUM
          const: 1
        - title: LOW
          const: 0
          description: Can wait forever
`)

	items, ok := isConstOneOfEnum(sch)
	require.True(t, ok)
	require.Len(t, items, 3)
	assert.Equal(t, "HIGH", items[0].Title)
	assert.Equal(t, "2", items[0].Value)
	assert.Equal(t, "An urgent problem", items[0].Doc)
	assert.Equal(t, "MEDIUM", items[1].Title)
	assert.Equal(t, "1", items[1].Value)
	assert.Equal(t, "", items[1].Doc)
	assert.Equal(t, "LOW", items[2].Title)
	assert.Equal(t, "0", items[2].Value)
	assert.Equal(t, "Can wait forever", items[2].Doc)
}

func TestIsConstOneOfEnum_String(t *testing.T) {
	sch := targetSchema(t, `      type: string
      oneOf:
        - title: Red
          const: r
        - title: Green
          const: g
        - title: Blue
          const: b
`)

	items, ok := isConstOneOfEnum(sch)
	require.True(t, ok)
	require.Len(t, items, 3)
	assert.Equal(t, "Red", items[0].Title)
	assert.Equal(t, "r", items[0].Value)
}

func TestIsConstOneOfEnum_MissingTitle(t *testing.T) {
	sch := targetSchema(t, `      type: integer
      oneOf:
        - title: HIGH
          const: 2
        - const: 1
`)

	_, ok := isConstOneOfEnum(sch)
	assert.False(t, ok, "missing title on one branch must disqualify the idiom")
}

func TestIsConstOneOfEnum_MissingConst(t *testing.T) {
	sch := targetSchema(t, `      type: integer
      oneOf:
        - title: HIGH
          const: 2
        - title: MEDIUM
`)

	_, ok := isConstOneOfEnum(sch)
	assert.False(t, ok, "missing const on one branch must disqualify the idiom")
}

func TestIsConstOneOfEnum_NonScalarOuterType(t *testing.T) {
	sch := targetSchema(t, `      type: object
      oneOf:
        - title: HIGH
          const: 2
        - title: LOW
          const: 0
`)

	_, ok := isConstOneOfEnum(sch)
	assert.False(t, ok, "object outer type must disqualify the idiom")
}

func TestIsConstOneOfEnum_EmptyOneOf(t *testing.T) {
	sch := targetSchema(t, `      type: integer
`)

	_, ok := isConstOneOfEnum(sch)
	assert.False(t, ok, "no oneOf means no idiom")
}

func TestIsConstOneOfEnum_NestedComposition(t *testing.T) {
	sch := targetSchema(t, `      type: integer
      oneOf:
        - title: HIGH
          const: 2
          oneOf:
            - const: 3
            - const: 4
        - title: LOW
          const: 0
`)

	_, ok := isConstOneOfEnum(sch)
	assert.False(t, ok, "a branch with nested composition must disqualify the idiom")
}

func TestIsConstOneOfEnum_NilSchema(t *testing.T) {
	_, ok := isConstOneOfEnum(nil)
	assert.False(t, ok)
}
