package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

// TestArrayOfAnyOf tests marshaling/unmarshaling of arrays with anyOf items.
func TestArrayOfAnyOf(t *testing.T) {
	t.Run("unmarshal string items", func(t *testing.T) {
		input := `["hello", "world"]`
		var arr ArrayOfAnyOf
		err := json.Unmarshal([]byte(input), &arr)
		require.NoError(t, err)
		require.Len(t, arr, 2)

		s0, err := arr[0].AsString0()
		require.NoError(t, err)
		assert.Equal(t, "hello", s0)

		s1, err := arr[1].AsString0()
		require.NoError(t, err)
		assert.Equal(t, "world", s1)
	})

	t.Run("unmarshal object item", func(t *testing.T) {
		input := `[{"id": 42}]`
		var arr ArrayOfAnyOf
		err := json.Unmarshal([]byte(input), &arr)
		require.NoError(t, err)
		require.Len(t, arr, 1)

		obj, err := arr[0].AsArrayOfAnyOfAnyOf1()
		require.NoError(t, err)
		require.NotNil(t, obj.ID)
		assert.Equal(t, 42, *obj.ID)
	})

	t.Run("unmarshal mixed items", func(t *testing.T) {
		input := `["hello", {"id": 1}, "world", {"id": 2}]`
		var arr ArrayOfAnyOf
		err := json.Unmarshal([]byte(input), &arr)
		require.NoError(t, err)
		require.Len(t, arr, 4)

		s0, err := arr[0].AsString0()
		require.NoError(t, err)
		assert.Equal(t, "hello", s0)

		obj1, err := arr[1].AsArrayOfAnyOfAnyOf1()
		require.NoError(t, err)
		require.NotNil(t, obj1.ID)
		assert.Equal(t, 1, *obj1.ID)

		s2, err := arr[2].AsString0()
		require.NoError(t, err)
		assert.Equal(t, "world", s2)

		obj3, err := arr[3].AsArrayOfAnyOfAnyOf1()
		require.NoError(t, err)
		require.NotNil(t, obj3.ID)
		assert.Equal(t, 2, *obj3.ID)
	})

	t.Run("marshal string item", func(t *testing.T) {
		var item ArrayOfAnyOfItem
		err := item.FromString0("hello")
		require.NoError(t, err)

		arr := ArrayOfAnyOf{item}
		data, err := json.Marshal(arr)
		require.NoError(t, err)
		assert.JSONEq(t, `["hello"]`, string(data))
	})

	t.Run("marshal object item", func(t *testing.T) {
		var item ArrayOfAnyOfItem
		err := item.FromArrayOfAnyOfAnyOf1(ArrayOfAnyOfAnyOf1{ID: ptr(42)})
		require.NoError(t, err)

		arr := ArrayOfAnyOf{item}
		data, err := json.Marshal(arr)
		require.NoError(t, err)
		assert.JSONEq(t, `[{"id": 42}]`, string(data))
	})

	t.Run("round trip mixed", func(t *testing.T) {
		var strItem ArrayOfAnyOfItem
		err := strItem.FromString0("test")
		require.NoError(t, err)

		var objItem ArrayOfAnyOfItem
		err = objItem.FromArrayOfAnyOfAnyOf1(ArrayOfAnyOfAnyOf1{ID: ptr(99)})
		require.NoError(t, err)

		original := ArrayOfAnyOf{strItem, objItem}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded ArrayOfAnyOf
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.Len(t, decoded, 2)

		s, err := decoded[0].AsString0()
		require.NoError(t, err)
		assert.Equal(t, "test", s)

		obj, err := decoded[1].AsArrayOfAnyOfAnyOf1()
		require.NoError(t, err)
		require.NotNil(t, obj.ID)
		assert.Equal(t, 99, *obj.ID)
	})
}

