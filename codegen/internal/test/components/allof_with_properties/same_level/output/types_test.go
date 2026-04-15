package output

import (
	"encoding/json"
	"testing"
)

// TestXHasAllFields verifies that schema X has properties from both its own
// definition AND from the allOf reference to YBase.
func TestXHasAllFields(t *testing.T) {
	a := "a-value"
	b := 42
	baseField := "base-value"

	x := X{
		A:         &a,
		B:         &b,
		BaseField: &baseField,
	}

	if *x.A != "a-value" {
		t.Errorf("X.A = %q, want %q", *x.A, "a-value")
	}
	if *x.B != 42 {
		t.Errorf("X.B = %d, want %d", *x.B, 42)
	}
	if *x.BaseField != "base-value" {
		t.Errorf("X.BaseField = %q, want %q", *x.BaseField, "base-value")
	}
}

func TestXJSONRoundTrip(t *testing.T) {
	a := "a-value"
	b := 42
	baseField := "base-value"

	original := X{A: &a, B: &b, BaseField: &baseField}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded X
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if *decoded.A != *original.A || *decoded.B != *original.B || *decoded.BaseField != *original.BaseField {
		t.Errorf("Round trip failed: got %+v, want %+v", decoded, original)
	}
}

// TestBarHasBothProperties verifies that Bar has both foo (from allOf) and bar
// (from direct properties).
func TestBarHasBothProperties(t *testing.T) {
	bar := Bar{Foo: "test-foo", Bar: "test-bar"}

	data, err := json.Marshal(bar)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded Bar
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Foo != "test-foo" {
		t.Errorf("Foo = %q, want %q", decoded.Foo, "test-foo")
	}
	if decoded.Bar != "test-bar" {
		t.Errorf("Bar = %q, want %q", decoded.Bar, "test-bar")
	}
}

// TestBarRequiredFields verifies that both foo and bar are required.
func TestBarRequiredFields(t *testing.T) {
	bar := Bar{}
	data, err := json.Marshal(bar)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	expected := `{"bar":"","foo":""}`
	if string(data) != expected {
		t.Errorf("got %s, want %s", string(data), expected)
	}
}

func TestApplyDefaults(t *testing.T) {
	(&X{}).ApplyDefaults()
	(&YBase{}).ApplyDefaults()
	(&Foo{}).ApplyDefaults()
	(&Bar{}).ApplyDefaults()
}
