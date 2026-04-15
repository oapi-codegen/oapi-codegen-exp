package output

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilterValue_StringVariant(t *testing.T) {
	var fv FilterValue
	require.NoError(t, fv.FromString1("hello"))

	got, err := fv.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "hello", got)

	data, err := json.Marshal(fv)
	require.NoError(t, err)
	assert.JSONEq(t, `"hello"`, string(data))

	var decoded FilterValue
	require.NoError(t, json.Unmarshal(data, &decoded))
	got2, err := decoded.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "hello", got2)
}

func TestFilterValue_FloatVariant(t *testing.T) {
	var fv FilterValue
	require.NoError(t, fv.FromFloat320(3.14))

	got, err := fv.AsFloat320()
	require.NoError(t, err)
	assert.InDelta(t, float32(3.14), got, 0.001)

	data, err := json.Marshal(fv)
	require.NoError(t, err)

	var decoded FilterValue
	require.NoError(t, json.Unmarshal(data, &decoded))
	got2, err := decoded.AsFloat320()
	require.NoError(t, err)
	assert.InDelta(t, float32(3.14), got2, 0.001)
}

func TestFilterValue_BoolVariant(t *testing.T) {
	var fv FilterValue
	require.NoError(t, fv.FromBool2(true))

	got, err := fv.AsBool2()
	require.NoError(t, err)
	assert.True(t, got)

	data, err := json.Marshal(fv)
	require.NoError(t, err)
	assert.Equal(t, "true", string(data))

	var decoded FilterValue
	require.NoError(t, json.Unmarshal(data, &decoded))
	got2, err := decoded.AsBool2()
	require.NoError(t, err)
	assert.True(t, got2)
}

func TestFilterRangeValue_StringVariant(t *testing.T) {
	var frv FilterRangeValue
	require.NoError(t, frv.FromString1("2023-01-01"))

	got, err := frv.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "2023-01-01", got)

	data, err := json.Marshal(frv)
	require.NoError(t, err)
	assert.JSONEq(t, `"2023-01-01"`, string(data))
}

func TestFilterRangeValue_FloatVariant(t *testing.T) {
	var frv FilterRangeValue
	require.NoError(t, frv.FromFloat320(99.5))

	got, err := frv.AsFloat320()
	require.NoError(t, err)
	assert.InDelta(t, float32(99.5), got, 0.001)

	data, err := json.Marshal(frv)
	require.NoError(t, err)

	var decoded FilterRangeValue
	require.NoError(t, json.Unmarshal(data, &decoded))
	got2, err := decoded.AsFloat320()
	require.NoError(t, err)
	assert.InDelta(t, float32(99.5), got2, 0.001)
}

func TestFilterPredicate_FilterValueVariant(t *testing.T) {
	var fv FilterValue
	require.NoError(t, fv.FromString1("match-me"))

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterValue(fv))

	// Round-trip through JSON.
	data, err := json.Marshal(fp)
	require.NoError(t, err)
	assert.JSONEq(t, `"match-me"`, string(data))

	var decoded FilterPredicate
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotFV, err := decoded.AsFilterValue()
	require.NoError(t, err)
	gotStr, err := gotFV.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "match-me", gotStr)
}

func TestFilterColumnIncludes_RoundTrip(t *testing.T) {
	// Build: FilterColumnIncludes -> FilterPredicate (FilterValue variant) + additional properties.
	var fv FilterValue
	require.NoError(t, fv.FromString1("match-me"))

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterValue(fv))

	original := FilterColumnIncludes{
		DollarSignIncludes: &fp,
		AdditionalProperties: map[string]any{
			"extra": "data",
		},
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded FilterColumnIncludes
	require.NoError(t, json.Unmarshal(data, &decoded))

	require.NotNil(t, decoded.DollarSignIncludes)
	gotFV, err := decoded.DollarSignIncludes.AsFilterValue()
	require.NoError(t, err)
	gotStr, err := gotFV.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "match-me", gotStr)

	assert.Equal(t, "data", decoded.AdditionalProperties["extra"])
}

func TestFilterColumnIncludes_AdditionalPropertiesOnly(t *testing.T) {
	fci := FilterColumnIncludes{
		AdditionalProperties: map[string]any{
			"customField": "customValue",
		},
	}

	data, err := json.Marshal(fci)
	require.NoError(t, err)

	var decoded FilterColumnIncludes
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, "customValue", decoded.AdditionalProperties["customField"])
}

