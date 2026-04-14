package params

//oapi-runtime:function params/StyleMatrixParam

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

// StyleMatrixParam serializes a value using matrix style (RFC 6570).
// Matrix style prefixes values with ;paramName=. Path parameters only.
//
// Non-explode: Primitives: ;p=val   Arrays: ;p=a,b,c           Objects: ;p=k1,v1,k2,v2
// Explode:     Primitives: ;p=val   Arrays: ;p=a;p=b;p=c       Objects: ;k1=v1;k2=v2
func StyleMatrixParam(paramName string, value any, opts ParameterOptions) (string, error) {
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
			return fmt.Sprintf(";%s=%s", paramName, escapeParameterString(string(b), opts.ParamLocation)), nil
		}
	}

	switch t.Kind() {
	case reflect.Slice:
		n := v.Len()
		sliceVal := make([]any, n)
		for i := 0; i < n; i++ {
			sliceVal[i] = v.Index(i).Interface()
		}
		return styleMatrixSlice(paramName, opts, sliceVal)
	case reflect.Struct:
		return styleMatrixStruct(paramName, opts, value)
	case reflect.Map:
		return styleMatrixMap(paramName, opts, value)
	default:
		strVal, err := primitiveToString(value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf(";%s=%s", paramName, escapeParameterString(strVal, opts.ParamLocation)), nil
	}
}

func styleMatrixSlice(paramName string, opts ParameterOptions, values []any) (string, error) {
	prefix := fmt.Sprintf(";%s=", paramName)
	parts := make([]string, len(values))
	for i, v := range values {
		part, err := primitiveToString(v)
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		parts[i] = escapeParameterString(part, opts.ParamLocation)
	}
	if opts.Explode {
		// ;paramName=a;paramName=b;paramName=c
		return prefix + strings.Join(parts, prefix), nil
	}
	// ;paramName=a,b,c
	return prefix + strings.Join(parts, ","), nil
}

func styleMatrixStruct(paramName string, opts ParameterOptions, value any) (string, error) {
	if timeVal, ok := marshalKnownTypes(value); ok {
		return fmt.Sprintf(";%s=%s", paramName, escapeParameterString(timeVal, opts.ParamLocation)), nil
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
		return StyleMatrixParam(paramName, i2, opts)
	}

	fieldDict, err := structToFieldDict(value)
	if err != nil {
		return "", err
	}

	if opts.Explode {
		// ;key1=value1;key2=value2
		var parts []string
		for _, k := range sortedKeys(fieldDict) {
			v := escapeParameterString(fieldDict[k], opts.ParamLocation)
			parts = append(parts, k+"="+v)
		}
		return ";" + strings.Join(parts, ";"), nil
	}

	// ;paramName=key1,value1,key2,value2
	prefix := fmt.Sprintf(";%s=", paramName)
	var parts []string
	for _, k := range sortedKeys(fieldDict) {
		v := escapeParameterString(fieldDict[k], opts.ParamLocation)
		parts = append(parts, k, v)
	}
	return prefix + strings.Join(parts, ","), nil
}

func styleMatrixMap(paramName string, opts ParameterOptions, value any) (string, error) {
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

	if opts.Explode {
		var parts []string
		for _, k := range sortedKeys(fieldDict) {
			v := escapeParameterString(fieldDict[k], opts.ParamLocation)
			parts = append(parts, k+"="+v)
		}
		return ";" + strings.Join(parts, ";"), nil
	}

	prefix := fmt.Sprintf(";%s=", paramName)
	var parts []string
	for _, k := range sortedKeys(fieldDict) {
		v := escapeParameterString(fieldDict[k], opts.ParamLocation)
		parts = append(parts, k, v)
	}
	return prefix + strings.Join(parts, ","), nil
}
