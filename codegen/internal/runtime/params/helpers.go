package params

//oapi-runtime:function params/ParamHelpers

import (
	"bytes"
	"encoding"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oapi-codegen/oapi-codegen-exp/codegen/internal/runtime/types"
)

// ParamLocation indicates where a parameter is located in an HTTP request.
type ParamLocation int

const (
	ParamLocationUndefined ParamLocation = iota
	ParamLocationQuery
	ParamLocationPath
	ParamLocationHeader
	ParamLocationCookie
)

// Binder is an interface for types that can bind themselves from a string value.
type Binder interface {
	Bind(value string) error
}

// MissingRequiredParameterError is returned when a required parameter is not
// present in the request. Upper layers can use errors.As to detect this and
// produce an appropriate HTTP error response.
type MissingRequiredParameterError struct {
	ParamName string
}

func (e *MissingRequiredParameterError) Error() string {
	return fmt.Sprintf("parameter '%s' is required", e.ParamName)
}

// primitiveToString converts a primitive value to a string representation.
// It handles basic Go types, time.Time, types.Date, and types that implement
// json.Marshaler or fmt.Stringer.
func primitiveToString(value any) (string, error) {
	// Check for known types first (time, date, uuid)
	if res, ok := marshalKnownTypes(value); ok {
		return res, nil
	}

	// Dereference pointers for optional values
	v := reflect.Indirect(reflect.ValueOf(value))
	t := v.Type()
	kind := t.Kind()

	switch kind {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), nil
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), nil
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		}
		return "false", nil
	case reflect.String:
		return v.String(), nil
	case reflect.Struct:
		// Check if it's a UUID
		if u, ok := value.(uuid.UUID); ok {
			return u.String(), nil
		}
		// Check if it implements json.Marshaler
		if m, ok := value.(json.Marshaler); ok {
			buf, err := m.MarshalJSON()
			if err != nil {
				return "", fmt.Errorf("failed to marshal to JSON: %w", err)
			}
			e := json.NewDecoder(bytes.NewReader(buf))
			e.UseNumber()
			var i2 any
			if err = e.Decode(&i2); err != nil {
				return "", fmt.Errorf("failed to decode JSON: %w", err)
			}
			return primitiveToString(i2)
		}
		fallthrough
	default:
		if s, ok := value.(fmt.Stringer); ok {
			return s.String(), nil
		}
		return "", fmt.Errorf("unsupported type %s", reflect.TypeOf(value).String())
	}
}

// marshalKnownTypes checks for special types (time.Time, Date, UUID) and marshals them.
func marshalKnownTypes(value any) (string, bool) {
	v := reflect.Indirect(reflect.ValueOf(value))
	t := v.Type()

	if t.ConvertibleTo(reflect.TypeOf(time.Time{})) {
		tt := v.Convert(reflect.TypeOf(time.Time{}))
		timeVal := tt.Interface().(time.Time)
		return timeVal.Format(time.RFC3339Nano), true
	}

	if t.ConvertibleTo(reflect.TypeOf(types.Date{})) {
		d := v.Convert(reflect.TypeOf(types.Date{}))
		dateVal := d.Interface().(types.Date)
		return dateVal.Format(types.DateFormat), true
	}

	if t.ConvertibleTo(reflect.TypeOf(uuid.UUID{})) {
		u := v.Convert(reflect.TypeOf(uuid.UUID{}))
		uuidVal := u.Interface().(uuid.UUID)
		return uuidVal.String(), true
	}

	return "", false
}

// escapeParameterName escapes a parameter name for use in query strings and
// paths. This ensures characters like [] in parameter names (e.g. user_ids[])
// are properly percent-encoded per RFC 3986.
func escapeParameterName(name string, paramLocation ParamLocation) string {
	// Parameter names should always be encoded regardless of allowReserved,
	// which only applies to values per the OpenAPI spec.
	return escapeParameterString(name, paramLocation, false)
}

// escapeParameterString escapes a parameter value based on its location.
// Query and path parameters need URL escaping; headers and cookies do not.
// When allowReserved is true and the location is query, RFC 3986 reserved
// characters are left unencoded per the OpenAPI allowReserved specification.
func escapeParameterString(value string, paramLocation ParamLocation, allowReserved bool) string {
	switch paramLocation {
	case ParamLocationQuery:
		if allowReserved {
			return escapeQueryAllowReserved(value)
		}
		return url.QueryEscape(value)
	case ParamLocationPath:
		return url.PathEscape(value)
	default:
		return value
	}
}

