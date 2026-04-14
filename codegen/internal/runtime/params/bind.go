package params

//oapi-runtime:function params/BindParameter

import (
	"encoding"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// BindParameter binds a styled parameter from a single string value to a Go
// object. This is the entry point for path, header, and cookie parameters
// where the HTTP framework has already extracted the raw value.
//
// The Style field in opts selects how the value is split into parts (simple,
// label, matrix, form). If Style is empty, "simple" is assumed.
func BindParameter(paramName string, value string, dest any, opts ParameterOptions) error {
	style := opts.Style
	if style == "" {
		style = "simple"
	}

	if value == "" {
		if opts.Required {
			return &MissingRequiredParameterError{ParamName: paramName}
		}
		return nil
	}

	// Unescape based on parameter location.
	var err error
	value, err = unescapeParameterString(value, opts.ParamLocation)
	if err != nil {
		return fmt.Errorf("error unescaping parameter '%s': %w", paramName, err)
	}

	// If the destination implements encoding.TextUnmarshaler, use it directly.
	if tu, ok := dest.(encoding.TextUnmarshaler); ok {
		if err := tu.UnmarshalText([]byte(value)); err != nil {
			return fmt.Errorf("error unmarshaling '%s' text as %T: %w", value, dest, err)
		}
		return nil
	}

	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	if t.Kind() == reflect.Struct || t.Kind() == reflect.Map {
		parts, err := splitStyledParameter(style, opts.Explode, true, paramName, value)
		if err != nil {
			return err
		}
		return bindSplitPartsToDestinationStruct(paramName, parts, opts.Explode, dest)
	}

	if t.Kind() == reflect.Slice {
		if opts.Format == "byte" && isByteSlice(t) {
			parts, err := splitStyledParameter(style, opts.Explode, false, paramName, value)
			if err != nil {
				return fmt.Errorf("error splitting input '%s' into parts: %w", value, err)
			}
			if len(parts) != 1 {
				return fmt.Errorf("expected single base64 value for byte slice parameter '%s', got %d parts", paramName, len(parts))
			}
			decoded, err := base64Decode(parts[0])
			if err != nil {
				return fmt.Errorf("error decoding base64 parameter '%s': %w", paramName, err)
			}
			v.SetBytes(decoded)
			return nil
		}

		parts, err := splitStyledParameter(style, opts.Explode, false, paramName, value)
		if err != nil {
			return fmt.Errorf("error splitting input '%s' into parts: %w", value, err)
		}
		return bindSplitPartsToDestinationArray(parts, dest)
	}

	// Primitive types need style-specific prefix stripping before binding.
	// Label and matrix use splitStyledParameter for their prefix formats.
	// Form style adds a "name=" prefix (e.g. "p=5") which is meaningful in
	// query strings but must be stripped for cookie/header values. We use
	// TrimPrefix instead of splitStyledParameter to avoid splitting on commas,
	// which would break string primitives containing literal commas.
	switch style {
	case "label", "matrix":
		parts, err := splitStyledParameter(style, opts.Explode, false, paramName, value)
		if err != nil {
			return fmt.Errorf("error splitting parameter '%s': %w", paramName, err)
		}
		if len(parts) != 1 {
			return fmt.Errorf("parameter '%s': expected single value, got %d parts", paramName, len(parts))
		}
		value = parts[0]
	case "form":
		value = strings.TrimPrefix(value, paramName+"=")
	}
	return BindStringToObject(value, dest)
}

// BindQueryParameter binds a query parameter from pre-parsed url.Values.
// The Style field in opts selects parsing behavior. If Style is empty, "form"
// is assumed. Supports form, spaceDelimited, pipeDelimited, and deepObject.
func BindQueryParameter(paramName string, queryParams url.Values, dest any, opts ParameterOptions) error {
	style := opts.Style
	if style == "" {
		style = "form"
	}

	// Destination value management for optional (pointer) parameters.
	dv := reflect.Indirect(reflect.ValueOf(dest))
	v := dv
	var output any
	extraIndirect := !opts.Required && v.Kind() == reflect.Pointer
	if !extraIndirect {
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

	switch style {
	case "form", "spaceDelimited", "pipeDelimited":
		if opts.Explode {
			// Exploded: each value is a separate key=value pair.
			// spaceDelimited and pipeDelimited with explode=true are
			// serialized identically to form explode=true.
			values, found := queryParams[paramName]
			var err error

			switch k {
			case reflect.Slice:
				if !found {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				if opts.Format == "byte" && isByteSlice(t) {
					if len(values) != 1 {
						return fmt.Errorf("expected single base64 value for byte slice parameter '%s', got %d values", paramName, len(values))
					}
					decoded, decErr := base64Decode(values[0])
					if decErr != nil {
						return fmt.Errorf("error decoding base64 parameter '%s': %w", paramName, decErr)
					}
					v.SetBytes(decoded)
				} else {
					err = bindSplitPartsToDestinationArray(values, output)
				}
			case reflect.Struct:
				var fieldsPresent bool
				fieldsPresent, err = bindParamsToExplodedObject(paramName, queryParams, output)
				if !fieldsPresent {
					return nil
				}
			default:
				if len(values) == 0 {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				if len(values) != 1 {
					return fmt.Errorf("multiple values for single value parameter '%s'", paramName)
				}
				if !found {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				err = BindStringToObject(values[0], output)
			}
			if err != nil {
				return err
			}
			if extraIndirect {
				dv.Set(reflect.ValueOf(output))
			}
			return nil
		}

		// Non-exploded: single value, delimiter-separated.
		values, found := queryParams[paramName]
		if !found {
			if opts.Required {
				return &MissingRequiredParameterError{ParamName: paramName}
			}
			return nil
		}
		if len(values) != 1 {
			return fmt.Errorf("parameter '%s' is not exploded, but is specified multiple times", paramName)
		}

		// Primitive types: use the raw value as-is without splitting.
		if k != reflect.Slice && k != reflect.Struct && k != reflect.Map {
			err := BindStringToObject(values[0], output)
			if err != nil {
				return err
			}
			if extraIndirect {
				dv.Set(reflect.ValueOf(output))
			}
			return nil
		}

		var parts []string
		switch style {
		case "spaceDelimited":
			parts = strings.Split(values[0], " ")
		case "pipeDelimited":
			parts = strings.Split(values[0], "|")
		default:
			parts = strings.Split(values[0], ",")
		}

		var err error
		switch k {
		case reflect.Slice:
			if opts.Format == "byte" && isByteSlice(t) {
				raw := strings.Join(parts, ",")
				decoded, decErr := base64Decode(raw)
				if decErr != nil {
					return fmt.Errorf("error decoding base64 parameter '%s': %w", paramName, decErr)
				}
				v.SetBytes(decoded)
			} else {
				err = bindSplitPartsToDestinationArray(parts, output)
			}
		case reflect.Struct, reflect.Map:
			// Some struct types (e.g. types.Date, time.Time) are scalar values
			// that should be bound from a single string, not decomposed as
			// key-value objects.
			switch bv := output.(type) {
			case Binder:
				if len(parts) != 1 {
					return fmt.Errorf("multiple values for single value parameter '%s'", paramName)
				}
				err = bv.Bind(parts[0])
			case encoding.TextUnmarshaler:
				if len(parts) != 1 {
					return fmt.Errorf("multiple values for single value parameter '%s'", paramName)
				}
				err = bv.UnmarshalText([]byte(parts[0]))
			default:
				err = bindSplitPartsToDestinationStruct(paramName, parts, opts.Explode, output)
			}
		}
		if err != nil {
			return err
		}
		if extraIndirect {
			dv.Set(reflect.ValueOf(output))
		}
		return nil

	case "deepObject":
		if !opts.Explode {
			return errors.New("deepObjects must be exploded")
		}
		return unmarshalDeepObject(dest, paramName, queryParams, opts.Required)

	default:
		return fmt.Errorf("style '%s' on parameter '%s' is invalid", style, paramName)
	}
}

// BindRawQueryParameter works like BindQueryParameter but operates on the raw
// (undecoded) query string. This correctly handles form/explode=false
// parameters whose values contain literal commas encoded as %2C — something
// that BindQueryParameter cannot do because url.Values has already decoded
// %2C to ',' before we can split on the delimiter comma.
func BindRawQueryParameter(paramName string, rawQuery string, dest any, opts ParameterOptions) error {
	style := opts.Style
	if style == "" {
		style = "form"
	}

	dv := reflect.Indirect(reflect.ValueOf(dest))
	v := dv
	var output any
	extraIndirect := !opts.Required && v.Kind() == reflect.Pointer
	if !extraIndirect {
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

	switch style {
	case "form", "spaceDelimited", "pipeDelimited":
		if opts.Explode {
			// For explode, url.ParseQuery is fine — no delimiter commas to
			// confuse with literal commas.
			queryParams, err := url.ParseQuery(rawQuery)
			if err != nil {
				return fmt.Errorf("error parsing query string: %w", err)
			}
			values, found := queryParams[paramName]

			switch k {
			case reflect.Slice:
				if !found {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				err = bindSplitPartsToDestinationArray(values, output)
			case reflect.Struct:
				var fieldsPresent bool
				fieldsPresent, err = bindParamsToExplodedObject(paramName, queryParams, output)
				if !fieldsPresent {
					return nil
				}
			default:
				if len(values) == 0 {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				if len(values) != 1 {
					return fmt.Errorf("multiple values for single value parameter '%s'", paramName)
				}
				if !found {
					if opts.Required {
						return &MissingRequiredParameterError{ParamName: paramName}
					}
					return nil
				}
				err = BindStringToObject(values[0], output)
			}
			if err != nil {
				return err
			}
			if extraIndirect {
				dv.Set(reflect.ValueOf(output))
			}
			return nil
		}

		// explode=false — use findRawQueryParam to get the still-encoded
		// value, split on the style-specific delimiter, then URL-decode
		// each resulting part individually.
		rawValues, found := findRawQueryParam(rawQuery, paramName)
		if !found {
			if opts.Required {
				return &MissingRequiredParameterError{ParamName: paramName}
			}
			return nil
		}
		if len(rawValues) != 1 {
			return fmt.Errorf("parameter '%s' is not exploded, but is specified multiple times", paramName)
		}

		// Primitive types: decode as-is without splitting.
		if k != reflect.Slice && k != reflect.Struct && k != reflect.Map {
			decoded, err := url.QueryUnescape(rawValues[0])
			if err != nil {
				return fmt.Errorf("error decoding query parameter '%s' value %q: %w", paramName, rawValues[0], err)
			}
			err = BindStringToObject(decoded, output)
			if err != nil {
				return err
			}
			if extraIndirect {
				dv.Set(reflect.ValueOf(output))
			}
			return nil
		}

		var rawParts []string
		switch style {
		case "spaceDelimited":
			normalized := strings.ReplaceAll(rawValues[0], "+", "%20")
			normalized = strings.ReplaceAll(normalized, " ", "%20")
			rawParts = strings.Split(normalized, "%20")
		case "pipeDelimited":
			rawParts = strings.Split(rawValues[0], "|")
		default:
			rawParts = strings.Split(rawValues[0], ",")
		}

		parts := make([]string, len(rawParts))
		for i, rp := range rawParts {
			decoded, err := url.QueryUnescape(rp)
			if err != nil {
				return fmt.Errorf("error decoding query parameter '%s' part %q: %w", paramName, rp, err)
			}
			parts[i] = decoded
		}

		var err error
		switch k {
		case reflect.Slice:
			err = bindSplitPartsToDestinationArray(parts, output)
		case reflect.Struct:
			err = bindSplitPartsToDestinationStruct(paramName, parts, opts.Explode, output)
		}
		if err != nil {
			return err
		}
		if extraIndirect {
			dv.Set(reflect.ValueOf(output))
		}
		return nil

	case "deepObject":
		if !opts.Explode {
			return errors.New("deepObjects must be exploded")
		}
		queryParams, err := url.ParseQuery(rawQuery)
		if err != nil {
			return fmt.Errorf("error parsing query string: %w", err)
		}
		return unmarshalDeepObject(dest, paramName, queryParams, opts.Required)

	default:
		return fmt.Errorf("style '%s' on parameter '%s' is invalid", style, paramName)
	}
}

// ---------------------------------------------------------------------------
// Deep object internals
// ---------------------------------------------------------------------------

// unmarshalDeepObject is the internal implementation of deep object
// unmarshaling that supports the required parameter.
func unmarshalDeepObject(dst any, paramName string, params url.Values, required bool) error {
	var fieldNames []string
	var fieldValues []string
	searchStr := paramName + "["

	for pName, pValues := range params {
		if strings.HasPrefix(pName, searchStr) {
			pName = pName[len(paramName):]
			if len(pValues) == 1 {
				fieldNames = append(fieldNames, pName)
				fieldValues = append(fieldValues, pValues[0])
			} else {
				for i, value := range pValues {
					fieldNames = append(fieldNames, pName+"["+strconv.Itoa(i)+"]")
					fieldValues = append(fieldValues, value)
				}
			}
		}
	}

	if len(fieldNames) == 0 {
		if required {
			return &MissingRequiredParameterError{ParamName: paramName}
		}
		return nil
	}

	paths := make([][]string, len(fieldNames))
	for i, path := range fieldNames {
		path = strings.TrimLeft(path, "[")
		path = strings.TrimRight(path, "]")
		paths[i] = strings.Split(path, "][")
	}

	fieldPaths := makeFieldOrValue(paths, fieldValues)
	err := assignPathValues(dst, fieldPaths)
	if err != nil {
		return fmt.Errorf("error assigning value to destination: %w", err)
	}
	return nil
}

// UnmarshalDeepObject unmarshals deepObject-style query parameters to a
// destination. Exported for use by generated code and tests.
func UnmarshalDeepObject(dst any, paramName string, params url.Values) error {
	return unmarshalDeepObject(dst, paramName, params, false)
}

type fieldOrValue struct {
	fields map[string]fieldOrValue
	value  string
}

func (f *fieldOrValue) appendPathValue(path []string, value string) {
	fieldName := path[0]
	if len(path) == 1 {
		f.fields[fieldName] = fieldOrValue{value: value}
		return
	}

	pv, found := f.fields[fieldName]
	if !found {
		pv = fieldOrValue{
			fields: make(map[string]fieldOrValue),
		}
		f.fields[fieldName] = pv
	}
	pv.appendPathValue(path[1:], value)
}

func makeFieldOrValue(paths [][]string, values []string) fieldOrValue {
	f := fieldOrValue{
		fields: make(map[string]fieldOrValue),
	}
	for i := range paths {
		f.appendPathValue(paths[i], values[i])
	}
	return f
}

func getFieldName(f reflect.StructField) string {
	n := f.Name
	tag, found := f.Tag.Lookup("json")
	if found {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			n = parts[0]
		}
	}
	return n
}

func fieldIndicesByJsonTag(i any) (map[string]int, error) {
	t := reflect.TypeOf(i)
	if t.Kind() != reflect.Struct {
		return nil, errors.New("expected a struct as input")
	}

	n := t.NumField()
	fieldMap := make(map[string]int)
	for i := 0; i < n; i++ {
		field := t.Field(i)
		fieldName := getFieldName(field)
		fieldMap[fieldName] = i
	}
	return fieldMap, nil
}

func assignPathValues(dst any, pathValues fieldOrValue) error {
	v := reflect.ValueOf(dst)
	iv := reflect.Indirect(v)
	it := iv.Type()

	switch it.Kind() {
	case reflect.Map:
		dstMap := reflect.MakeMap(iv.Type())
		for key, value := range pathValues.fields {
			dstKey := reflect.ValueOf(key)
			dstVal := reflect.New(iv.Type().Elem())
			err := assignPathValues(dstVal.Interface(), value)
			if err != nil {
				return fmt.Errorf("error binding map: %w", err)
			}
			dstMap.SetMapIndex(dstKey, dstVal.Elem())
		}
		iv.Set(dstMap)
		return nil

	case reflect.Slice:
		sliceLength := len(pathValues.fields)
		dstSlice := reflect.MakeSlice(it, sliceLength, sliceLength)
		err := assignDeepObjectSlice(dstSlice, pathValues)
		if err != nil {
			return fmt.Errorf("error assigning slice: %w", err)
		}
		iv.Set(dstSlice)
		return nil

	case reflect.Struct:
		if dst, isBinder := v.Interface().(Binder); isBinder {
			return dst.Bind(pathValues.value)
		}

		if it.ConvertibleTo(reflect.TypeOf(types.Date{})) {
			var date types.Date
			var err error
			date.Time, err = time.Parse(types.DateFormat, pathValues.value)
			if err != nil {
				return fmt.Errorf("invalid date format: %w", err)
			}
			dst := iv
			if it != reflect.TypeOf(types.Date{}) {
				ivPtr := iv.Addr()
				aPtr := ivPtr.Convert(reflect.TypeOf(&types.Date{}))
				dst = reflect.Indirect(aPtr)
			}
			dst.Set(reflect.ValueOf(date))
			return nil
		}

		if it.ConvertibleTo(reflect.TypeOf(time.Time{})) {
			tm, err := time.Parse(time.RFC3339Nano, pathValues.value)
			if err != nil {
				tm, err = time.Parse(types.DateFormat, pathValues.value)
				if err != nil {
					return fmt.Errorf("error parsing '%s' as RFC3339 or date: %w", pathValues.value, err)
				}
			}
			dst := iv
			if it != reflect.TypeOf(time.Time{}) {
				ivPtr := iv.Addr()
				aPtr := ivPtr.Convert(reflect.TypeOf(&time.Time{}))
				dst = reflect.Indirect(aPtr)
			}
			dst.Set(reflect.ValueOf(tm))
			return nil
		}

		fieldMap, err := fieldIndicesByJsonTag(iv.Interface())
		if err != nil {
			return fmt.Errorf("failed enumerating fields: %w", err)
		}
		for _, fieldName := range sortedFieldOrValueKeys(pathValues.fields) {
			fieldValue := pathValues.fields[fieldName]
			fieldIndex, found := fieldMap[fieldName]
			if !found {
				return fmt.Errorf("field [%s] is not present in destination object", fieldName)
			}
			field := iv.Field(fieldIndex)
			err = assignPathValues(field.Addr().Interface(), fieldValue)
			if err != nil {
				return fmt.Errorf("error assigning field [%s]: %w", fieldName, err)
			}
		}
		return nil

	case reflect.Ptr:
		dstVal := reflect.New(it.Elem())
		dstPtr := dstVal.Interface()
		err := assignPathValues(dstPtr, pathValues)
		iv.Set(dstVal)
		return err

	case reflect.Bool:
		val, err := strconv.ParseBool(pathValues.value)
		if err != nil {
			return fmt.Errorf("expected a valid bool, got %s", pathValues.value)
		}
		iv.SetBool(val)
		return nil

	case reflect.Float32:
		val, err := strconv.ParseFloat(pathValues.value, 32)
		if err != nil {
			return fmt.Errorf("expected a valid float, got %s", pathValues.value)
		}
		iv.SetFloat(val)
		return nil

	case reflect.Float64:
		val, err := strconv.ParseFloat(pathValues.value, 64)
		if err != nil {
			return fmt.Errorf("expected a valid float, got %s", pathValues.value)
		}
		iv.SetFloat(val)
		return nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(pathValues.value, 10, 64)
		if err != nil {
			return fmt.Errorf("expected a valid int, got %s", pathValues.value)
		}
		iv.SetInt(val)
		return nil

	case reflect.String:
		iv.SetString(pathValues.value)
		return nil

	default:
		return errors.New("unhandled type: " + it.String())
	}
}

func assignDeepObjectSlice(dst reflect.Value, pathValues fieldOrValue) error {
	nValues := len(pathValues.fields)
	values := make([]string, nValues)
	for i := 0; i < nValues; i++ {
		indexStr := strconv.Itoa(i)
		fv, found := pathValues.fields[indexStr]
		if !found {
			return errors.New("array deepObjects must have consecutive indices")
		}
		values[i] = fv.value
	}

	for i := 0; i < nValues; i++ {
		dstElem := dst.Index(i).Addr()
		err := assignPathValues(dstElem.Interface(), fieldOrValue{value: values[i]})
		if err != nil {
			return fmt.Errorf("error binding array: %w", err)
		}
	}
	return nil
}

func sortedFieldOrValueKeys(m map[string]fieldOrValue) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ---------------------------------------------------------------------------
// Exploded object binding
// ---------------------------------------------------------------------------

// bindParamsToExplodedObject reflects the destination structure and pulls the
// value for each settable field from the given query parameters. Returns
// whether any fields were bound.
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
