package params

//oapi-runtime:function params/ParameterOptions

// ParameterOptions carries OpenAPI parameter metadata to bind and style
// functions so they can handle explode, required, type-aware coercions,
// and location-aware escaping from a single uniform call site.
type ParameterOptions struct {
	ParamLocation ParamLocation
	Explode       bool
	Required      bool
	Type          string // OpenAPI type: "string", "integer", "array", "object"
	Format        string // OpenAPI format: "int32", "date-time", etc.
}