// escapeQueryAllowReserved percent-encodes a query parameter value while
// leaving RFC 3986 reserved characters (:/?#[]@!$&'()*+,;=) unencoded, as
// specified by OpenAPI's allowReserved parameter option.
func escapeQueryAllowReserved(value string) string {
	const reserved = `:/?#[]@!$&'()*+,;=`

	var buf strings.Builder
	for _, b := range []byte(value) {
		if isUnreserved(b) || strings.IndexByte(reserved, b) >= 0 {
			buf.WriteByte(b)
		} else {
			fmt.Fprintf(&buf, "%%%02X", b)
		}
	}
	return buf.String()
}

// isUnreserved reports whether the byte is an RFC 3986 unreserved character:
// ALPHA / DIGIT / "-" / "." / "_" / "~"
func isUnreserved(c byte) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '.' || c == '_' || c == '~'
}

// unescapeParameterString unescapes a parameter value based on its location.
func unescapeParameterString(value string, paramLocation ParamLocation) (string, error) {
	switch paramLocation {
	case ParamLocationQuery, ParamLocationUndefined:
		return url.QueryUnescape(value)
	case ParamLocationPath:
		return url.PathUnescape(value)
	default:
		return value, nil
	}
}

// sortedKeys returns the keys of a map in sorted order.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// BindStringToObject binds a string value to a destination object.
// It handles primitives, encoding.TextUnmarshaler, and the Binder interface.
func BindStringToObject(src string, dst any) error {
	// Check for TextUnmarshaler
	if tu, ok := dst.(encoding.TextUnmarshaler); ok {
		return tu.UnmarshalText([]byte(src))
	}

	// Check for Binder interface
	if b, ok := dst.(Binder); ok {
		return b.Bind(src)
	}

	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("dst must be a pointer, got %T", dst)
	}
	v = v.Elem()

	switch v.Kind() {
	case reflect.String:
		v.SetString(src)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(src, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse int: %w", err)
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u, err := strconv.ParseUint(src, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse uint: %w", err)
		}
		v.SetUint(u)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(src, 64)
		if err != nil {
			return fmt.Errorf("failed to parse float: %w", err)
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(src)
		if err != nil {
			return fmt.Errorf("failed to parse bool: %w", err)
		}
		v.SetBool(b)
	default:
		// Try JSON unmarshal as a fallback
		return json.Unmarshal([]byte(src), dst)
	}
	return nil
}

// bindSplitPartsToDestinationArray binds a slice of string parts to a destination slice.
func bindSplitPartsToDestinationArray(parts []string, dest any) error {
	v := reflect.Indirect(reflect.ValueOf(dest))
	t := v.Type()

	newArray := reflect.MakeSlice(t, len(parts), len(parts))
	for i, p := range parts {
		err := BindStringToObject(p, newArray.Index(i).Addr().Interface())
		if err != nil {
			return fmt.Errorf("error setting array element: %w", err)
		}
	}
	v.Set(newArray)
	return nil
}

// bindSplitPartsToDestinationStruct binds string parts to a destination struct via JSON.
func bindSplitPartsToDestinationStruct(paramName string, parts []string, explode bool, dest any) error {
	var fields []string
	if explode {
		fields = make([]string, len(parts))
		for i, property := range parts {
			propertyParts := strings.Split(property, "=")
			if len(propertyParts) != 2 {
				return fmt.Errorf("parameter '%s' has invalid exploded format", paramName)
			}
			fields[i] = "\"" + propertyParts[0] + "\":\"" + propertyParts[1] + "\""
		}
	} else {
		if len(parts)%2 != 0 {
			return fmt.Errorf("parameter '%s' has invalid format, property/values need to be pairs", paramName)
		}
		fields = make([]string, len(parts)/2)
		for i := 0; i < len(parts); i += 2 {
			key := parts[i]
			value := parts[i+1]
			fields[i/2] = "\"" + key + "\":\"" + value + "\""
		}
	}
	jsonParam := "{" + strings.Join(fields, ",") + "}"
	return json.Unmarshal([]byte(jsonParam), dest)
}

// splitStyledParameter splits a styled parameter string value into parts based
// on the OpenAPI style. The object flag indicates whether the destination is a
// struct/map (affects matrix explode handling).
func splitStyledParameter(style string, explode bool, object bool, paramName string, value string) ([]string, error) {
	switch style {
	case "simple":
		// In the simple case, we always split on comma
		return strings.Split(value, ","), nil
	case "label":
		if explode {
			// Exploded: .a.b.c or .key=value.key=value
			parts := strings.Split(value, ".")
			if parts[0] != "" {
				return nil, fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
			}
			return parts[1:], nil
		}
		// Unexploded: .a,b,c
		if value[0] != '.' {
			return nil, fmt.Errorf("invalid format for label parameter '%s', should start with '.'", paramName)
		}
		return strings.Split(value[1:], ","), nil
	case "matrix":
		if explode {
			// Exploded: ;a;b;c or ;key=value;key=value
			parts := strings.Split(value, ";")
			if parts[0] != "" {
				return nil, fmt.Errorf("invalid format for matrix parameter '%s', should start with ';'", paramName)
			}
			parts = parts[1:]
			if !object {
				prefix := paramName + "="
				for i := range parts {
					parts[i] = strings.TrimPrefix(parts[i], prefix)
				}
			}
			return parts, nil
		}
		// Unexploded: ;paramName=a,b,c
		prefix := ";" + paramName + "="
		if !strings.HasPrefix(value, prefix) {
			return nil, fmt.Errorf("expected parameter '%s' to start with %s", paramName, prefix)
		}
		return strings.Split(strings.TrimPrefix(value, prefix), ","), nil
	case "form":
		if explode {
			parts := strings.Split(value, "&")
			if !object {
				prefix := paramName + "="
				for i := range parts {
					parts[i] = strings.TrimPrefix(parts[i], prefix)
				}
			}
			return parts, nil
		}
		parts := strings.Split(value, ",")
		prefix := paramName + "="
		for i := range parts {
			parts[i] = strings.TrimPrefix(parts[i], prefix)
		}
		return parts, nil
	}

	return nil, fmt.Errorf("unhandled parameter style: %s", style)
}

// findRawQueryParam extracts values for a named parameter from a raw
// (undecoded) query string. The parameter key is decoded for comparison
// purposes, but the returned values remain in their original encoded form.
func findRawQueryParam(rawQuery, paramName string) (values []string, found bool) {
	for rawQuery != "" {
		var part string
		if i := strings.IndexByte(rawQuery, '&'); i >= 0 {
			part = rawQuery[:i]
			rawQuery = rawQuery[i+1:]
		} else {
			part = rawQuery
			rawQuery = ""
		}
		if part == "" {
			continue
		}
		key := part
		var val string
		if i := strings.IndexByte(part, '='); i >= 0 {
			key = part[:i]
			val = part[i+1:]
		}
		decodedKey, err := url.QueryUnescape(key)
		if err != nil {
			// Skip malformed keys.
			continue
		}
		if decodedKey == paramName {
			values = append(values, val)
			found = true
		}
	}
	return values, found
}

// isByteSlice reports whether t is []byte (or equivalently []uint8).
func isByteSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8
}

