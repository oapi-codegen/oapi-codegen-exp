package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStyleLabelParam_Primitive(t *testing.T) {
	result, err := StyleLabelParam("id", 5, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, ".5", result)
}

func TestStyleLabelParam_String(t *testing.T) {
	result, err := StyleLabelParam("color", "blue", ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, ".blue", result)
}

func TestStyleLabelParam_StringSlice(t *testing.T) {
	result, err := StyleLabelParam("tags", []string{"a", "b", "c"}, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, ".a,b,c", result)
}

func TestStyleLabelParam_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
	}
	result, err := StyleLabelParam("filter", obj{Color: "red", Size: 10}, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, ".color,red,size,10", result)
}

func TestStyleLabelParam_Roundtrip_Primitive(t *testing.T) {
	styled, err := StyleLabelParam("id", 42, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result int
	err = BindLabelParam("id", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestStyleLabelParam_Roundtrip_StringSlice(t *testing.T) {
	original := []string{"x", "y", "z"}
	styled, err := StyleLabelParam("items", original, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result []string
	err = BindLabelParam("items", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleLabelParam_Roundtrip_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	}
	original := obj{Color: "blue", Size: "large"}
	styled, err := StyleLabelParam("filter", original, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)

	var result obj
	err = BindLabelParam("filter", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleLabelExplodeParam_Primitive(t *testing.T) {
	result, err := StyleLabelParam("id", 5, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, ".5", result)
}

func TestStyleLabelExplodeParam_StringSlice(t *testing.T) {
	result, err := StyleLabelParam("tags", []string{"a", "b", "c"}, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, ".a.b.c", result)
}

func TestStyleLabelExplodeParam_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  int    `json:"size"`
	}
	result, err := StyleLabelParam("filter", obj{Color: "red", Size: 10}, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, ".color=red.size=10", result)
}

func TestStyleLabelExplodeParam_Roundtrip_Primitive(t *testing.T) {
	styled, err := StyleLabelParam("id", 42, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)

	var result int
	err = BindLabelParam("id", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestStyleLabelExplodeParam_Roundtrip_StringSlice(t *testing.T) {
	original := []string{"x", "y", "z"}
	styled, err := StyleLabelParam("items", original, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)

	var result []string
	err = BindLabelParam("items", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}

func TestStyleLabelExplodeParam_Roundtrip_Struct(t *testing.T) {
	type obj struct {
		Color string `json:"color"`
		Size  string `json:"size"`
	}
	original := obj{Color: "blue", Size: "large"}
	styled, err := StyleLabelParam("filter", original, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)

	var result obj
	err = BindLabelParam("filter", styled, &result, ParameterOptions{ParamLocation: ParamLocationPath, Explode: true})
	require.NoError(t, err)
	assert.Equal(t, original, result)
}