func TestFilterPredicate_PredicateOpVariant(t *testing.T) {
	// FilterPredicate -> FilterPredicateOp with $any containing an array of FilterPredicates.
	var innerFV FilterValue
	require.NoError(t, innerFV.FromFloat320(42))

	var innerFP FilterPredicate
	require.NoError(t, innerFP.FromFilterValue(innerFV))

	var opAny FilterPredicateOpAny
	require.NoError(t, opAny.FromFilterPredicateOpAnyOneOf0(FilterPredicateOpAnyOneOf0{innerFP}))

	op := FilterPredicateOp{
		DollarSignAny: &opAny,
	}

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterPredicateOp(op))

	data, err := json.Marshal(fp)
	require.NoError(t, err)

	var decoded FilterPredicate
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotOp, err := decoded.AsFilterPredicateOp()
	require.NoError(t, err)
	require.NotNil(t, gotOp.DollarSignAny)

	gotSlice, err := gotOp.DollarSignAny.AsFilterPredicateOpAnyOneOf0()
	require.NoError(t, err)
	require.Len(t, gotSlice, 1)

	gotFV, err := gotSlice[0].AsFilterValue()
	require.NoError(t, err)
	gotFloat, err := gotFV.AsFloat320()
	require.NoError(t, err)
	assert.InDelta(t, float32(42), gotFloat, 0.001)
}

func TestFilterPredicate_ArrayVariant(t *testing.T) {
	// FilterPredicate -> FilterPredicateOneOf1 ([]FilterPredicate).
	var fv1 FilterValue
	require.NoError(t, fv1.FromString1("a"))
	var fv2 FilterValue
	require.NoError(t, fv2.FromString1("b"))

	var fp1 FilterPredicate
	require.NoError(t, fp1.FromFilterValue(fv1))
	var fp2 FilterPredicate
	require.NoError(t, fp2.FromFilterValue(fv2))

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterPredicateOneOf1(FilterPredicateOneOf1{fp1, fp2}))

	data, err := json.Marshal(fp)
	require.NoError(t, err)

	var decoded FilterPredicate
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotSlice, err := decoded.AsFilterPredicateOneOf1()
	require.NoError(t, err)
	require.Len(t, gotSlice, 2)

	s1, err := gotSlice[0].AsFilterValue()
	require.NoError(t, err)
	str1, err := s1.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "a", str1)

	s2, err := gotSlice[1].AsFilterValue()
	require.NoError(t, err)
	str2, err := s2.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "b", str2)
}

func TestFilterPredicate_RangeOpVariant(t *testing.T) {
	var rangeVal FilterRangeValue
	require.NoError(t, rangeVal.FromString1("100"))

	rangeOp := FilterPredicateRangeOp{
		DollarSignLt: &rangeVal,
	}

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterPredicateRangeOp(rangeOp))

	data, err := json.Marshal(fp)
	require.NoError(t, err)

	var decoded FilterPredicate
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotRangeOp, err := decoded.AsFilterPredicateRangeOp()
	require.NoError(t, err)
	require.NotNil(t, gotRangeOp.DollarSignLt)

	gotStr, err := gotRangeOp.DollarSignLt.AsString1()
	require.NoError(t, err)
	assert.Equal(t, "100", gotStr)
}

func TestFilterPredicateOp_AdditionalProperties(t *testing.T) {
	op := FilterPredicateOp{
		AdditionalProperties: map[string]any{
			"$custom": "value",
		},
	}

	data, err := json.Marshal(op)
	require.NoError(t, err)

	var decoded FilterPredicateOp
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, "value", decoded.AdditionalProperties["$custom"])
}

func TestFilterPredicateOpNone_PredicateVariant(t *testing.T) {
	var fv FilterValue
	require.NoError(t, fv.FromBool2(false))

	var fp FilterPredicate
	require.NoError(t, fp.FromFilterValue(fv))

	var none FilterPredicateOpNone
	require.NoError(t, none.FromFilterPredicate(fp))

	data, err := json.Marshal(none)
	require.NoError(t, err)

	var decoded FilterPredicateOpNone
	require.NoError(t, json.Unmarshal(data, &decoded))

	gotFP, err := decoded.AsFilterPredicate()
	require.NoError(t, err)
	gotFV, err := gotFP.AsFilterValue()
	require.NoError(t, err)
	gotBool, err := gotFV.AsBool2()
	require.NoError(t, err)
	assert.False(t, gotBool)
}

func TestApplyDefaults(t *testing.T) {
	// ApplyDefaults should be callable on all types without panic.
	fci := &FilterColumnIncludes{}
	fci.ApplyDefaults()

	fv := &FilterValue{}
	fv.ApplyDefaults()

	frv := &FilterRangeValue{}
	frv.ApplyDefaults()

	fp := &FilterPredicate{}
	fp.ApplyDefaults()

	fpo := &FilterPredicateOp{}
	fpo.ApplyDefaults()

	fpr := &FilterPredicateRangeOp{}
	fpr.ApplyDefaults()

	fpa := &FilterPredicateOpAny{}
	fpa.ApplyDefaults()

	fpn := &FilterPredicateOpNone{}
	fpn.ApplyDefaults()
}

func TestGetOpenAPISpecJSON(t *testing.T) {
	data, err := GetOpenAPISpecJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}
