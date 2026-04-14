package params

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// BindParameter (path/header/cookie — string value)
// ---------------------------------------------------------------------------

func TestBindParameter_Simple_Roundtrip(t *testing.T) {
	pathOpts := func(explode bool) ParameterOptions {
		return ParameterOptions{Style: "simple", ParamLocation: ParamLocationPath, Explode: explode}
	}

	t.Run("primitive", func(t *testing.T) {
		styled, err := StyleParameter("id", 42, pathOpts(false))
		require.NoError(t, err)

		var result int
		err = BindParameter("id", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	})
	t.Run("string_slice", func(t *testing.T) {
		original := []string{"x", "y", "z"}
		styled, err := StyleParameter("items", original, pathOpts(false))
		require.NoError(t, err)

		var result []string
		err = BindParameter("items", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  string `json:"size"`
		}
		original := obj{Color: "blue", Size: "large"}
		styled, err := StyleParameter("filter", original, pathOpts(false))
		require.NoError(t, err)

		var result obj
		err = BindParameter("filter", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  string `json:"size"`
		}
		original := obj{Color: "blue", Size: "large"}
		styled, err := StyleParameter("filter", original, pathOpts(true))
		require.NoError(t, err)

		var result obj
		err = BindParameter("filter", styled, &result, pathOpts(true))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		original := []string{"a", "b", "c"}
		styled, err := StyleParameter("items", original, pathOpts(true))
		require.NoError(t, err)

		var result []string
		err = BindParameter("items", styled, &result, pathOpts(true))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
}

func TestBindParameter_Label_Roundtrip(t *testing.T) {
	pathOpts := func(explode bool) ParameterOptions {
		return ParameterOptions{Style: "label", ParamLocation: ParamLocationPath, Explode: explode}
	}

	t.Run("primitive", func(t *testing.T) {
		styled, err := StyleParameter("id", 42, pathOpts(false))
		require.NoError(t, err)

		var result int
		err = BindParameter("id", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	})
	t.Run("string_slice", func(t *testing.T) {
		original := []string{"x", "y", "z"}
		styled, err := StyleParameter("items", original, pathOpts(false))
		require.NoError(t, err)

		var result []string
		err = BindParameter("items", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  string `json:"size"`
		}
		original := obj{Color: "blue", Size: "large"}
		styled, err := StyleParameter("filter", original, pathOpts(false))
		require.NoError(t, err)

		var result obj
		err = BindParameter("filter", styled, &result, pathOpts(false))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("explode_primitive", func(t *testing.T) {
		styled, err := StyleParameter("id", 42, pathOpts(true))
		require.NoError(t, err)

		var result int
		err = BindParameter("id", styled, &result, pathOpts(true))
		require.NoError(t, err)
		assert.Equal(t, 42, result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		original := []string{"x", "y", "z"}
		styled, err := StyleParameter("items", original, pathOpts(true))
		require.NoError(t, err)

		var result []string
		err = BindParameter("items", styled, &result, pathOpts(true))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  string `json:"size"`
		}
		original := obj{Color: "blue", Size: "large"}
		styled, err := StyleParameter("filter", original, pathOpts(true))
		require.NoError(t, err)

		var result obj
		err = BindParameter("filter", styled, &result, pathOpts(true))
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
}

func TestBindParameter_RequiredEmpty(t *testing.T) {
	var result string
	err := BindParameter("name", "", &result, ParameterOptions{Style: "simple", Required: true})
	require.Error(t, err)
	var reqErr *MissingRequiredParameterError
	assert.True(t, errors.As(err, &reqErr))
	assert.Equal(t, "name", reqErr.ParamName)
}

func TestBindParameter_OptionalEmpty(t *testing.T) {
	var result string
	err := BindParameter("name", "", &result, ParameterOptions{Style: "simple"})
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

// ---------------------------------------------------------------------------
// BindQueryParameter (query — url.Values)
// ---------------------------------------------------------------------------

func TestBindQueryParameter_Form_Roundtrip(t *testing.T) {
	t.Run("explode_primitive", func(t *testing.T) {
		styled, err := StyleParameter("color", "blue", ParameterOptions{Style: "form", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result string
		err = BindQueryParameter("color", vals, &result, ParameterOptions{Style: "form", Explode: true, Required: true})
		require.NoError(t, err)
		assert.Equal(t, "blue", result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		original := []string{"a", "b", "c"}
		styled, err := StyleParameter("items", original, ParameterOptions{Style: "form", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result []string
		err = BindQueryParameter("items", vals, &result, ParameterOptions{Style: "form", Explode: true, Required: true})
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  string `json:"size"`
		}
		original := obj{Color: "blue", Size: "large"}
		styled, err := StyleParameter("filter", original, ParameterOptions{Style: "form", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result obj
		err = BindQueryParameter("filter", vals, &result, ParameterOptions{Style: "form", Explode: true, Required: true})
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
}

func TestBindQueryParameter_OptionalMissing(t *testing.T) {
	vals := url.Values{}

	var result *string
	err := BindQueryParameter("missing", vals, &result, ParameterOptions{Style: "form", Explode: true})
	require.NoError(t, err)
	assert.Nil(t, result)
}

func TestBindQueryParameter_RequiredMissing(t *testing.T) {
	vals := url.Values{}

	var result string
	err := BindQueryParameter("required", vals, &result, ParameterOptions{Style: "form", Explode: true, Required: true})
	require.Error(t, err)
	var reqErr *MissingRequiredParameterError
	assert.True(t, errors.As(err, &reqErr))
}

func TestBindQueryParameter_DeepObject_Roundtrip(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		original := obj{Color: "blue", Size: 42}
		styled, err := StyleParameter("filter", original, ParameterOptions{Style: "deepObject", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result obj
		err = BindQueryParameter("filter", vals, &result, ParameterOptions{Style: "deepObject", Explode: true})
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("nested_struct", func(t *testing.T) {
		type inner struct {
			City string `json:"city"`
		}
		type outer struct {
			Name    string `json:"name"`
			Address inner  `json:"address"`
		}
		original := outer{Name: "alice", Address: inner{City: "NYC"}}
		styled, err := StyleParameter("user", original, ParameterOptions{Style: "deepObject", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result outer
		err = BindQueryParameter("user", vals, &result, ParameterOptions{Style: "deepObject", Explode: true})
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
	t.Run("with_slice", func(t *testing.T) {
		type obj struct {
			Tags []string `json:"tags"`
		}
		original := obj{Tags: []string{"a", "b"}}
		styled, err := StyleParameter("filter", original, ParameterOptions{Style: "deepObject", ParamLocation: ParamLocationQuery, Explode: true})
		require.NoError(t, err)

		vals, err := url.ParseQuery(styled)
		require.NoError(t, err)

		var result obj
		err = BindQueryParameter("filter", vals, &result, ParameterOptions{Style: "deepObject", Explode: true})
		require.NoError(t, err)
		assert.Equal(t, original, result)
	})
}
