package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oapiCodegenTypesPkg "github.com/oapi-codegen/oapi-codegen-exp/runtime/types"
)

func ptrTo[T any](v T) *T {
	return &v
}

// TestAnyOfWithSingleRef verifies that anyOf with a single $ref generates
// correct types that can be used.
// https://github.com/oapi-codegen/oapi-codegen/issues/502
func TestAnyOfWithSingleRef(t *testing.T) {
	claims := OptionalClaims{
		IDToken:     ptrTo("id-token-value"),
		AccessToken: ptrTo("access-token-value"),
	}

	assert.Equal(t, "id-token-value", *claims.IDToken)
	assert.Equal(t, "access-token-value", *claims.AccessToken)
}

func TestApplicationWithNonNullOptionalClaims(t *testing.T) {
	// Build the union variant via FromOptionalClaims.
	var unionVal ApplicationOptionalClaims
	require.NoError(t, unionVal.FromOptionalClaims(OptionalClaims{
		IDToken: ptrTo("token"),
	}))

	app := Application{
		Name:           ptrTo("my-app"),
		OptionalClaims: oapiCodegenTypesPkg.NewNullableWithValue(unionVal),
	}

	assert.Equal(t, "my-app", *app.Name)
	require.True(t, app.OptionalClaims.IsSpecified(), "OptionalClaims should be specified")
	require.False(t, app.OptionalClaims.IsNull(), "OptionalClaims should not be null")

	got := app.OptionalClaims.MustGet()
	extracted, err := got.AsOptionalClaims()
	require.NoError(t, err)
	require.NotNil(t, extracted.IDToken)
	assert.Equal(t, "token", *extracted.IDToken)
}

func TestApplicationJSONRoundTrip(t *testing.T) {
	var unionVal ApplicationOptionalClaims
	require.NoError(t, unionVal.FromOptionalClaims(OptionalClaims{
		IDToken:     ptrTo("id"),
		AccessToken: ptrTo("access"),
	}))

	original := Application{
		Name:           ptrTo("test-app"),
		OptionalClaims: oapiCodegenTypesPkg.NewNullableWithValue(unionVal),
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded Application
	require.NoError(t, json.Unmarshal(data, &decoded))

	assert.Equal(t, *original.Name, *decoded.Name)
	require.True(t, decoded.OptionalClaims.IsSpecified(), "OptionalClaims should be specified after round trip")
	require.False(t, decoded.OptionalClaims.IsNull(), "OptionalClaims should not be null after round trip")

	got := decoded.OptionalClaims.MustGet()
	extracted, err := got.AsOptionalClaims()
	require.NoError(t, err)
	require.NotNil(t, extracted.IDToken)
	assert.Equal(t, "id", *extracted.IDToken)
	require.NotNil(t, extracted.AccessToken)
	assert.Equal(t, "access", *extracted.AccessToken)
}

func TestApplicationNullOptionalClaims(t *testing.T) {
	app := Application{
		Name:           ptrTo("null-test-app"),
		OptionalClaims: oapiCodegenTypesPkg.NewNullNullable[ApplicationOptionalClaims](),
	}

	require.True(t, app.OptionalClaims.IsSpecified(), "OptionalClaims should be specified")
	require.True(t, app.OptionalClaims.IsNull(), "OptionalClaims should be null")

	data, err := json.Marshal(app)
	require.NoError(t, err)

	// The JSON should contain "optionalClaims":null.
	var raw map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(data, &raw))
	assert.JSONEq(t, "null", string(raw["optionalClaims"]))

	// Round-trip: unmarshalling null JSON should produce a null Nullable.
	var decoded Application
	require.NoError(t, json.Unmarshal(data, &decoded))
	require.True(t, decoded.OptionalClaims.IsSpecified())
	assert.True(t, decoded.OptionalClaims.IsNull())
}

func TestApplicationOptionalClaimsUnionFromAsRoundTrip(t *testing.T) {
	original := OptionalClaims{
		IDToken:     ptrTo("my-id"),
		AccessToken: ptrTo("my-access"),
	}

	var union ApplicationOptionalClaims
	require.NoError(t, union.FromOptionalClaims(original))

	// Marshal and unmarshal the union to simulate a JSON round-trip.
	data, err := union.MarshalJSON()
	require.NoError(t, err)

	var decoded ApplicationOptionalClaims
	require.NoError(t, decoded.UnmarshalJSON(data))

	extracted, err := decoded.AsOptionalClaims()
	require.NoError(t, err)
	require.NotNil(t, extracted.IDToken)
	assert.Equal(t, "my-id", *extracted.IDToken)
	require.NotNil(t, extracted.AccessToken)
	assert.Equal(t, "my-access", *extracted.AccessToken)
}

func TestApplyDefaults(t *testing.T) {
	// ApplyDefaults should be callable without panic on all types,
	// even when fields are zero-valued.
	t.Run("OptionalClaims", func(t *testing.T) {
		var c OptionalClaims
		assert.NotPanics(t, func() { c.ApplyDefaults() })
		assert.Nil(t, c.IDToken)
		assert.Nil(t, c.AccessToken)
	})

	t.Run("Application", func(t *testing.T) {
		var a Application
		assert.NotPanics(t, func() { a.ApplyDefaults() })
		assert.Nil(t, a.Name)
		assert.False(t, a.OptionalClaims.IsSpecified())
	})

	t.Run("ApplicationOptionalClaims", func(t *testing.T) {
		var u ApplicationOptionalClaims
		assert.NotPanics(t, func() { u.ApplyDefaults() })
	})
}
