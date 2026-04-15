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

// ===== Scenario 1: Array with anyOf items =====

func TestArrayOfAnyOf(t *testing.T) {
	t.Run("unmarshal string items", func(t *testing.T) {
		input := `["hello", "world"]`
		var arr ArrayOfAnyOf
		require.NoError(t, json.Unmarshal([]byte(input), &arr))
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
		require.NoError(t, json.Unmarshal([]byte(input), &arr))
		require.Len(t, arr, 1)

		obj, err := arr[0].AsArrayOfAnyOfAnyOf1()
		require.NoError(t, err)
		require.NotNil(t, obj.ID)
		assert.Equal(t, 42, *obj.ID)
	})

	t.Run("unmarshal mixed items", func(t *testing.T) {
		input := `["hello", {"id": 1}, "world", {"id": 2}]`
		var arr ArrayOfAnyOf
		require.NoError(t, json.Unmarshal([]byte(input), &arr))
		require.Len(t, arr, 4)

		s0, _ := arr[0].AsString0()
		assert.Equal(t, "hello", s0)
		obj1, _ := arr[1].AsArrayOfAnyOfAnyOf1()
		assert.Equal(t, 1, *obj1.ID)
		s2, _ := arr[2].AsString0()
		assert.Equal(t, "world", s2)
		obj3, _ := arr[3].AsArrayOfAnyOfAnyOf1()
		assert.Equal(t, 2, *obj3.ID)
	})

	t.Run("round trip mixed", func(t *testing.T) {
		var strItem ArrayOfAnyOfItem
		require.NoError(t, strItem.FromString0("test"))
		var objItem ArrayOfAnyOfItem
		require.NoError(t, objItem.FromArrayOfAnyOfAnyOf1(ArrayOfAnyOfAnyOf1{ID: ptr(99)}))

		data, err := json.Marshal(ArrayOfAnyOf{strItem, objItem})
		require.NoError(t, err)

		var decoded ArrayOfAnyOf
		require.NoError(t, json.Unmarshal(data, &decoded))
		require.Len(t, decoded, 2)

		s, _ := decoded[0].AsString0()
		assert.Equal(t, "test", s)
		obj, _ := decoded[1].AsArrayOfAnyOfAnyOf1()
		assert.Equal(t, 99, *obj.ID)
	})
}

// ===== Scenario 2: Object with anyOf property =====

func TestObjectWithAnyOfProperty(t *testing.T) {
	t.Run("unmarshal string value", func(t *testing.T) {
		var obj ObjectWithAnyOfProperty
		require.NoError(t, json.Unmarshal([]byte(`{"value": "hello"}`), &obj))
		require.NotNil(t, obj.Value)
		s, err := obj.Value.AsString0()
		require.NoError(t, err)
		assert.Equal(t, "hello", s)
	})

	t.Run("unmarshal integer value", func(t *testing.T) {
		var obj ObjectWithAnyOfProperty
		require.NoError(t, json.Unmarshal([]byte(`{"value": 42}`), &obj))
		require.NotNil(t, obj.Value)
		i, err := obj.Value.AsInt1()
		require.NoError(t, err)
		assert.Equal(t, 42, i)
	})

	t.Run("round trip string", func(t *testing.T) {
		var val ObjectWithAnyOfPropertyValue
		require.NoError(t, val.FromString0("test"))
		original := ObjectWithAnyOfProperty{Value: &val}
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded ObjectWithAnyOfProperty
		require.NoError(t, json.Unmarshal(data, &decoded))
		s, _ := decoded.Value.AsString0()
		assert.Equal(t, "test", s)
	})
}

// ===== Scenario 3: Object with oneOf property =====