// base64Decode decodes s as base64.
//
// Per OpenAPI 3.0, format: byte uses RFC 4648 Section 4 (standard alphabet,
// padded). We use padding presence to select the right decoder, rather than
// blindly cascading (which can produce corrupt output when RawStdEncoding
// silently accepts padded input and treats '=' as data).
func base64Decode(s string) ([]byte, error) {
	if s == "" {
		return []byte{}, nil
	}

	if strings.ContainsRune(s, '=') {
		if strings.ContainsAny(s, "-_") {
			return base64Decode1(base64.URLEncoding, s)
		}
		return base64Decode1(base64.StdEncoding, s)
	}

	if strings.ContainsAny(s, "-_") {
		return base64Decode1(base64.RawURLEncoding, s)
	}
	return base64Decode1(base64.RawStdEncoding, s)
}

func base64Decode1(enc *base64.Encoding, s string) ([]byte, error) {
	b, err := enc.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode string %q: %w", s, err)
	}
	return b, nil
}

// structToFieldDict converts a struct to a map of field names to string values.
func structToFieldDict(value any) (map[string]string, error) {
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

		// Skip nil optional fields
		if f.Type().Kind() == reflect.Ptr && f.IsNil() {
			continue
		}
		str, err := primitiveToString(f.Interface())
		if err != nil {
			return nil, fmt.Errorf("error formatting field '%s': %w", fieldName, err)
		}
		fieldDict[fieldName] = str
	}
	return fieldDict, nil
}
