package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

// TestEnumGenerated verifies that the enum type is generated for properties inside anyOf.
// Issue 1429: enum type was not being generated when used inside anyOf.
func TestEnumGenerated(t *testing.T) {
	assert.Equal(t, TestAnyOf1FieldA("foo"), Foo)
	assert.Equal(t, TestAnyOf1FieldA("bar"), Bar)
}

// TestFromAsTestAnyOf0 tests the variant with a plain string field.
func TestFromAsTestAnyOf0(t *testing.T) {
	var u Test
	err := u.FromTestAnyOf0(TestAnyOf0{FieldA: ptr("hello")})
	require.NoError(t, err)

	data, err := u.MarshalJSON()
	require.NoError(t, err)

	var decoded Test
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	got, err := decoded.AsTestAnyOf0()
	require.NoError(t, err)
	require.NotNil(t, got.FieldA)
	assert.Equal(t, "hello", *got.FieldA)
}

// TestFromAsTestAnyOf1 tests the variant with an enum-constrained string field.
func TestFromAsTestAnyOf1(t *testing.T) {
	var u Test
	err := u.FromTestAnyOf1(TestAnyOf1{FieldA: ptr(string(Foo))})
	require.NoError(t, err)

	data, err := u.MarshalJSON()
	require.NoError(t, err)

	var decoded Test
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	got, err := decoded.AsTestAnyOf1()
	require.NoError(t, err)
	require.NotNil(t, got.FieldA)
	assert.Equal(t, "foo", *got.FieldA)
}

// TestUnmarshalJSONIntoUnion verifies that raw JSON unmarshals into the union
// and can be read back via both As* accessors (anyOf semantics).
func TestUnmarshalJSONIntoUnion(t *testing.T) {
	input := `{"fieldA":"bar"}`

	var u Test
	err := u.UnmarshalJSON([]byte(input))
	require.NoError(t, err)

	v0, err := u.AsTestAnyOf0()
	require.NoError(t, err)
	require.NotNil(t, v0.FieldA)
	assert.Equal(t, "bar", *v0.FieldA)

	v1, err := u.AsTestAnyOf1()
	require.NoError(t, err)
	require.NotNil(t, v1.FieldA)
	assert.Equal(t, "bar", *v1.FieldA)
}

func TestApplyDefaults(t *testing.T) {
	(&Test{}).ApplyDefaults()
	(&TestAnyOf0{}).ApplyDefaults()
	(&TestAnyOf1{}).ApplyDefaults()
}
