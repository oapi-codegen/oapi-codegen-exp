package params

//oapi-runtime:function params/StyleSimpleParam

import (
	"bytes"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// StyleSimpleParam serializes a value using simple style (RFC 6570).
// Simple style is the default for path and header parameters.
//
// Non-explode: Arrays: a,b,c  Objects: key1,value1,key2,value2
// Explode:     Arrays: a,b,c  Objects: key1=value1,key2=value2
func StyleSimpleParam(paramName string, value any, opts ParameterOptions) (string, error) {
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", fmt.Errorf("value is a nil pointer")
		}
		v = reflect.Indirect(v)
		t = v.Type()
	}

	if tu, ok := value.(encoding.TextMarshaler); ok {
		innerT := reflect.Indirect(reflect.ValueOf(value)).Type()
		if !innerT.ConvertibleTo(reflect.TypeOf(time.Time{})) && !innerT.ConvertibleTo(reflect.TypeOf(types.Date{})) {
			b, err := tu.MarshalText()
			if err != nil {
				return "", fmt.Errorf("error marshaling '%s' as text: %w", value, err)
			}
			return escapeParameterString(string(b), opts.ParamLocation), nil
		}
	}

	switch t.Kind() {
	case reflect.Slice:
		n := v.Len()
		sliceVal := make([]any, n)
		for i := 0; i < n; i++ {
			sliceVal[i] = v.Index(i).Interface()
		}
		return styleSimpleSlice(paramName, opts, sliceVal)
	case reflect.Struct:
		return styleSimpleStruct(paramName, opts, value)
	case reflect.Map:
		return styleSimpleMap(paramName, opts, value)
	default:
		strVal, err := primitiveToString(value)
		if err != nil {
			return "", err
		}
		return escapeParameterString(strVal, opts.ParamLocation), nil
	}
}

func styleSimpleSlice(paramName string, opts ParameterOptions, values []any) (string, error) {
	// Simple arrays are always comma-separated regardless of explode
	parts := make([]string, len(values))
	for i, v := range values {
		part, err := primitiveToString(v)
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		parts[i] = escapeParameterString(part, opts.ParamLocation)
	}
	return strings.Join(parts, ","), nil
}

func styleSimpleStruct(paramName string, opts ParameterOptions, value any) (string, error) {
	if timeVal, ok := marshalKnownTypes(value); ok {
		return escapeParameterString(timeVal, opts.ParamLocation), nil
	}

	if m, ok := value.(json.Marshaler); ok {
		buf, err := m.MarshalJSON()
		if err != nil {
			return "", fmt.Errorf("failed to marshal to JSON: %w", err)
		}
		var i2 any
		e := json.NewDecoder(bytes.NewReader(buf))
		e.UseNumber()
		if err = e.Decode(&i2); err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return StyleSimpleParam(paramName, i2, opts)
	}

	fieldDict, err := structToFieldDict(value)
	if err != nil {
		return "", err
	}

	var parts []string
	for _, k := range sortedKeys(fieldDict) {
		v := escapeParameterString(fieldDict[k], opts.ParamLocation)
		if opts.Explode {
			parts = append(parts, k+"="+v)
		} else {
			parts = append(parts, k, v)
		}
	}
	return strings.Join(parts, ","), nil
}

func styleSimpleMap(paramName string, opts ParameterOptions, value any) (string, error) {
	dict, ok := value.(map[string]any)
	if !ok {
		return "", errors.New("map not of type map[string]any")
	}

	fieldDict := make(map[string]string)
	for fieldName, val := range dict {
		str, err := primitiveToString(val)
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		fieldDict[fieldName] = str
	}

	var parts []string
	for _, k := range sortedKeys(fieldDict) {
		v := escapeParameterString(fieldDict[k], opts.ParamLocation)
		if opts.Explode {
			parts = append(parts, k+"="+v)
		} else {
			parts = append(parts, k, v)
		}
	}
	return strings.Join(parts, ","), nil
}