func TestObjectWithOneOfProperty(t *testing.T) {
	t.Run("marshal variant 0", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		require.NoError(t, variant.FromObjectWithOneOfPropertyVariantOneOf0(ObjectWithOneOfPropertyVariantOneOf0{
			Kind: ptr("person"), Name: ptr("Alice"),
		}))
		data, err := json.Marshal(ObjectWithOneOfProperty{Variant: &variant})
		require.NoError(t, err)
		assert.JSONEq(t, `{"variant": {"kind": "person", "name": "Alice"}}`, string(data))
	})

	t.Run("marshal variant 1", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		require.NoError(t, variant.FromObjectWithOneOfPropertyVariantOneOf1(ObjectWithOneOfPropertyVariantOneOf1{
			Kind: ptr("counter"), Count: ptr(10),
		}))
		data, err := json.Marshal(ObjectWithOneOfProperty{Variant: &variant})
		require.NoError(t, err)
		assert.JSONEq(t, `{"variant": {"kind": "counter", "count": 10}}`, string(data))
	})

	t.Run("round trip variant 0", func(t *testing.T) {
		var variant ObjectWithOneOfPropertyVariant
		require.NoError(t, variant.FromObjectWithOneOfPropertyVariantOneOf0(ObjectWithOneOfPropertyVariantOneOf0{
			Kind: ptr("person"), Name: ptr("Bob"),
		}))
		data, err := json.Marshal(ObjectWithOneOfProperty{Variant: &variant})
		require.NoError(t, err)

		var decoded ObjectWithOneOfProperty
		require.NoError(t, json.Unmarshal(data, &decoded))
		v0, _ := decoded.Variant.AsObjectWithOneOfPropertyVariantOneOf0()
		assert.Equal(t, "Bob", *v0.Name)
	})
}

// ===== Scenario 4: allOf containing oneOf =====

func TestAllOfWithOneOf(t *testing.T) {
	t.Run("marshal with optionA", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		require.NoError(t, union.FromAllOfWithOneOfAllOf1OneOf0(AllOfWithOneOfAllOf1OneOf0{OptionA: ptr(true)}))
		data, err := json.Marshal(AllOfWithOneOf{Base: ptr("test"), AllOfWithOneOfAllOf1: &union})
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(data, &m))
		assert.Equal(t, "test", m["base"])
		assert.Equal(t, true, m["optionA"])
	})

	t.Run("marshal with optionB", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		require.NoError(t, union.FromAllOfWithOneOfAllOf1OneOf1(AllOfWithOneOfAllOf1OneOf1{OptionB: ptr(42)}))
		data, err := json.Marshal(AllOfWithOneOf{Base: ptr("test"), AllOfWithOneOfAllOf1: &union})
		require.NoError(t, err)

		var m map[string]any
		require.NoError(t, json.Unmarshal(data, &m))
		assert.Equal(t, "test", m["base"])
		assert.Equal(t, float64(42), m["optionB"])
	})

	t.Run("unmarshal with optionA", func(t *testing.T) {
		var obj AllOfWithOneOf
		require.NoError(t, json.Unmarshal([]byte(`{"base": "test", "optionA": true}`), &obj))
		assert.Equal(t, "test", *obj.Base)
		require.NotNil(t, obj.AllOfWithOneOfAllOf1)
		v0, _ := obj.AllOfWithOneOfAllOf1.AsAllOfWithOneOfAllOf1OneOf0()
		assert.Equal(t, true, *v0.OptionA)
	})

	t.Run("round trip with optionA", func(t *testing.T) {
		var union AllOfWithOneOfAllOf1
		require.NoError(t, union.FromAllOfWithOneOfAllOf1OneOf0(AllOfWithOneOfAllOf1OneOf0{OptionA: ptr(false)}))
		original := AllOfWithOneOf{Base: ptr("round"), AllOfWithOneOfAllOf1: &union}
		data, err := json.Marshal(original)
		require.NoError(t, err)

		var decoded AllOfWithOneOf
		require.NoError(t, json.Unmarshal(data, &decoded))
		assert.Equal(t, "round", *decoded.Base)
		v0, _ := decoded.AllOfWithOneOfAllOf1.AsAllOfWithOneOfAllOf1OneOf0()
		assert.Equal(t, false, *v0.OptionA)
	})
}

// ===== Scenario 5: oneOf with nested allOf, field preservation =====

