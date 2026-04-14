package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStyleParameter_Simple(t *testing.T) {
	opts := func(extra ...func(*ParameterOptions)) ParameterOptions {
		o := ParameterOptions{Style: "simple", ParamLocation: ParamLocationPath}
		for _, f := range extra {
			f(&o)
		}
		return o
	}
	explode := func(o *ParameterOptions) { o.Explode = true }

	t.Run("primitive", func(t *testing.T) {
		result, err := StyleParameter("id", 5, opts())
		require.NoError(t, err)
		assert.Equal(t, "5", result)
	})
	t.Run("string", func(t *testing.T) {
		result, err := StyleParameter("name", "hello", opts())
		require.NoError(t, err)
		assert.Equal(t, "hello", result)
	})
	t.Run("string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts())
		require.NoError(t, err)
		assert.Equal(t, "a,b,c", result)
	})
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts())
		require.NoError(t, err)
		assert.Equal(t, "color,red,size,10", result)
	})
	t.Run("explode_primitive", func(t *testing.T) {
		result, err := StyleParameter("id", 5, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "5", result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "a,b,c", result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "color=red,size=10", result)
	})
}

func TestStyleParameter_Form(t *testing.T) {
	opts := func(extra ...func(*ParameterOptions)) ParameterOptions {
		o := ParameterOptions{Style: "form", ParamLocation: ParamLocationQuery}
		for _, f := range extra {
			f(&o)
		}
		return o
	}
	explode := func(o *ParameterOptions) { o.Explode = true }

	t.Run("primitive", func(t *testing.T) {
		result, err := StyleParameter("color", "blue", opts())
		require.NoError(t, err)
		assert.Equal(t, "color=blue", result)
	})
	t.Run("int", func(t *testing.T) {
		result, err := StyleParameter("count", 5, opts())
		require.NoError(t, err)
		assert.Equal(t, "count=5", result)
	})
	t.Run("string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts())
		require.NoError(t, err)
		assert.Equal(t, "tags=a,b,c", result)
	})
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts())
		require.NoError(t, err)
		assert.Equal(t, "filter=color,red,size,10", result)
	})
	t.Run("explode_primitive", func(t *testing.T) {
		result, err := StyleParameter("color", "blue", opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "color=blue", result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "tags=a&tags=b&tags=c", result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, "color=red&size=10", result)
	})
}

func TestStyleParameter_Label(t *testing.T) {
	opts := func(extra ...func(*ParameterOptions)) ParameterOptions {
		o := ParameterOptions{Style: "label", ParamLocation: ParamLocationPath}
		for _, f := range extra {
			f(&o)
		}
		return o
	}
	explode := func(o *ParameterOptions) { o.Explode = true }

	t.Run("primitive", func(t *testing.T) {
		result, err := StyleParameter("id", 5, opts())
		require.NoError(t, err)
		assert.Equal(t, ".5", result)
	})
	t.Run("string", func(t *testing.T) {
		result, err := StyleParameter("color", "blue", opts())
		require.NoError(t, err)
		assert.Equal(t, ".blue", result)
	})
	t.Run("string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts())
		require.NoError(t, err)
		assert.Equal(t, ".a,b,c", result)
	})
	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts())
		require.NoError(t, err)
		assert.Equal(t, ".color,red,size,10", result)
	})
	t.Run("explode_primitive", func(t *testing.T) {
		result, err := StyleParameter("id", 5, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, ".5", result)
	})
	t.Run("explode_string_slice", func(t *testing.T) {
		result, err := StyleParameter("tags", []string{"a", "b", "c"}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, ".a.b.c", result)
	})
	t.Run("explode_struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts(explode))
		require.NoError(t, err)
		assert.Equal(t, ".color=red.size=10", result)
	})
}

func TestStyleParameter_DeepObject(t *testing.T) {
	opts := ParameterOptions{Style: "deepObject", ParamLocation: ParamLocationQuery, Explode: true}

	t.Run("struct", func(t *testing.T) {
		type obj struct {
			Color string `json:"color"`
			Size  int    `json:"size"`
		}
		result, err := StyleParameter("filter", obj{Color: "red", Size: 10}, opts)
		require.NoError(t, err)
		assert.Equal(t, "filter[color]=red&filter[size]=10", result)
	})
	t.Run("nested_struct", func(t *testing.T) {
		type inner struct {
			City string `json:"city"`
		}
		type outer struct {
			Name    string `json:"name"`
			Address inner  `json:"address"`
		}
		result, err := StyleParameter("user", outer{
			Name:    "alice",
			Address: inner{City: "NYC"},
		}, opts)
		require.NoError(t, err)
		assert.Equal(t, "user[address][city]=NYC&user[name]=alice", result)
	})
	t.Run("with_slice", func(t *testing.T) {
		type obj struct {
			Tags []string `json:"tags"`
		}
		result, err := StyleParameter("filter", obj{Tags: []string{"a", "b"}}, opts)
		require.NoError(t, err)
		assert.Equal(t, "filter[tags][0]=a&filter[tags][1]=b", result)
	})
}
