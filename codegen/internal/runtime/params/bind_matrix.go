package params

//oapi-runtime:function params/BindMatrixParam

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// BindMatrixParam binds a matrix-style parameter to a destination.
// Matrix style values are prefixed with semicolons. Path parameters only.
//
// Non-explode:
//
//	Primitives: ;paramName=value
//	Arrays:     ;paramName=a,b,c
//	Objects:    ;paramName=key1,value1,key2,value2
//
// Explode:
//
//	Primitives: ;paramName=value
//	Arrays:     ;paramName=a;paramName=b;paramName=c
//	Objects:    ;key1=value1;key2=value2
func BindMatrixParam(paramName string, value string, dest any, opts ParameterOptions) error {
	if value == "" {
		return fmt.Errorf("parameter '%s' is empty, can't bind its value", paramName)
	}

	var err error
	value, err = unescapeParameterString(value, opts.ParamLocation)
	if err != nil {
		return fmt.Errorf("error unescaping parameter '%s': %w", paramName, err)
	}

	if opts.Explode {
		return bindMatrixExplode(paramName, value, dest)
	}
	return bindMatrixNoExplode(paramName, value, dest)
}

func bindMatrixNoExplode(paramName string, value string, dest any) error {
	prefix := ";" + paramName + "="
	if !strings.HasPrefix(value, prefix) {
		return fmt.Errorf("expected parameter '%s' to start with %s", paramName, prefix)
	}
	stripped := strings.TrimPrefix(value, prefix)

	if tu, ok := dest.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(stripped))
	}

	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	switch t.Kind() {
	case reflect.Struct:
		parts := strings.Split(stripped, ",")
		return bindSplitPartsToDestinationStruct(paramName, parts, false, dest)
	case reflect.Slice:
		parts := strings.Split(stripped, ",")
		return bindSplitPartsToDestinationArray(parts, dest)
	default:
		return BindStringToObject(stripped, dest)
	}
}

func bindMatrixExplode(paramName string, value string, dest any) error {
	parts := strings.Split(value, ";")
	if parts[0] != "" {
		return fmt.Errorf("invalid format for matrix parameter '%s', should start with ';'", paramName)
	}
	parts = parts[1:]

	if tu, ok := dest.(encoding.TextUnmarshaler); ok {
		if len(parts) == 1 {
			kv := strings.SplitN(parts[0], "=", 2)
			if len(kv) == 2 && kv[0] == paramName {
				return tu.UnmarshalText([]byte(kv[1]))
			}
		}
		return fmt.Errorf("invalid format for matrix parameter '%s'", paramName)
	}

	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	switch t.Kind() {
	case reflect.Struct:
		return bindSplitPartsToDestinationStruct(paramName, parts, true, dest)
	case reflect.Slice:
		prefix := paramName + "="
		values := make([]string, len(parts))
		for i, part := range parts {
			values[i] = strings.TrimPrefix(part, prefix)
		}
		return bindSplitPartsToDestinationArray(values, dest)
	default:
		if len(parts) == 1 {
			kv := strings.SplitN(parts[0], "=", 2)
			if len(kv) == 2 && kv[0] == paramName {
				return BindStringToObject(kv[1], dest)
			}
		}
		return fmt.Errorf("invalid format for matrix parameter '%s'", paramName)
	}
}
