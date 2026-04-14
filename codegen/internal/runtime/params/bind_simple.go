package params

//oapi-runtime:function params/BindSimpleParam

import (
	"encoding"
	"fmt"
	"reflect"
	"strings"
)

// BindSimpleParam binds a simple-style parameter to a destination.
// Simple style is the default for path and header parameters.
// Only used as a single-value entry point (no query variant needed).
//
// Non-explode: Arrays: a,b,c  Objects: key1,value1,key2,value2
// Explode:     Arrays: a,b,c  Objects: key1=value1,key2=value2
func BindSimpleParam(paramName string, value string, dest any, opts ParameterOptions) error {
	if value == "" {
		return fmt.Errorf("parameter '%s' is empty, can't bind its value", paramName)
	}

	var err error
	value, err = unescapeParameterString(value, opts.ParamLocation)
	if err != nil {
		return fmt.Errorf("error unescaping parameter '%s': %w", paramName, err)
	}

	if tu, ok := dest.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(value))
	}

	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	switch t.Kind() {
	case reflect.Struct:
		parts := strings.Split(value, ",")
		return bindSplitPartsToDestinationStruct(paramName, parts, opts.Explode, dest)
	case reflect.Slice:
		parts := strings.Split(value, ",")
		return bindSplitPartsToDestinationArray(parts, dest)
	default:
		return BindStringToObject(value, dest)
	}
}
