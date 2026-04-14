package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStyleSimpleParam_Primitive(t *testing.T) {
	result, err := StyleSimpleParam("id", 5, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, "5", result)
}

func TestStyleSimpleParam_String(t *testing.T) {
	result, err := StyleSimpleParam("name", "hello", ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, "hello", result)
}

func TestStyleSimpleParam_StringSlice(t *testing.T) {
	result, err := StyleSimpleParam("tags", []string{"a", "b", "c"}, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, "a,b,c", result)
}

func TestStyleSimpleParam_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
	}
	result, err := StyleSimpleParam("filter", obj{Color: "red", Size: 10}, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, "color,red,size,10", result)
}

func TestStyleSimpleParam_Roundtrip_Primitive(t *testing.T) {
	styled, err := StyleSimpleParam("id", 42, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result int
	err = BindSimpleParam("id", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestStyleSimpleParam_Roundtrip_StringSlice(t *testing.T) {
	original := []string{"x", "y", "z"}
	styled, err := StyleSimpleParam("items", original, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result []string
	err = BindSimpleParam("items", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleSimpleParam_Roundtrip_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	}
	original := obj{Color: "blue", Size: "large"}
	styled, err := StyleSimpleParam("filter", original, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result obj
	err = BindSimpleParam("filter", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleSimpleExplodeParam_Primitive(t *testing.T) {
	result, err := StyleSimpleParam("id", 5, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "5", result)
}

func TestStyleSimpleExplodeParam_StringSlice(t *testing.T) {
	result, err := StyleSimpleParam("tags", []string{"a", "b", "c"}, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "a,b,c", result)
}

func TestStyleSimpleExplodeParam_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
	}
	result, err := StyleSimpleParam("filter", obj{Color: "red", Size: 10}, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, "color=red,size=10", result)
}

func TestStyleSimpleExplodeParam_Roundtrip_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	}
	original := obj{Color: "blue", Size: "large"}
	styled, err := StyleSimpleParam("filter", original, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)

	var result obj
	err = BindSimpleParam("filter", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleSimpleExplodeParam_Roundtrip_StringSlice(t *testing.T) {
	original := []string{"a", "b", "c"}
	styled, err := StyleSimpleParam("items", original, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)

	var result []string
	err = BindSimpleParam("items", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}
