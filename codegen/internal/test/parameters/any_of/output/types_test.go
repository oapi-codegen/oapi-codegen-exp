package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFromTestAnyOf0AsTestAnyOf0RoundTrip verifies that setting the first
// anyOf variant via FromTestAnyOf0 and reading it back via AsTestAnyOf0
// preserves the data.
func TestFromTestAnyOf0AsTestAnyOf0RoundTrip(t *testing.T) {
	var u Test
	err := u.FromTestAnyOf0(TestAnyOf0{
		Item1: "value1",
		Item2: "value2",
	})
	require.NoError(t, err)

	got, err := u.AsTestAnyOf0()
	require.NoError(t, err)
	assert.Equal(t, "value1", got.Item1)
	assert.Equal(t, "value2", got.Item2)
}

// TestFromTestAnyOf1AsTestAnyOf1RoundTrip verifies the second anyOf variant
// round-trips through From/As.
func TestFromTestAnyOf1AsTestAnyOf1RoundTrip(t *testing.T) {
	item2 := "hello"
	item3 := "world"

	var u Test
	err := u.FromTestAnyOf1(TestAnyOf1{
		Item2: &item2,
		Item3: &item3,
	})
	require.NoError(t, err)

	got, err := u.AsTestAnyOf1()
	require.NoError(t, err)
	require.NotNil(t, got.Item2)
	require.NotNil(t, got.Item3)
	assert.Equal(t, "hello", *got.Item2)
	assert.Equal(t, "world", *got.Item3)
}

// TestTestMarshalBothVariantsMerged verifies that for an anyOf union, merging
// both variants produces JSON containing fields from both.
func TestTestMarshalBothVariantsMerged(t *testing.T) {
	item2 := "shared"
	item3 := "extra"

	var u Test
	err := u.FromTestAnyOf0(TestAnyOf0{Item1: "first", Item2: "second"})
	require.NoError(t, err)
	err = u.MergeTestAnyOf1(TestAnyOf1{Item2: &item2, Item3: &item3})
	require.NoError(t, err)

	data, err := json.Marshal(u)
	require.NoError(t, err)

	var m map[string]any
	err = json.Unmarshal(data, &m)
	require.NoError(t, err)

	assert.Equal(t, "first", m["item1"])
	assert.Equal(t, "extra", m["item3"])
	// item2 should be present (the merge overwrites with "shared")
	assert.Equal(t, "shared", m["item2"])
}

// TestTestUnmarshalBothVariants verifies that unmarshaling JSON with fields
// from both anyOf variants allows reading either variant back.
func TestTestUnmarshalBothVariants(t *testing.T) {
	input := `{"item1":"a","item2":"b","item3":"c"}`

	var u Test
	err := json.Unmarshal([]byte(input), &u)
	require.NoError(t, err)

	// AsTestAnyOf0: item1 and item2 should be populated
	v0, err := u.AsTestAnyOf0()
	require.NoError(t, err)
	assert.Equal(t, "a", v0.Item1)
	assert.Equal(t, "b", v0.Item2)

	// AsTestAnyOf1: item2 and item3 should be populated
	v1, err := u.AsTestAnyOf1()
	require.NoError(t, err)
	require.NotNil(t, v1.Item2)
	require.NotNil(t, v1.Item3)
	assert.Equal(t, "b", *v1.Item2)
	assert.Equal(t, "c", *v1.Item3)
}

// TestTest2FromInt0AsInt0RoundTrip verifies int round-trip through the
// oneOf union.
func TestTest2FromInt0AsInt0RoundTrip(t *testing.T) {
	var u Test2
	err := u.FromInt0(42)
	require.NoError(t, err)

	got, err := u.AsInt0()
	require.NoError(t, err)
	assert.Equal(t, 42, got)
}

// TestTest2FromString1AsString1RoundTrip verifies string round-trip through
// the oneOf union.
func TestTest2FromString1AsString1RoundTrip(t *testing.T) {
	var u Test2
	err := u.FromString1("hello")
	require.NoError(t, err)

	got, err := u.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "hello", got)
}

// TestTest2MarshalInt verifies that a Test2 holding an int marshals to a bare
// JSON number.
func TestTest2MarshalInt(t *testing.T) {
	var u Test2
	err := u.FromInt0(99)
	require.NoError(t, err)

	data, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Equal(t, "99", string(data))
}

// TestTest2MarshalString verifies that a Test2 holding a string marshals to a
// JSON string.
func TestTest2MarshalString(t *testing.T) {
	var u Test2
	err := u.FromString1("world")
	require.NoError(t, err)

	data, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Equal(t, `"world"`, string(data))
}

// TestTest2MarshalBothSetLastWins verifies the behavior when both variants are
// set sequentially -- the last From* call wins since each overwrites the union.
func TestTest2MarshalBothSetLastWins(t *testing.T) {
	var u Test2
	err := u.FromInt0(1)
	require.NoError(t, err)
	err = u.FromString1("one")
	require.NoError(t, err)

	data, err := json.Marshal(u)
	require.NoError(t, err)
	// The last From call (FromString1) should win
	assert.Equal(t, `"one"`, string(data))
}

// TestTest2UnmarshalInt verifies that unmarshaling a JSON integer populates
// the int variant.
func TestTest2UnmarshalInt(t *testing.T) {
	var u Test2
	err := json.Unmarshal([]byte(`42`), &u)
	require.NoError(t, err)

	got, err := u.AsInt0()
	require.NoError(t, err)
	assert.Equal(t, 42, got)
}

// TestTest2UnmarshalString verifies that unmarshaling a JSON string populates
// the string variant.
func TestTest2UnmarshalString(t *testing.T) {
	var u Test2
	err := json.Unmarshal([]byte(`"world"`), &u)
	require.NoError(t, err)

	got, err := u.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "world", got)
}

// TestGetTestParameterAlias verifies the GetTestParameter type alias is a
// slice of Test2 and works as expected.
func TestGetTestParameterAlias(t *testing.T) {
	var params GetTestParameter

	var elem Test2
	err := elem.FromInt0(10)
	require.NoError(t, err)

	params = append(params, elem)
	require.Len(t, params, 1)

	got, err := params[0].AsInt0()
	require.NoError(t, err)
	assert.Equal(t, 10, got)
}

// TestApplyDefaults verifies that ApplyDefaults can be called on all types
// without panic.
func TestApplyDefaults(t *testing.T) {
	u := &Test{}
	u.ApplyDefaults()

	v0 := &TestAnyOf0{}
	v0.ApplyDefaults()

	v1 := &TestAnyOf1{}
	v1.ApplyDefaults()

	u2 := &Test2{}
	u2.ApplyDefaults()
}