func TestPromptOneOf0HasPromptField(t *testing.T) {
	v := PromptOneOf0{
		Type: ptr("chat"), Name: "my-prompt", Version: 1,
		Prompt: []ChatMessage{{Role: "user", Content: "hello"}},
	}
	assert.Equal(t, "chat", *v.Type)
	require.Len(t, v.Prompt, 1)
	assert.Equal(t, "user", v.Prompt[0].Role)
}

func TestPromptOneOf1HasPromptField(t *testing.T) {
	v := PromptOneOf1{Type: ptr("text"), Name: "text-prompt", Version: 2, Prompt: "Write a poem"}
	assert.Equal(t, "text", *v.Type)
	assert.Equal(t, "Write a poem", v.Prompt)
}

func TestPromptUnionChatRoundTrip(t *testing.T) {
	var u Prompt
	require.NoError(t, u.FromPromptOneOf0(PromptOneOf0{
		Type: ptr("chat"), Name: "chat-prompt", Version: 1,
		Prompt: []ChatMessage{{Role: "system", Content: "You are helpful"}, {Role: "user", Content: "Hi"}},
	}))
	data, err := u.MarshalJSON()
	require.NoError(t, err)

	var m map[string]any
	require.NoError(t, json.Unmarshal(data, &m))
	assert.NotNil(t, m["prompt"], "prompt field must be present in JSON")

	var decoded Prompt
	require.NoError(t, decoded.UnmarshalJSON(data))
	got, _ := decoded.AsPromptOneOf0()
	assert.Equal(t, "chat-prompt", got.Name)
	require.Len(t, got.Prompt, 2)
}

func TestPromptUnionTextRoundTrip(t *testing.T) {
	var u Prompt
	require.NoError(t, u.FromPromptOneOf1(PromptOneOf1{
		Type: ptr("text"), Name: "text-prompt", Version: 3, Prompt: "Tell me a joke",
	}))
	data, err := u.MarshalJSON()
	require.NoError(t, err)

	var decoded Prompt
	require.NoError(t, decoded.UnmarshalJSON(data))
	got, _ := decoded.AsPromptOneOf1()
	assert.Equal(t, "Tell me a joke", got.Prompt)
}

func TestPromptDiscriminatorConstants(t *testing.T) {
	assert.Equal(t, PromptOneOf0AllOf0Type("chat"), Chat)
	assert.Equal(t, PromptOneOf1AllOf0Type("text"), Text)
}

// ===== Scenario 6: anyOf/allOf/oneOf with duplicate enum variants =====

func TestCompositionEnumFieldA_StringRoundTrip(t *testing.T) {
	var fa CompositionEnumTestFieldA
	require.NoError(t, fa.FromString0("plain-string"))
	data, err := fa.MarshalJSON()
	require.NoError(t, err)

	var fa2 CompositionEnumTestFieldA
	require.NoError(t, fa2.UnmarshalJSON(data))
	got, err := fa2.AsString0()
	require.NoError(t, err)
	assert.Equal(t, "plain-string", got)
}

func TestCompositionEnumFieldA_EnumRoundTrip(t *testing.T) {
	var fa CompositionEnumTestFieldA
	require.NoError(t, fa.FromCompositionEnumTestFieldAAnyOf1(CompositionEnumTestFieldAAnyOf1Foo))
	data, err := fa.MarshalJSON()
	require.NoError(t, err)

	var fa2 CompositionEnumTestFieldA
	require.NoError(t, fa2.UnmarshalJSON(data))
	got, err := fa2.AsCompositionEnumTestFieldAAnyOf1()
	require.NoError(t, err)
	assert.Equal(t, CompositionEnumTestFieldAAnyOf1Foo, got)
}

func TestCompositionEnumFieldC_StringRoundTrip(t *testing.T) {
	var fc CompositionEnumTestFieldC
	require.NoError(t, fc.FromString0("one-of-string"))
	data, err := fc.MarshalJSON()
	require.NoError(t, err)

	var fc2 CompositionEnumTestFieldC
	require.NoError(t, fc2.UnmarshalJSON(data))
	got, err := fc2.AsString0()
	require.NoError(t, err)
	assert.Equal(t, "one-of-string", got)
}

