package dce

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEliminateDeadCode_RemovesOrphanedComments(t *testing.T) {
	src := `package gen

func usedFunc() {
	helperFunc()
}

// --- oapi-runtime begin ---

// helperFunc is used by usedFunc.
func helperFunc() {}

// unusedFunc does nothing useful.
// It has a multi-line doc comment.
func unusedFunc() {}

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	assert.Contains(t, result, "helperFunc")
	assert.NotContains(t, result, "unusedFunc")
}

func TestEliminateDeadCode_NoMarkers(t *testing.T) {
	src := `package gen

func foo() {}
`
	result, err := EliminateDeadCode(src)
	require.NoError(t, err)
	assert.Equal(t, src, result)
}

func TestEliminateDeadCode_AllUsed(t *testing.T) {
	src := `package gen

func caller() {
	a()
	b()
}

// --- oapi-runtime begin ---

// a does something.
func a() {}

// b does something else.
func b() {}

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	assert.Contains(t, result, "func a()")
	assert.Contains(t, result, "func b()")
	assert.Contains(t, result, "a does something")
	assert.Contains(t, result, "b does something else")
}

func TestEliminateDeadCode_NoneUsed(t *testing.T) {
	src := `package gen

func caller() {}

// --- oapi-runtime begin ---

// orphanA is unused.
func orphanA() {}

// orphanB is also unused.
func orphanB() {}

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	assert.NotContains(t, result, "orphanA")
	assert.NotContains(t, result, "orphanB")
}

func TestEliminateDeadCode_InlineComments(t *testing.T) {
	src := `package gen

func caller() {
	used()
}

// --- oapi-runtime begin ---

func used() {} // keep this

func unused() {} // drop this

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	assert.Contains(t, result, "keep this")
	assert.NotContains(t, result, "drop this")
}

func TestEliminateDeadCode_TypeDecl(t *testing.T) {
	src := `package gen

// UsedType is referenced.
type UsedType struct{}

// --- oapi-runtime begin ---

// helperType is used by UsedType (transitively).
type helperType = UsedType

// unusedType is not referenced.
type unusedType struct {
	field string
}

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	// helperType is reachable because it references UsedType which is in roots,
	// and helperType itself may or may not be reachable depending on the algorithm.
	// The key assertion: unusedType and its comment should be gone.
	assert.NotContains(t, result, "unusedType")
	assert.NotContains(t, result, "unusedType is not referenced")
}

func TestEliminateDeadCode_PreservesNonRuntimeComments(t *testing.T) {
	src := `package gen

// Package-level comment that should survive.

// caller calls helper.
func caller() {
	helper()
}

// --- oapi-runtime begin ---

// helper is needed.
func helper() {}

// unused is not needed.
func unused() {}

// --- oapi-runtime end ---
`

	result, err := EliminateDeadCode(src)
	require.NoError(t, err)

	assert.Contains(t, result, "Package-level comment")
	assert.Contains(t, result, "caller calls helper")
	assert.Contains(t, result, "helper is needed")
	assert.NotContains(t, result, "unused is not needed")
	// The word "unused" should only not appear as a function or its comment
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "unused") && !strings.Contains(trimmed, "nolint") {
			t.Errorf("unexpected reference to 'unused' in output: %s", trimmed)
		}
	}
}
