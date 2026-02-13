package codegen

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRuntime(t *testing.T) {
	rt, err := GenerateRuntime("github.com/example/runtime")
	require.NoError(t, err)
	require.NotNil(t, rt)

	t.Run("types", func(t *testing.T) {
		code := rt.Types
		require.NotEmpty(t, code)

		assert.Contains(t, code, "package types")
		assert.Contains(t, code, "type Date struct")
		assert.Contains(t, code, "type Nullable[")
		assert.Contains(t, code, "type Email ")
		assert.Contains(t, code, "type File struct")
		assert.True(t, strings.HasPrefix(code, "// Code generated"))
	})

	t.Run("params", func(t *testing.T) {
		code := rt.Params
		require.NotEmpty(t, code)

		assert.Contains(t, code, "package params")
		assert.Contains(t, code, "type ParamLocation int")
		assert.Contains(t, code, "func BindStringToObject(")
		assert.Contains(t, code, "func marshalKnownTypes(")
		assert.Contains(t, code, "func StyleSimpleParam(")
		assert.Contains(t, code, "func StyleFormParam(")
		assert.Contains(t, code, "func BindSimpleParam(")
		assert.Contains(t, code, "func BindFormParam(")

		// Params should reference types.Date, not bare Date
		assert.Contains(t, code, "types.Date{}")
		assert.Contains(t, code, "types.DateFormat")

		assert.True(t, strings.HasPrefix(code, "// Code generated"))
	})

	t.Run("helpers", func(t *testing.T) {
		code := rt.Helpers
		require.NotEmpty(t, code)

		assert.Contains(t, code, "package helpers")
		assert.Contains(t, code, "func MarshalForm(")
		assert.True(t, strings.HasPrefix(code, "// Code generated"))
	})
}

func TestGenerateRuntimeEmptyPath(t *testing.T) {
	_, err := GenerateRuntime("")
	assert.Error(t, err)
}
