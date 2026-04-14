package params

//oapi-runtime:function params/BindPipeDelimitedQueryParam

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// BindPipeDelimitedQueryParam binds a pipeDelimited-style query parameter.
// Pipe-delimited style uses pipes as array separators. Query only.
//
// Non-explode: ?param=a|b|c -> []string{"a", "b", "c"}
// Explode:     ?param=a&param=b -> []string{"a", "b"} (same as form explode)
func BindPipeDelimitedQueryParam(paramName string, queryParams url.Values, dest any, opts ParameterOptions) error {
	if opts.Explode {
		// Exploded pipe-delimited is same as exploded form
		return BindFormQueryParam(paramName, queryParams, dest, opts)
	}

	value := queryParams.Get(paramName)
	if value == "" {
		if opts.Required {
			return fmt.Errorf("query parameter '%s' is required", paramName)
		}
		return nil
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
	case reflect.Slice:
		parts := strings.Split(value, "|")
		return bindSplitPartsToDestinationArray(parts, dest)
	default:
		return BindStringToObject(value, dest)
	}
}
