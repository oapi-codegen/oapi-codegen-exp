package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEnumConstantValues verifies that all oneOf enum constants have the correct
// string values.
// https://github.com/oapi-codegen/oapi-codegen/issues/1029
func TestEnumConstantValues(t *testing.T) {
	assert.Equal(t, "undefined", string(Undefined))
	assert.Equal(t, "registered", string(Registered))
	assert.Equal(t, "pending", string(Pending))
	assert.Equal(t, "active", string(Active))
}

func TestRegistrationStateOneOf0_RoundTrip(t *testing.T) {
	var state RegistrationState
	err := state.FromRegistrationStateOneOf0(Undefined)
	require.NoError(t, err)

	data, err := json.Marshal(state)
	require.NoError(t, err)
	assert.JSONEq(t, `"undefined"`, string(data))

	var decoded RegistrationState
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	got, err := decoded.AsRegistrationStateOneOf0()
	require.NoError(t, err)
	assert.Equal(t, Undefined, got)
}

func TestRegistrationStateOneOf1_RoundTrip(t *testing.T) {
	var state RegistrationState
	err := state.FromRegistrationStateOneOf1(Registered)
	require.NoError(t, err)

	data, err := json.Marshal(state)
	require.NoError(t, err)
	assert.JSONEq(t, `"registered"`, string(data))

	var decoded RegistrationState
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	got, err := decoded.AsRegistrationStateOneOf1()
	require.NoError(t, err)
	assert.Equal(t, Registered, got)
}

func TestRegistrationStateOneOf2_RoundTrip(t *testing.T) {
	var state RegistrationState
	err := state.FromRegistrationStateOneOf2(Pending)
	require.NoError(t, err)

	data, err := json.Marshal(state)
	require.NoError(t, err)
	assert.JSONEq(t, `"pending"`, string(data))

	var decoded RegistrationState
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	got, err := decoded.AsRegistrationStateOneOf2()
	require.NoError(t, err)
	assert.Equal(t, Pending, got)
}

func TestRegistrationStateOneOf3_RoundTrip(t *testing.T) {
	var state RegistrationState
	err := state.FromRegistrationStateOneOf3(Active)
	require.NoError(t, err)

	data, err := json.Marshal(state)
	require.NoError(t, err)
	assert.JSONEq(t, `"active"`, string(data))

	var decoded RegistrationState
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	got, err := decoded.AsRegistrationStateOneOf3()
	require.NoError(t, err)
	assert.Equal(t, Active, got)
}

func TestRegistrationWithState_RoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*RegistrationState) error
		expected string
	}{
		{
			name: "undefined",
			setup: func(s *RegistrationState) error {
				return s.FromRegistrationStateOneOf0(Undefined)
			},
			expected: `{"state":"undefined"}`,
		},
		{
			name: "registered",
			setup: func(s *RegistrationState) error {
				return s.FromRegistrationStateOneOf1(Registered)
			},
			expected: `{"state":"registered"}`,
		},
		{
			name: "pending",
			setup: func(s *RegistrationState) error {
				return s.FromRegistrationStateOneOf2(Pending)
			},
			expected: `{"state":"pending"}`,
		},
		{
			name: "active",
			setup: func(s *RegistrationState) error {
				return s.FromRegistrationStateOneOf3(Active)
			},
			expected: `{"state":"active"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var state RegistrationState
			err := tt.setup(&state)
			require.NoError(t, err)

			reg := Registration{State: &state}

			data, err := json.Marshal(reg)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(data))

			var decoded Registration
			err = json.Unmarshal(data, &decoded)
			require.NoError(t, err)
			require.NotNil(t, decoded.State)

			reEncoded, err := json.Marshal(decoded)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(reEncoded))
		})
	}
}

func TestRegistrationWithNilState(t *testing.T) {
	reg := Registration{}

	data, err := json.Marshal(reg)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(data))

	var decoded Registration
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	assert.Nil(t, decoded.State)
}

func TestApplyDefaults(t *testing.T) {
	reg := Registration{}
	reg.ApplyDefaults()
	// ApplyDefaults is a no-op for this schema, but verify it doesn't panic.
	assert.Nil(t, reg.State)

	var state RegistrationState
	state.ApplyDefaults()
	// Also a no-op, verify no panic.
}
