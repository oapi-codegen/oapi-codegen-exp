package params

//oapi-runtime:function params/StylePipeDelimitedParam

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// StylePipeDelimitedParam serializes a value using pipeDelimited style.
// Pipe-delimited style is used for query parameters with array values.
//
// Non-explode: paramName=a|b|c
// Explode:     paramName=a&paramName=b&paramName=c
func StylePipeDelimitedParam(paramName string, value any, opts ParameterOptions) (string, error) {
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
			return fmt.Sprintf("%s=%s", paramName, escapeParameterString(string(b), opts.ParamLocation)), nil
		}
	}

	switch t.Kind() {
	case reflect.Slice:
		n := v.Len()
		sliceVal := make([]any, n)
		for i := 0; i < n; i++ {
			sliceVal[i] = v.Index(i).Interface()
		}
		prefix := fmt.Sprintf("%s=", paramName)
		parts := make([]string, len(sliceVal))
		for i, sv := range sliceVal {
			part, err := primitiveToString(sv)
			if err != nil {
				return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
			}
			parts[i] = escapeParameterString(part, opts.ParamLocation)
		}
		if opts.Explode {
			return prefix + strings.Join(parts, "&"+prefix), nil
		}
		return prefix + strings.Join(parts, "|"), nil
	default:
		strVal, err := primitiveToString(value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s=%s", paramName, escapeParameterString(strVal, opts.ParamLocation)), nil
	}
}
