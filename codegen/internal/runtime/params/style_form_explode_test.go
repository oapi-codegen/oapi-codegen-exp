package params

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStyleFormExplodeParam_Primitive(t *testing.T) {
	result, err := StyleFormParam("color", "blue", ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "color=blue", result)
}

func TestStyleFormExplodeParam_Int(t *testing.T) {
	result, err := StyleFormParam("count", 5, ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "count=5", result)
}

func TestStyleFormExplodeParam_StringSlice(t *testing.T) {
	result, err := StyleFormParam("tags", []string{"a", "b", "c"}, ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "tags=a&tags=b&tags=c", result)
}

func TestStyleFormExplodeParam_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
	}
	result, err := StyleFormParam("filter", obj{Color: "red", Size: 10}, ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "color=red&size=10", result)
}

func TestStyleFormExplodeParam_Roundtrip_Primitive(t *testing.T) {
	styled, err := StyleFormParam("color", "blue", ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)

	// Parse to url.Values
	vals, err := url.ParseQuery(styled)
	require.NoError(t, err)

	var result string
	err = BindFormQueryParam("color", vals, &result, ParameterOptions{Explode: true, Required: true})
	require.NoError(t, err)
	assert.Equal(t, "blue", result)
}

func TestStyleFormExplodeParam_Roundtrip_StringSlice(t *testing.T) {
	original := []string{"a", "b", "c"}
	styled, err := StyleFormParam("items", original, ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)

	// Parse to url.Values
	vals, err := url.ParseQuery(styled)
	require.NoError(t, err)

	var result []string
	err = BindFormQueryParam("items", vals, &result, ParameterOptions{Explode: true, Required: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleFormExplodeParam_Roundtrip_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	}
	original := obj{Color: "blue", Size: "large"}
	styled, err := StyleFormParam("filter", original, ParameterOptions{ParamLocation: ParamLocationQuery, Explode: true})
	require.NoError(t, err)

	// Parse to url.Values
	vals, err := url.ParseQuery(styled)
	require.NoError(t, err)

	var result obj
	err = BindFormQueryParam("filter", vals, &result, ParameterOptions{Explode: true, Required: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestBindFormExplodeParam_OptionalMissing(t *testing.T) {
	vals := url.Values{}

	var result *string
	err := BindFormQueryParam("missing", vals, &result, ParameterOptions{Explode: true})
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBindFormExplodeParam_RequiredMissing(t *testing.T) {
	vals := url.Values{}

	var result string
	err := BindFormQueryParam("required", vals, &result, ParameterOptions{Explode: true, Required: true})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}
