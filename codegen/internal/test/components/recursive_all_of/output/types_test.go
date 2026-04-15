package output

import (
	"encoding/json"
	"testing"
)

// TestRecursiveObjectFields verifies that the recursive allOf self-reference
// is resolved without a stack overflow, and that fields from both the
// non-recursive base and the inline properties are present.
func TestRecursiveObjectFields(t *testing.T) {
	nonRec := "base"
	rec := "inline"
	obj := RecursiveObject{
		FieldInNonRecursive: &nonRec,
		FieldInRecursive:    &rec,
	}

	if *obj.FieldInNonRecursive != "base" {
		t.Errorf("FieldInNonRecursive = %q, want %q", *obj.FieldInNonRecursive, "base")
	}
	if *obj.FieldInRecursive != "inline" {
		t.Errorf("FieldInRecursive = %q, want %q", *obj.FieldInRecursive, "inline")
	}
}

func TestRecursiveObjectJSONRoundTrip(t *testing.T) {
	nonRec := "hello"
	rec := "world"
	original := RecursiveObject{
		FieldInNonRecursive: &nonRec,
		FieldInRecursive:    &rec,
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var decoded RecursiveObject
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if *decoded.FieldInNonRecursive != "hello" {
		t.Errorf("FieldInNonRecursive = %q, want %q", *decoded.FieldInNonRecursive, "hello")
	}
	if *decoded.FieldInRecursive != "world" {
		t.Errorf("FieldInRecursive = %q, want %q", *decoded.FieldInRecursive, "world")
	}
}

func TestNonRecursiveObjectFields(t *testing.T) {
	val := "test"
	obj := NonRecursiveObject{
		FieldInNonRecursive: &val,
	}
	if *obj.FieldInNonRecursive != "test" {
		t.Errorf("FieldInNonRecursive = %q, want %q", *obj.FieldInNonRecursive, "test")
	}
}

func TestApplyDefaults(t *testing.T) {
	r := &RecursiveObject{}
	r.ApplyDefaults()
	nr := &NonRecursiveObject{}
	nr.ApplyDefaults()
}
