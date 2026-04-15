package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldA_StringRoundTrip(t *testing.T) {
	var fa TestFieldA
	require.NoError(t, fa.FromString0("plain-string"))

	data, err := fa.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `"plain-string"`, string(data))

	var fa2 TestFieldA
	require.NoError(t, fa2.UnmarshalJSON(data))

	got, err := fa2.AsString0()
	require.NoError(t, err)
	assert.Equal(t, "plain-string", got)
}

func TestFieldA_EnumRoundTrip(t *testing.T) {
	var fa TestFieldA
	require.NoError(t, fa.FromTestFieldAAnyOf1(TestFieldAAnyOf1Foo))

	data, err := fa.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `"foo"`, string(data))

	var fa2 TestFieldA
	require.NoError(t, fa2.UnmarshalJSON(data))

	got, err := fa2.AsTestFieldAAnyOf1()
	require.NoError(t, err)
	assert.Equal(t, TestFieldAAnyOf1Foo, got)
}

func TestFieldC_StringRoundTrip(t *testing.T) {
	var fc TestFieldC
	require.NoError(t, fc.FromString0("one-of-string"))

	data, err := fc.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `"one-of-string"`, string(data))

	var fc2 TestFieldC
	require.NoError(t, fc2.UnmarshalJSON(data))

	got, err := fc2.AsString0()
	require.NoError(t, err)
	assert.Equal(t, "one-of-string", got)
}

func TestFieldC_EnumRoundTrip(t *testing.T) {
	var fc TestFieldC
	require.NoError(t, fc.FromTestFieldCOneOf1(TestFieldCOneOf1Bar))

	data, err := fc.MarshalJSON()
	require.NoError(t, err)
	assert.JSONEq(t, `"bar"`, string(data))

	var fc2 TestFieldC
	require.NoError(t, fc2.UnmarshalJSON(data))

	got, err := fc2.AsTestFieldCOneOf1()
	require.NoError(t, err)
	assert.Equal(t, TestFieldCOneOf1Bar, got)
}

func TestEnumConstants(t *testing.T) {
	assert.Equal(t, TestFieldAAnyOf1("foo"), TestFieldAAnyOf1Foo)
	assert.Equal(t, TestFieldAAnyOf1("bar"), TestFieldAAnyOf1Bar)
	assert.Equal(t, TestFieldCOneOf1("foo"), TestFieldCOneOf1Foo)
	assert.Equal(t, TestFieldCOneOf1("bar"), TestFieldCOneOf1Bar)
}

func TestFullStructJSONRoundTrip(t *testing.T) {
	var fieldA TestFieldA
	require.NoError(t, fieldA.FromString0("hello"))

	var fieldC TestFieldC
	require.NoError(t, fieldC.FromTestFieldCOneOf1(TestFieldCOneOf1Foo))

	original := Test{
		FieldA: &fieldA,
		FieldB: &TestFieldB{},
		FieldC: &fieldC,
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded Test
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.NotNil(t, decoded.FieldA)
	gotA, err := decoded.FieldA.AsString0()
	require.NoError(t, err)
	assert.Equal(t, "hello", gotA)

	require.NotNil(t, decoded.FieldB)

	require.NotNil(t, decoded.FieldC)
	gotC, err := decoded.FieldC.AsTestFieldCOneOf1()
	require.NoError(t, err)
	assert.Equal(t, TestFieldCOneOf1Foo, gotC)
}

func TestApplyDefaults(t *testing.T) {
	test := &Test{}
	test.ApplyDefaults()

	fa := &TestFieldA{}
	fa.ApplyDefaults()

	fb := &TestFieldB{}
	fb.ApplyDefaults()

	fc := &TestFieldC{}
	fc.ApplyDefaults()
}

func TestGetOpenAPISpecJSON(t *testing.T) {
	data, err := GetOpenAPISpecJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}