func TestCompositionEnumFieldC_EnumRoundTrip(t *testing.T) {
	var fc CompositionEnumTestFieldC
	require.NoError(t, fc.FromCompositionEnumTestFieldCOneOf1(CompositionEnumTestFieldCOneOf1Bar))
	data, err := fc.MarshalJSON()
	require.NoError(t, err)

	var fc2 CompositionEnumTestFieldC
	require.NoError(t, fc2.UnmarshalJSON(data))
	got, err := fc2.AsCompositionEnumTestFieldCOneOf1()
	require.NoError(t, err)
	assert.Equal(t, CompositionEnumTestFieldCOneOf1Bar, got)
}

func TestCompositionEnumConstants(t *testing.T) {
	assert.Equal(t, CompositionEnumTestFieldAAnyOf1("foo"), CompositionEnumTestFieldAAnyOf1Foo)
	assert.Equal(t, CompositionEnumTestFieldAAnyOf1("bar"), CompositionEnumTestFieldAAnyOf1Bar)
	assert.Equal(t, CompositionEnumTestFieldCOneOf1("foo"), CompositionEnumTestFieldCOneOf1Foo)
	assert.Equal(t, CompositionEnumTestFieldCOneOf1("bar"), CompositionEnumTestFieldCOneOf1Bar)
}

func TestCompositionEnumFullStructRoundTrip(t *testing.T) {
	var fieldA CompositionEnumTestFieldA
	require.NoError(t, fieldA.FromString0("hello"))
	var fieldC CompositionEnumTestFieldC
	require.NoError(t, fieldC.FromCompositionEnumTestFieldCOneOf1(CompositionEnumTestFieldCOneOf1Foo))

	original := CompositionEnumTest{
		FieldA: &fieldA,
		FieldB: &CompositionEnumTestFieldB{},
		FieldC: &fieldC,
	}
	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded CompositionEnumTest
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotA, _ := decoded.FieldA.AsString0()
	assert.Equal(t, "hello", gotA)
	require.NotNil(t, decoded.FieldB)
	gotC, _ := decoded.FieldC.AsCompositionEnumTestFieldCOneOf1()
	assert.Equal(t, CompositionEnumTestFieldCOneOf1Foo, gotC)
}

// ===== Shared =====

func TestApplyDefaults(t *testing.T) {
	// Scenario 1-4
	(&ArrayOfAnyOfItem{}).ApplyDefaults()
	(&ArrayOfAnyOfAnyOf1{}).ApplyDefaults()
	(&ObjectWithAnyOfProperty{}).ApplyDefaults()
	(&ObjectWithAnyOfPropertyValue{}).ApplyDefaults()
	(&ObjectWithOneOfProperty{}).ApplyDefaults()
	(&ObjectWithOneOfPropertyVariant{}).ApplyDefaults()
	(&ObjectWithOneOfPropertyVariantOneOf0{}).ApplyDefaults()
	(&ObjectWithOneOfPropertyVariantOneOf1{}).ApplyDefaults()
	(&AllOfWithOneOf{}).ApplyDefaults()
	(&AllOfWithOneOfAllOf1{}).ApplyDefaults()
	(&AllOfWithOneOfAllOf1OneOf0{}).ApplyDefaults()
	(&AllOfWithOneOfAllOf1OneOf1{}).ApplyDefaults()
	// Scenario 5
	(&BasePrompt{}).ApplyDefaults()
	(&TextPrompt{}).ApplyDefaults()
	(&ChatMessage{}).ApplyDefaults()
	(&ChatPrompt{}).ApplyDefaults()
	(&Prompt{}).ApplyDefaults()
	(&PromptOneOf0{}).ApplyDefaults()
	(&PromptOneOf1{}).ApplyDefaults()
	// Scenario 6
	(&CompositionEnumTest{}).ApplyDefaults()
	(&CompositionEnumTestFieldA{}).ApplyDefaults()
	(&CompositionEnumTestFieldB{}).ApplyDefaults()
	(&CompositionEnumTestFieldC{}).ApplyDefaults()
}
