package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSeverityConstants verifies that the integer enum-via-oneOf idiom
// produces a `type Severity int` with titled constants and the right values.
func TestSeverityConstants(t *testing.T) {
	assert.Equal(t, 2, int(HIGH))
	assert.Equal(t, 1, int(MEDIUM))
	assert.Equal(t, 0, int(LOW))
}

// TestSeverityJSONRoundTrip verifies Severity marshals as its integer value.
func TestSeverityJSONRoundTrip(t *testing.T) {
	data, err := json.Marshal(HIGH)
	require.NoError(t, err)
	assert.JSONEq(t, `2`, string(data))

	var got Severity
	require.NoError(t, json.Unmarshal([]byte(`1`), &got))
	assert.Equal(t, MEDIUM, got)
}

// TestColorConstants verifies that the string enum-via-oneOf idiom produces
// a `type Color string` with titled constants and the right string values.
func TestColorConstants(t *testing.T) {
	assert.Equal(t, "r", string(Red))
	assert.Equal(t, "g", string(Green))
	assert.Equal(t, "b", string(Blue))
}

// TestColorJSONRoundTrip verifies Color marshals as its string value.
func TestColorJSONRoundTrip(t *testing.T) {
	data, err := json.Marshal(Red)
	require.NoError(t, err)
	assert.JSONEq(t, `"r"`, string(data))

	var got Color
	require.NoError(t, json.Unmarshal([]byte(`"b"`), &got))
	assert.Equal(t, Blue, got)
}

// TestMixedOneOfStillUnion verifies that a oneOf whose branches do NOT all
// have `title` + `const` (here, the second branch lacks `title`) falls
// through to the standard union generator rather than becoming an enum.
// If MixedOneOf were mis-detected as an enum, this file would fail to
// compile (the union methods would not exist).
func TestMixedOneOfStillUnion(t *testing.T) {
	var m MixedOneOf
	require.NoError(t, m.FromAny0("a"))

	data, err := json.Marshal(m)
	require.NoError(t, err)
	assert.JSONEq(t, `"a"`, string(data))
}
