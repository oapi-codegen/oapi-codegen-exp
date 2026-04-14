package params

//oapi-runtime:function params/ParameterOptions

// ParameterOptions carries OpenAPI parameter metadata to bind and style
// functions so they can handle style dispatch, explode, required,
// type-aware coercions, and location-aware escaping from a single
// uniform call site. All fields have sensible zero-value defaults.
type ParameterOptions struct {
	Style         string        // OpenAPI style: "simple", "form", "label", "matrix", "deepObject", "pipeDelimited", "spaceDelimited"
	ParamLocation ParamLocation // Where the parameter appears: query, path, header, cookie
	Explode       bool
	Required      bool
	Type          string // OpenAPI type: "string", "integer", "array", "object"
	Format        string // OpenAPI format: "int32", "date-time", etc.
	AllowReserved bool   // When true, reserved characters in query values are not percent-encoded
}