// TestObjectWithAnyOfProperty tests marshaling/unmarshaling of objects with anyOf properties.
func TestObjectWithAnyOfProperty(t *testing.T) {
	t.Run("unmarshal string value", func(t *testing.T) {
		input := `{"value": "hello"}`
		var obj ObjectWithAnyOfProperty
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)
		require.NotNil(t, obj.Value)

		s, err := obj.Value.AsString0()
		require.NoError(t, err)
		assert.Equal(t, "hello", s)
	})

	t.Run("unmarshal integer value", func(t *testing.T) {
		input := `{"value": 42}`
		var obj ObjectWithAnyOfProperty
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)
		require.NotNil(t, obj.Value)

		i, err := obj.Value.AsInt1()
		require.NoError(t, err)
		assert.Equal(t, 42, i)
	})

	t.Run("marshal string value", func(t *testing.T) {
		var val ObjectWithAnyOfPropertyValue
		err := val.FromString0("hello")
		require.NoError(t, err)

		obj := ObjectWithAnyOfProperty{Value: &val}
		data, err := json.Marshal(obj)
		require.NoError(t, err)
		assert.JSONEq(t, `{"value": "hello"}`, string(data))
	})

	t.Run("marshal integer value", func(t *testing.T) {
		var val ObjectWithAnyOfPropertyValue
		err := val.FromInt1(42)
		require.NoError(t, err)

		obj := ObjectWithAnyOfProperty{Value: &val}
		data, err := json.Marshal(obj)
		require.NoError(t, err)
		assert.JSONEq(t, `{"value": 42}`, string(data))
	})

	t.Run("round trip string", func(t *testing.T) {
		var val ObjectWithAnyOfPropertyValue
		err := val.FromString0("test")
		require.NoError(t, err)

		original := ObjectWithAnyOfProperty{Value: &val}
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded ObjectWithAnyOfProperty
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.NotNil(t, decoded.Value)

		s, err := decoded.Value.AsString0()
		require.NoError(t, err)
		assert.Equal(t, "test", s)
	})
}

// TestObjectWithOneOfProperty tests marshaling/unmarshaling of objects with oneOf properties.
func TestObjectWithOneOfProperty(t *testing.T) {
	t.Run("marshal variant 0", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		err := variant.FromObjectWithOneOfPropertyVariantOneOf0(ObjectWithOneOfPropertyVariantOneOf0{
			Kind: ptr("person"),
			Name: ptr("Alice"),
		})
		require.NoError(t, err)

		obj := ObjectWithOneOfProperty{Variant: &variant}
		data, err := json.Marshal(obj)
		require.NoError(t, err)
		assert.JSONEq(t, `{"variant": {"kind": "person", "name": "Alice"}}`, string(data))
	})

	t.Run("marshal variant 1", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		err := variant.FromObjectWithOneOfPropertyVariantOneOf1(ObjectWithOneOfPropertyVariantOneOf1{
			Kind:  ptr("counter"),
			Count: ptr(10),
		})
		require.NoError(t, err)

		obj := ObjectWithOneOfProperty{Variant: &variant}
		data, err := json.Marshal(obj)
		require.NoError(t, err)
		assert.JSONEq(t, `{"variant": {"kind": "counter", "count": 10}}`, string(data))
	})

	t.Run("unmarshal and extract variant 0", func(t *testing.T) {
		input := `{"variant": {"kind": "person", "name": "Alice"}}`
		var obj ObjectWithOneOfProperty
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)
		require.NotNil(t, obj.Variant)

		v0, err := obj.Variant.AsObjectWithOneOfPropertyVariantOneOf0()
		require.NoError(t, err)
		require.NotNil(t, v0.Kind)
		assert.Equal(t, "person", *v0.Kind)
		require.NotNil(t, v0.Name)
		assert.Equal(t, "Alice", *v0.Name)
	})

	t.Run("unmarshal and extract variant 1", func(t *testing.T) {
		input := `{"variant": {"kind": "counter", "count": 10}}`
		var obj ObjectWithOneOfProperty
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)
		require.NotNil(t, obj.Variant)

		v1, err := obj.Variant.AsObjectWithOneOfPropertyVariantOneOf1()
		require.NoError(t, err)
		require.NotNil(t, v1.Kind)
		assert.Equal(t, "counter", *v1.Kind)
		require.NotNil(t, v1.Count)
		assert.Equal(t, 10, *v1.Count)
	})

	t.Run("unmarshal ambiguous input still stores raw JSON", func(t *testing.T) {
		// With V3 union types, UnmarshalJSON stores raw bytes without
		// oneOf validation. Both variants can be extracted -- the caller
		// is responsible for choosing the right As* method.
		input := `{"variant": {"kind": "ambiguous"}}`
		var obj ObjectWithOneOfProperty
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)
		require.NotNil(t, obj.Variant)

		// Both As* calls succeed because the raw JSON is valid for either.
		v0, err := obj.Variant.AsObjectWithOneOfPropertyVariantOneOf0()
		require.NoError(t, err)
		assert.Equal(t, "ambiguous", *v0.Kind)

		v1, err := obj.Variant.AsObjectWithOneOfPropertyVariantOneOf1()
		require.NoError(t, err)
		assert.Equal(t, "ambiguous", *v1.Kind)
	})

	t.Run("round trip variant 0", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		err := variant.FromObjectWithOneOfPropertyVariantOneOf0(ObjectWithOneOfPropertyVariantOneOf0{
			Kind: ptr("person"),
			Name: ptr("Bob"),
		})
		require.NoError(t, err)

		original := ObjectWithOneOfProperty{Variant: &variant}
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded ObjectWithOneOfProperty
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)
		require.NotNil(t, decoded.Variant)

		v0, err := decoded.Variant.AsObjectWithOneOfPropertyVariantOneOf0()
		require.NoError(t, err)
		assert.Equal(t, "Bob", *v0.Name)
	})
}

