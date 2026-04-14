package params

//oapi-runtime:function params/StyleParameter

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// StyleParameter serializes a Go value into an OpenAPI-styled parameter string.
// This is the entry point for client-side parameter serialization. The Style
// field in opts selects the serialization format. If Style is empty, "simple"
// is assumed.
func StyleParameter(paramName string, value any, opts ParameterOptions) (string, error) {
	style := opts.Style
	if style == "" {
		style = "simple"
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	// Dereference pointers; error on nil.
	if t.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", fmt.Errorf("value is a nil pointer")
		}
		v = reflect.Indirect(v)
		t = v.Type()
	}

	// If the value implements encoding.TextMarshaler, use it — but not for
	// time.Time or types.Date which have their own formatting logic.
	if tu, ok := value.(encoding.TextMarshaler); ok {
		it := reflect.Indirect(reflect.ValueOf(value)).Type()
		if !it.ConvertibleTo(reflect.TypeOf(time.Time{})) && !it.ConvertibleTo(reflect.TypeOf(types.Date{})) {
			b, err := tu.MarshalText()
			if err != nil {
				return "", fmt.Errorf("error marshaling '%s' as text: %w", value, err)
			}
			return stylePrimitive(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, string(b))
		}
	}

	switch t.Kind() {
	case reflect.Slice:
		if opts.Format == "byte" && isByteSlice(t) {
			encoded := base64.StdEncoding.EncodeToString(v.Bytes())
			return stylePrimitive(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, encoded)
		}
		n := v.Len()
		sliceVal := make([]any, n)
		for i := 0; i < n; i++ {
			sliceVal[i] = v.Index(i).Interface()
		}
		return styleSlice(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, sliceVal)
	case reflect.Struct:
		return styleStruct(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, value)
	case reflect.Map:
		return styleMap(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, value)
	default:
		return stylePrimitive(style, opts.Explode, paramName, opts.ParamLocation, opts.AllowReserved, value)
	}
}

// ---------------------------------------------------------------------------
// Internal style helpers
// ---------------------------------------------------------------------------

func styleSlice(style string, explode bool, paramName string, paramLocation ParamLocation, allowReserved bool, values []any) (string, error) {
	if style == "deepObject" {
		if !explode {
			return "", errors.New("deepObjects must be exploded")
		}
		return MarshalDeepObject(values, paramName)
	}

	var prefix string
	var separator string

	escapedName := escapeParameterName(paramName, paramLocation)

	switch style {
	case "simple":
		separator = ","
	case "label":
		prefix = "."
		if explode {
			separator = "."
		} else {
			separator = ","
		}
	case "matrix":
		prefix = fmt.Sprintf(";%s=", escapedName)
		if explode {
			separator = prefix
		} else {
			separator = ","
		}
	case "form":
		prefix = fmt.Sprintf("%s=", escapedName)
		if explode {
			separator = "&" + prefix
		} else {
			separator = ","
		}
	case "spaceDelimited":
		prefix = fmt.Sprintf("%s=", escapedName)
		if explode {
			separator = "&" + prefix
		} else {
			separator = " "
		}
	case "pipeDelimited":
		prefix = fmt.Sprintf("%s=", escapedName)
		if explode {
			separator = "&" + prefix
		} else {
			separator = "|"
		}
	default:
		return "", fmt.Errorf("unsupported style '%s'", style)
	}

	parts := make([]string, len(values))
	for i, v := range values {
		part, err := primitiveToString(v)
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		parts[i] = escapeParameterString(part, paramLocation, allowReserved)
	}
	return prefix + strings.Join(parts, separator), nil
}

func styleStruct(style string, explode bool, paramName string, paramLocation ParamLocation, allowReserved bool, value any) (string, error) {
	if timeVal, ok := marshalKnownTypes(value); ok {
		return stylePrimitive(style, explode, paramName, paramLocation, allowReserved, timeVal)
	}

	if style == "deepObject" {
		if !explode {
			return "", errors.New("deepObjects must be exploded")
		}
		return MarshalDeepObject(value, paramName)
	}

	// If input implements json.Marshaler (e.g. objects with additional properties
	// or anyOf), marshal to JSON and re-style the generic structure.
	if m, ok := value.(json.Marshaler); ok {
		buf, err := m.MarshalJSON()
		if err != nil {
			return "", fmt.Errorf("failed to marshal input to JSON: %w", err)
		}
		e := json.NewDecoder(bytes.NewReader(buf))
		e.UseNumber()
		var i2 any
		err = e.Decode(&i2)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
		return StyleParameter(paramName, i2, ParameterOptions{
			Style:         style,
			ParamLocation: paramLocation,
			Explode:       explode,
			AllowReserved: allowReserved,
		})
	}

	// Build a dictionary of the struct's fields.
	v := reflect.ValueOf(value)
	t := reflect.TypeOf(value)
	fieldDict := make(map[string]string)

	for i := 0; i < t.NumField(); i++ {
		fieldT := t.Field(i)
		tag := fieldT.Tag.Get("json")
		fieldName := fieldT.Name
		if tag != "" {
			tagParts := strings.Split(tag, ",")
			if tagParts[0] != "" {
				fieldName = tagParts[0]
			}
		}
		f := v.Field(i)

		// Skip nil optional fields.
		if f.Type().Kind() == reflect.Ptr && f.IsNil() {
			continue
		}
		str, err := primitiveToString(f.Interface())
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		fieldDict[fieldName] = str
	}

	return processFieldDict(style, explode, paramName, paramLocation, allowReserved, fieldDict)
}

