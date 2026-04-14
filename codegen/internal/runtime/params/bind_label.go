package params

//oapi-runtime:function params/BindLabelParam

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// BindLabelParam binds a label-style parameter to a destination.
// Label style values are prefixed with a dot. Path parameters only.
//
// Non-explode:
//
//	Primitives: .value          Arrays: .a,b,c          Objects: .key1,value1,key2,value2
//
// Explode:
//
//	Primitives: .value          Arrays: .a.b.c          Objects: .key1=value1.key2=value2
func BindLabelParam(paramName string, value string, dest any, opts ParameterOptions) error {
	if value == "" {
		return fmt.Errorf("parameter '%s' is empty, can't bind its value", paramName)
	}

	var err error
	value, err = unescapeParameterString(value, opts.ParamLocation)
	if err != nil {
		return fmt.Errorf("error unescaping parameter '%s': %w", paramName, err)
	}

	if value[0] != '.' {
		return fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
	}

	if tu, ok := dest.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(value[1:]))
	}

	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	if opts.Explode {
		// Explode: split on dot, skip first empty part
		parts := strings.Split(value, ".")
		if parts[0] != "" {
			return fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
		}
		parts = parts[1:]

		switch t.Kind() {
		case reflect.Struct:
			return bindSplitPartsToDestinationStruct(paramName, parts, true, dest)
		case reflect.Slice:
			return bindSplitPartsToDestinationArray(parts, dest)
		default:
			return BindStringToObject(value[1:], dest)
		}
	}

	// Non-explode: strip leading dot, split on comma
	stripped := value[1:]

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
