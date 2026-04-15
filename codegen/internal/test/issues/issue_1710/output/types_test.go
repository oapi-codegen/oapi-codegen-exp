package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

// TestPromptOneOf0HasPromptField verifies the key regression: the 'prompt'
// field ([]ChatMessage) must be present on the chat variant after allOf/oneOf
// composition. This was the bug reported in issue 1710.
func TestPromptOneOf0HasPromptField(t *testing.T) {
	v := PromptOneOf0{
		Type:    ptr("chat"),
		Name:    "my-prompt",
		Version: 1,
		Prompt: []ChatMessage{
			{Role: "user", Content: "hello"},
		},
	}

	assert.Equal(t, "chat", *v.Type)
	assert.Equal(t, "my-prompt", v.Name)
	assert.Equal(t, 1, v.Version)
	require.Len(t, v.Prompt, 1)
	assert.Equal(t, "user", v.Prompt[0].Role)
	assert.Equal(t, "hello", v.Prompt[0].Content)
}

// TestPromptOneOf1HasPromptField verifies the text variant's string prompt field.
func TestPromptOneOf1HasPromptField(t *testing.T) {
	v := PromptOneOf1{
		Type:    ptr("text"),
		Name:    "text-prompt",
		Version: 2,
		Prompt:  "Write a poem",
	}

	assert.Equal(t, "text", *v.Type)
	assert.Equal(t, "Write a poem", v.Prompt)
}

// TestPromptUnionChatRoundTrip tests From/As round-trip for the chat variant.
func TestPromptUnionChatRoundTrip(t *testing.T) {
	chat := PromptOneOf0{
		Type:    ptr("chat"),
		Name:    "chat-prompt",
		Version: 1,
		Prompt: []ChatMessage{
			{Role: "system", Content: "You are helpful"},
			{Role: "user", Content: "Hi"},
		},
	}

	var u Prompt
	err := u.FromPromptOneOf0(chat)
	require.NoError(t, err)

	data, err := u.MarshalJSON()
	require.NoError(t, err)

	// Verify the prompt field is in the JSON
	var m map[string]any
	err = json.Unmarshal(data, &m)
	require.NoError(t, err)
	assert.Equal(t, "chat", m["type"])
	assert.Equal(t, "chat-prompt", m["name"])
	assert.NotNil(t, m["prompt"], "prompt field must be present in JSON")

	// Round-trip back
	var decoded Prompt
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	got, err := decoded.AsPromptOneOf0()
	require.NoError(t, err)
	assert.Equal(t, "chat-prompt", got.Name)
	assert.Equal(t, 1, got.Version)
	require.Len(t, got.Prompt, 2)
	assert.Equal(t, "system", got.Prompt[0].Role)
	assert.Equal(t, "user", got.Prompt[1].Role)
}

// TestPromptUnionTextRoundTrip tests From/As round-trip for the text variant.
func TestPromptUnionTextRoundTrip(t *testing.T) {
	text := PromptOneOf1{
		Type:    ptr("text"),
		Name:    "text-prompt",
		Version: 3,
		Prompt:  "Tell me a joke",
	}

	var u Prompt
	err := u.FromPromptOneOf1(text)
	require.NoError(t, err)

	data, err := u.MarshalJSON()
	require.NoError(t, err)

	var decoded Prompt
	err = decoded.UnmarshalJSON(data)
	require.NoError(t, err)

	got, err := decoded.AsPromptOneOf1()
	require.NoError(t, err)
	assert.Equal(t, "text-prompt", got.Name)
	assert.Equal(t, 3, got.Version)
	assert.Equal(t, "Tell me a joke", got.Prompt)
}

// TestTextPromptHasPromptField verifies the standalone TextPrompt type.
func TestTextPromptHasPromptField(t *testing.T) {
	tp := TextPrompt{Prompt: "my prompt", Name: "test", Version: 1}
	assert.Equal(t, "my prompt", tp.Prompt)
}

// TestChatPromptHasPromptField verifies the standalone ChatPrompt type.
func TestChatPromptHasPromptField(t *testing.T) {
	cp := ChatPrompt{
		Prompt: []ChatMessage{{Role: "assistant", Content: "hello"}},
		Name:   "test",
		Version: 1,
	}
	require.Len(t, cp.Prompt, 1)
	assert.Equal(t, "assistant", cp.Prompt[0].Role)
}

func TestDiscriminatorConstants(t *testing.T) {
	assert.Equal(t, PromptOneOf0AllOf0Type("chat"), Chat)
	assert.Equal(t, PromptOneOf1AllOf0Type("text"), Text)
}

func TestApplyDefaults(t *testing.T) {
	(&BasePrompt{}).ApplyDefaults()
	(&TextPrompt{}).ApplyDefaults()
	(&ChatMessage{}).ApplyDefaults()
	(&ChatPrompt{}).ApplyDefaults()
	(&Prompt{}).ApplyDefaults()
	(&PromptOneOf0{}).ApplyDefaults()
	(&PromptOneOf1{}).ApplyDefaults()
}