func styleMap(style string, explode bool, paramName string, paramLocation ParamLocation, allowReserved bool, value any) (string, error) {
	if style == "deepObject" {
		if !explode {
			return "", errors.New("deepObjects must be exploded")
		}
		return MarshalDeepObject(value, paramName)
	}
	v := reflect.ValueOf(value)

	fieldDict := make(map[string]string)
	for _, fieldName := range v.MapKeys() {
		str, err := primitiveToString(v.MapIndex(fieldName).Interface())
		if err != nil {
			return "", fmt.Errorf("error formatting '%s': %w", paramName, err)
		}
		fieldDict[fieldName.String()] = str
	}
	return processFieldDict(style, explode, paramName, paramLocation, allowReserved, fieldDict)
}

func processFieldDict(style string, explode bool, paramName string, paramLocation ParamLocation, allowReserved bool, fieldDict map[string]string) (string, error) {
	var parts []string

	if style != "deepObject" {
		if explode {
			for _, k := range sortedKeys(fieldDict) {
				v := escapeParameterString(fieldDict[k], paramLocation, allowReserved)
				parts = append(parts, k+"="+v)
			}
		} else {
			for _, k := range sortedKeys(fieldDict) {
				v := escapeParameterString(fieldDict[k], paramLocation, allowReserved)
				parts = append(parts, k)
				parts = append(parts, v)
			}
		}
	}

	escapedName := escapeParameterName(paramName, paramLocation)

	var prefix string
	var separator string

	switch style {
	case "simple":
		separator = ","
	case "label":
		prefix = "."
		if explode {
			separator = prefix
		} else {
			separator = ","
		}
	case "matrix":
		if explode {
			separator = ";"
			prefix = ";"
		} else {
			separator = ","
			prefix = fmt.Sprintf(";%s=", escapedName)
		}
	case "form":
		if explode {
			separator = "&"
		} else {
			prefix = fmt.Sprintf("%s=", escapedName)
			separator = ","
		}
	case "deepObject":
		if !explode {
			return "", fmt.Errorf("deepObject parameters must be exploded")
		}
		for _, k := range sortedKeys(fieldDict) {
			v := fieldDict[k]
			part := fmt.Sprintf("%s[%s]=%s", escapedName, k, v)
			parts = append(parts, part)
		}
		separator = "&"
	default:
		return "", fmt.Errorf("unsupported style '%s'", style)
	}

	return prefix + strings.Join(parts, separator), nil
}

func stylePrimitive(style string, explode bool, paramName string, paramLocation ParamLocation, allowReserved bool, value any) (string, error) {
	strVal, err := primitiveToString(value)
	if err != nil {
		return "", err
	}

	escapedName := escapeParameterName(paramName, paramLocation)

	var prefix string
	switch style {
	case "simple":
	case "label":
		prefix = "."
	case "matrix":
		prefix = fmt.Sprintf(";%s=", escapedName)
	case "form":
		prefix = fmt.Sprintf("%s=", escapedName)
	default:
		return "", fmt.Errorf("unsupported style '%s'", style)
	}
	return prefix + escapeParameterString(strVal, paramLocation, allowReserved), nil
}

// ---------------------------------------------------------------------------
// Deep object marshaling
// ---------------------------------------------------------------------------

// MarshalDeepObject marshals an object to deepObject style query parameters.
func MarshalDeepObject(i any, paramName string) (string, error) {
	buf, err := json.Marshal(i)
	if err != nil {
		return "", fmt.Errorf("failed to marshal input to JSON: %w", err)
	}
	var i2 any
	err = json.Unmarshal(buf, &i2)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	fields, err := marshalDeepObjectRecursive(i2, nil)
	if err != nil {
		return "", fmt.Errorf("error traversing JSON structure: %w", err)
	}

	for idx := range fields {
		fields[idx] = paramName + fields[idx]
	}
	return strings.Join(fields, "&"), nil
}

func marshalDeepObjectRecursive(in any, path []string) ([]string, error) {
	var result []string

	switch t := in.(type) {
	case []any:
		for i, iface := range t {
			newPath := append(path, strconv.Itoa(i))
			fields, err := marshalDeepObjectRecursive(iface, newPath)
			if err != nil {
				return nil, fmt.Errorf("error traversing array: %w", err)
			}
			result = append(result, fields...)
		}
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			newPath := append(path, k)
			fields, err := marshalDeepObjectRecursive(t[k], newPath)
			if err != nil {
				return nil, fmt.Errorf("error traversing map: %w", err)
			}
			result = append(result, fields...)
		}
	default:
		prefix := "[" + strings.Join(path, "][") + "]"
		result = []string{
			prefix + fmt.Sprintf("=%v", t),
		}
	}
	return result, nil
}