// TestAllOfWithOneOf tests marshaling/unmarshaling of allOf containing oneOf.
func TestAllOfWithOneOf(t *testing.T) {
	t.Run("marshal with optionA", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		err := union.FromAllOfWithOneOfAllOf1OneOf0(AllOfWithOneOfAllOf1OneOf0{
			OptionA: ptr(true),
		})
		require.NoError(t, err)

		obj := AllOfWithOneOf{
			Base:                 ptr("test"),
			AllOfWithOneOfAllOf1: &union,
		}

		data, err := json.Marshal(obj)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "test", m["base"])
		assert.Equal(t, true, m["optionA"])
	})

	t.Run("marshal with optionB", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		err := union.FromAllOfWithOneOfAllOf1OneOf1(AllOfWithOneOfAllOf1OneOf1{
			OptionB: ptr(42),
		})
		require.NoError(t, err)

		obj := AllOfWithOneOf{
			Base:                 ptr("test"),
			AllOfWithOneOfAllOf1: &union,
		}

		data, err := json.Marshal(obj)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "test", m["base"])
		assert.Equal(t, float64(42), m["optionB"]) // JSON numbers are float64
	})

	t.Run("marshal with nil union", func(t *testing.T) {
		obj := AllOfWithOneOf{
			Base: ptr("only-base"),
		}

		data, err := json.Marshal(obj)
		require.NoError(t, err)

		var m map[string]any
		err = json.Unmarshal(data, &m)
		require.NoError(t, err)

		assert.Equal(t, "only-base", m["base"])
		assert.NotContains(t, m, "optionA")
		assert.NotContains(t, m, "optionB")
	})

	t.Run("unmarshal with optionA", func(t *testing.T) {
		input := `{"base": "test", "optionA": true}`
		var obj AllOfWithOneOf
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)

		require.NotNil(t, obj.Base)
		assert.Equal(t, "test", *obj.Base)

		require.NotNil(t, obj.AllOfWithOneOfAllOf1)
		v0, err := obj.AllOfWithOneOfAllOf1.AsAllOfWithOneOfAllOf1OneOf0()
		require.NoError(t, err)
		require.NotNil(t, v0.OptionA)
		assert.Equal(t, true, *v0.OptionA)
	})

	t.Run("unmarshal with optionB", func(t *testing.T) {
		input := `{"base": "test", "optionB": 42}`
		var obj AllOfWithOneOf
		err := json.Unmarshal([]byte(input), &obj)
		require.NoError(t, err)

		require.NotNil(t, obj.Base)
		assert.Equal(t, "test", *obj.Base)

		require.NotNil(t, obj.AllOfWithOneOfAllOf1)
		v1, err := obj.AllOfWithOneOfAllOf1.AsAllOfWithOneOfAllOf1OneOf1()
		require.NoError(t, err)
		require.NotNil(t, v1.OptionB)
		assert.Equal(t, 42, *v1.OptionB)
	})

	t.Run("round trip with optionA", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		err := union.FromAllOfWithOneOfAllOf1OneOf0(AllOfWithOneOfAllOf1OneOf0{
			OptionA: ptr(false),
		})
		require.NoError(t, err)

		original := AllOfWithOneOf{
			Base:                 ptr("round"),
			AllOfWithOneOfAllOf1: &union,
		}

		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded AllOfWithOneOf
		err = json.Unmarshal(data, &decoded)
		require.NoError(t, err)

		require.NotNil(t, decoded.Base)
		assert.Equal(t, "round", *decoded.Base)

		require.NotNil(t, decoded.AllOfWithOneOfAllOf1)
		v0, err := decoded.AllOfWithOneOfAllOf1.AsAllOfWithOneOfAllOf1OneOf0()
		require.NoError(t, err)
		require.NotNil(t, v0.OptionA)
		assert.Equal(t, false, *v0.OptionA)
	})
}

