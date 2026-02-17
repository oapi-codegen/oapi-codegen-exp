package codegen

import (
	"bytes"
	"sort"
	"strings"
	"text/template"
)

// StructTagInfo contains the data available to struct tag templates.
// Extension-driven concerns (omitzero, json-ignore, omitempty overrides)
// are applied as post-processing in generateFieldTag, not exposed here.
type StructTagInfo struct {
	// FieldName is the JSON/YAML field name (from the OpenAPI property name)
	FieldName string
	// IsOptional is true if the field is optional (not required)
	IsOptional bool
}

// StructTagTemplate defines a single struct tag with a name and template.
type StructTagTemplate struct {
	// Name is the tag name (e.g., "json", "yaml", "form")
	Name string `yaml:"name"`
	// Template is a Go text/template that produces the tag value.
	// Available fields: .FieldName, .IsOptional
	// Example: `{{ .FieldName }}{{if .IsOptional}},omitempty{{end}}`
	Template string `yaml:"template"`
}

// StructTagsConfig configures struct tag generation.
type StructTagsConfig struct {
	// Tags is the list of tags to generate for struct fields.
	// Order is preserved in the generated output.
	Tags []StructTagTemplate `yaml:"tags,omitempty"`
}

// DefaultStructTagsConfig returns the default struct tag configuration.
// By default, json and form tags are generated. Extension-driven concerns
// (omitzero, json-ignore, omitempty overrides) are handled by post-processing
// in generateFieldTag, not in the templates.
func DefaultStructTagsConfig() StructTagsConfig {
	return StructTagsConfig{
		Tags: []StructTagTemplate{
			{
				Name:     "json",
				Template: `{{ .FieldName }}{{if .IsOptional}},omitempty{{end}}`,
			},
			{
				Name:     "form",
				Template: `{{ .FieldName }}{{if .IsOptional}},omitempty{{end}}`,
			},
		},
	}
}

// Merge merges user config on top of this config by name.
// User entries override matching defaults; new entries are appended.
func (c StructTagsConfig) Merge(other StructTagsConfig) StructTagsConfig {
	if len(other.Tags) == 0 {
		return c
	}
	// Start with defaults, override/append from user config
	merged := make(map[string]StructTagTemplate, len(c.Tags)+len(other.Tags))
	order := make([]string, 0, len(c.Tags)+len(other.Tags))
	for _, t := range c.Tags {
		merged[t.Name] = t
		order = append(order, t.Name)
	}
	for _, t := range other.Tags {
		if _, exists := merged[t.Name]; !exists {
			order = append(order, t.Name)
		}
		merged[t.Name] = t
	}
	result := StructTagsConfig{Tags: make([]StructTagTemplate, 0, len(order))}
	for _, name := range order {
		result.Tags = append(result.Tags, merged[name])
	}
	return result
}

// StructTagGenerator generates struct tags from templates.
type StructTagGenerator struct {
	templates []*tagTemplate
}

type tagTemplate struct {
	name string
	tmpl *template.Template
}

// NewStructTagGenerator creates a generator from the configuration.
// Invalid templates are silently skipped.
func NewStructTagGenerator(config StructTagsConfig) *StructTagGenerator {
	g := &StructTagGenerator{
		templates: make([]*tagTemplate, 0, len(config.Tags)),
	}

	for _, tag := range config.Tags {
		tmpl, err := template.New(tag.Name).Parse(tag.Template)
		if err != nil {
			// Skip invalid templates
			continue
		}
		g.templates = append(g.templates, &tagTemplate{
			name: tag.Name,
			tmpl: tmpl,
		})
	}

	return g
}

// GenerateTags generates the complete struct tag string for a field.
// Returns a string like `json:"name,omitempty" yaml:"name,omitempty"`.
func (g *StructTagGenerator) GenerateTags(info StructTagInfo) string {
	if len(g.templates) == 0 {
		return ""
	}

	var tags []string
	for _, t := range g.templates {
		var buf bytes.Buffer
		if err := t.tmpl.Execute(&buf, info); err != nil {
			// Skip tags that fail to render
			continue
		}
		value := buf.String()
		if value != "" {
			tags = append(tags, t.name+`:`+`"`+value+`"`)
		}
	}

	if len(tags) == 0 {
		return ""
	}

	return "`" + strings.Join(tags, " ") + "`"
}

// GenerateTagsMap generates tags as a map for cases where we need to add extra tags.
// Returns a map of tag name -> tag value (without quotes).
func (g *StructTagGenerator) GenerateTagsMap(info StructTagInfo) map[string]string {
	result := make(map[string]string)

	for _, t := range g.templates {
		var buf bytes.Buffer
		if err := t.tmpl.Execute(&buf, info); err != nil {
			continue
		}
		value := buf.String()
		if value != "" {
			result[t.name] = value
		}
	}

	return result
}

// FormatTagsMap formats a tag map into a struct tag string.
// Tags are sorted alphabetically by name for deterministic output.
func FormatTagsMap(tags map[string]string) string {
	if len(tags) == 0 {
		return ""
	}

	// Sort tag names for deterministic output
	names := make([]string, 0, len(tags))
	for name := range tags {
		names = append(names, name)
	}
	sort.Strings(names)

	var parts []string
	for _, name := range names {
		parts = append(parts, name+`:`+`"`+tags[name]+`"`)
	}

	return "`" + strings.Join(parts, " ") + "`"
}

// applyOmitEmptyOverride rebuilds tag values when extensions override the default omitempty.
func applyOmitEmptyOverride(tags map[string]string, fieldName string, omitEmpty bool, tagNames ...string) {
	for _, name := range tagNames {
		if _, ok := tags[name]; !ok {
			continue
		}
		v := fieldName
		if omitEmpty {
			v += ",omitempty"
		}
		tags[name] = v
	}
}
