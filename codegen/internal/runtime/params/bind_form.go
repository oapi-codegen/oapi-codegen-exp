package params

//oapi-runtime:function params/BindFormParam

import (
	"encoding"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// BindFormParam binds a form-style parameter from a single string value.
// Used for path, header, and cookie parameters where the value has already
// been extracted from the request.
//
// Non-explode (default for form):
//
//	Arrays:  a,b,c -> []string{"a", "b", "c"}
//	Objects: key1,value1,key2,value2 -> struct{Key1, Key2}
//
// Explode:
//
//	Primitives and arrays: same comma-separated format
//	Objects: key1=value1,key2=value2 -> struct{Key1, Key2}
func BindFormParam(paramName string, value string, dest any, opts ParameterOptions) error {
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

// BindFormQueryParam binds a form-style query parameter from url.Values.
// The function looks up the parameter by name and handles both exploded
// and non-exploded formats.
//
// Non-explode: ?param=a,b,c (single query key, comma-separated value)
// Explode:     ?param=a&param=b&param=c (multiple query keys)
func BindFormQueryParam(paramName string, queryParams url.Values, dest any, opts ParameterOptions) error {
	if opts.Explode {
		return bindFormExplodeQuery(paramName, queryParams, dest, opts)
	}
	// Non-explode: single value, comma-separated
	value := queryParams.Get(paramName)
	if value == "" {
		if opts.Required {
			return fmt.Errorf("query parameter '%s' is required", paramName)
		}
		return nil
	}
	return BindFormParam(paramName, value, dest, opts)
}

// bindFormExplodeQuery handles the exploded form-style query parameter case.
func bindFormExplodeQuery(paramName string, queryParams url.Values, dest any, opts ParameterOptions) error {
	dv := reflect.Indirect(reflect.ValueOf(dest))
	v := dv
	var output any

	if opts.Required {
		output = dest
	} else {
		if v.IsNil() {
			t := v.Type()
			newValue := reflect.New(t.Elem())
			output = newValue.Interface()
		} else {
			output = v.Interface()
		}
		v = reflect.Indirect(reflect.ValueOf(output))
	}

	t := v.Type()
	k := t.Kind()

	values, found := queryParams[paramName]

	switch k {
	case reflect.Slice:
		if !found {
			if opts.Required {
				return fmt.Errorf("query parameter '%s' is required", paramName)
			}
			return nil
		}
		err := bindSplitPartsToDestinationArray(values, output)
		if err != nil {
			return err
		}
	case reflect.Struct:
		fieldsPresent, err := bindParamsToExplodedObject(paramName, queryParams, output)
		if err != nil {
			return err
		}
		if !fieldsPresent {
			return nil
		}
	default:
		if len(values) == 0 {
			if opts.Required {
				return fmt.Errorf("query parameter '%s' is required", paramName)
			}
			return nil
		}
		if len(values) != 1 {
			return fmt.Errorf("multiple values for single value parameter '%s'", paramName)
		}
		if !found {
			if opts.Required {
				return fmt.Errorf("query parameter '%s' is required", paramName)
			}
			return nil
		}
		err := BindStringToObject(values[0], output)
		if err != nil {
			return err
		}
	}

	if !opts.Required {
		dv.Set(reflect.ValueOf(output))
	}
	return nil
}

// bindParamsToExplodedObject binds query params to struct fields for exploded objects.
func bindParamsToExplodedObject(paramName string, values url.Values, dest any) (bool, error) {
	binder, v, t := indirectBinder(dest)
	if binder != nil {
		_, found := values[paramName]
		if !found {
			return false, nil
		}
		return true, BindStringToObject(values.Get(paramName), dest)
	}
	if t.Kind() != reflect.Struct {
		return false, fmt.Errorf("unmarshaling query arg '%s' into wrong type", paramName)
	}

	fieldsPresent := false
	for i := 0; i < t.NumField(); i++ {
		fieldT := t.Field(i)
		if !v.Field(i).CanSet() {
			continue
		}

		tag := fieldT.Tag.Get("json")
		fieldName := fieldT.Name
		if tag != "" {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] != "" {
				fieldName = tagParts[0]
			}
		}

		fieldVal, found := values[fieldName]
		if found {
			if len(fieldVal) != 1 {
				return false, fmt.Errorf("field '%s' specified multiple times for param '%s'", fieldName, paramName)
			}
			err := BindStringToObject(fieldVal[0], v.Field(i).Addr().Interface())
			if err != nil {
				return false, fmt.Errorf("could not bind query arg '%s': %w", paramName, err)
			}
			fieldsPresent = true
		}
	}
	return fieldsPresent, nil
}

// indirectBinder checks if dest implements Binder and returns reflect values.
func indirectBinder(dest any) (any, reflect.Value, reflect.Type) {
	v := reflect.ValueOf(dest)
	if v.Type().NumMethod() > 0 && v.CanInterface() {
		if u, ok := v.Interface().(Binder); ok {
			return u, reflect.Value{}, nil
		}
	}
	v = reflect.Indirect(v)
	t := v.Type()
	if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
		return dest, reflect.Value{}, nil
	}
	if t.ConvertibleTo(reflect.TypeOf(types.Date{})) {
		return dest, reflect.Value{}, nil
	}
	return nil, v, t
}