// TestApplyDefaults verifies that ApplyDefaults is callable on all types without panicking.
func TestApplyDefaults(t *testing.T) {
	t.Run("ArrayOfAnyOfItem", func(t *testing.T) {
		var item ArrayOfAnyOfItem
		item.ApplyDefaults() // should not panic
	})

	t.Run("ArrayOfAnyOfAnyOf1", func(t *testing.T) {
		s := ArrayOfAnyOfAnyOf1{}
		s.ApplyDefaults()
	})

	t.Run("ObjectWithAnyOfProperty", func(t *testing.T) {
		s := ObjectWithAnyOfProperty{}
		s.ApplyDefaults()
	})

	t.Run("ObjectWithAnyOfPropertyValue", func(t *testing.T) {
		var v ObjectWithAnyOfPropertyValue
		v.ApplyDefaults()
	})

	t.Run("ObjectWithOneOfProperty", func(t *testing.T) {
		s := ObjectWithOneOfProperty{}
		s.ApplyDefaults()
	})

	t.Run("ObjectWithOneOfPropertyVariant", func(t *testing.T) {
		var v ObjectWithOneOfPropertyVariant
		v.ApplyDefaults()
	})

	t.Run("ObjectWithOneOfPropertyVariantOneOf0", func(t *testing.T) {
		s := ObjectWithOneOfPropertyVariantOneOf0{}
		s.ApplyDefaults()
	})

	t.Run("ObjectWithOneOfPropertyVariantOneOf1", func(t *testing.T) {
		s := ObjectWithOneOfPropertyVariantOneOf1{}
		s.ApplyDefaults()
	})

	t.Run("AllOfWithOneOf", func(t *testing.T) {
		s := AllOfWithOneOf{}
		s.ApplyDefaults()
	})

	t.Run("AllOfWithOneOfAllOf1", func(t *testing.T) {
		var v AllOfWithOneOfAllOf1
		v.ApplyDefaults()
	})

	t.Run("AllOfWithOneOfAllOf1OneOf0", func(t *testing.T) {
		s := AllOfWithOneOfAllOf1OneOf0{}
		s.ApplyDefaults()
	})

	t.Run("AllOfWithOneOfAllOf1OneOf1", func(t *testing.T) {
		s := AllOfWithOneOfAllOf1OneOf1{}
		s.ApplyDefaults()
	})
}
