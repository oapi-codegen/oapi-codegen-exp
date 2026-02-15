package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pb33f/libopenapi/datamodel/high/base"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// templateEntry describes a single template to load.
type templateEntry struct {
	Name     string // Template name for ExecuteTemplate
	Template string // Path within the templates FS (relative to "files/")
}

// loadTemplates parses one or more sets of template entries into the given template.
// Each set is a slice of templateEntry. All entries are loaded in order.
func loadTemplates(tmpl *template.Template, sets ...[]templateEntry) error {
	for _, set := range sets {
		for _, entry := range set {
			content, err := templates.TemplateFS.ReadFile("files/" + entry.Template)
			if err != nil {
				return fmt.Errorf("reading template %s: %w", entry.Template, err)
			}
			if _, err := tmpl.New(entry.Name).Parse(string(content)); err != nil {
				return fmt.Errorf("parsing template %s: %w", entry.Template, err)
			}
		}
	}
	return nil
}

// clientTemplateEntries converts ClientTemplates map to a slice of templateEntry.
func clientTemplateEntries() []templateEntry {
	entries := make([]templateEntry, 0, len(templates.ClientTemplates))
	for _, ct := range templates.ClientTemplates {
		entries = append(entries, templateEntry{Name: ct.Name, Template: ct.Template})
	}
	return entries
}

// initiatorTemplateEntries converts InitiatorTemplates map to a slice of templateEntry.
func initiatorTemplateEntries() []templateEntry {
	entries := make([]templateEntry, 0, len(templates.InitiatorTemplates))
	for _, it := range templates.InitiatorTemplates {
		entries = append(entries, templateEntry{Name: it.Name, Template: it.Template})
	}
	return entries
}

// senderTemplateEntries converts SenderTemplates map to a slice of templateEntry.
func senderTemplateEntries() []templateEntry {
	entries := make([]templateEntry, 0, len(templates.SenderTemplates))
	for _, st := range templates.SenderTemplates {
		entries = append(entries, templateEntry{Name: st.Name, Template: st.Template})
	}
	return entries
}

// SenderTemplateData is the unified template data for client and initiator templates.
// Templates use {{if .IsClient}} to branch on the few points where they diverge.
type SenderTemplateData struct {
	IsClient    bool                   // true for client, false for initiator
	Prefix      string                 // "" for client, "Webhook"/"Callback" for initiator
	PrefixLower string                 // "" for client, "webhook"/"callback" for initiator
	TypeName    string                 // "Client" or "WebhookInitiator"
	Receiver    string                 // "c" or "p"
	OptionType  string                 // "ClientOption" or "WebhookInitiatorOption"
	ErrorType   string                 // "ClientHttpError" or "WebhookHttpError"
	SimpleType  string                 // "SimpleClient" or "SimpleWebhookInitiator"
	Operations  []*OperationDescriptor // Operations to generate for
}

// sharedServerTemplateEntries converts SharedServerTemplates map to a slice of templateEntry.
func sharedServerTemplateEntries() []templateEntry {
	entries := make([]templateEntry, 0, len(templates.SharedServerTemplates))
	for _, st := range templates.SharedServerTemplates {
		entries = append(entries, templateEntry{Name: st.Name, Template: st.Template})
	}
	return entries
}

// ClientGenerator generates client code from operation descriptors.
type ClientGenerator struct {
	tmpl           *template.Template
	schemaIndex    map[string]*SchemaDescriptor
	generateSimple bool
	modelsPackage  *ModelsPackage
}

// NewClientGenerator creates a new client generator.
// modelsPackage can be nil if models are in the same package.
// rp holds the package prefixes for runtime sub-packages; all empty when embedded.
func NewClientGenerator(schemaIndex map[string]*SchemaDescriptor, generateSimple bool, modelsPackage *ModelsPackage, rp RuntimePrefixes, typeMapping TypeMapping) (*ClientGenerator, error) {
	tmpl := template.New("client").Funcs(templates.Funcs()).Funcs(clientFuncs(schemaIndex, modelsPackage, typeMapping)).Funcs(rp.FuncMap())

	if err := loadTemplates(tmpl, clientTemplateEntries(), senderTemplateEntries(), sharedServerTemplateEntries()); err != nil {
		return nil, err
	}

	return &ClientGenerator{
		tmpl:           tmpl,
		schemaIndex:    schemaIndex,
		generateSimple: generateSimple,
		modelsPackage:  modelsPackage,
	}, nil
}

// clientFuncs returns template functions specific to client generation.
func clientFuncs(schemaIndex map[string]*SchemaDescriptor, modelsPackage *ModelsPackage, typeMapping TypeMapping) template.FuncMap {
	return template.FuncMap{
		"pathFmt":                        pathFmt,
		"isSimpleOperation":              isSimpleOperation,
		"simpleOperationSuccessResponse": simpleOperationSuccessResponse,
		"errorResponseForOperation":      errorResponseForOperation,
		"defaultTypedBody": func(op *OperationDescriptor) *RequestBodyDescriptor {
			return op.DefaultTypedBody()
		},
		"goTypeForContent": func(content *ResponseContentDescriptor) string {
			return goTypeForContent(content, schemaIndex, modelsPackage, typeMapping)
		},
		"modelsPkg": func() string {
			return modelsPackage.Prefix()
		},
	}
}

// pathFmt converts a path with {param} placeholders to a format string.
// Example: "/pets/{petId}" -> "/pets/%s"
func pathFmt(path string) string {
	result := path
	for {
		start := strings.Index(result, "{")
		if start == -1 {
			break
		}
		end := strings.Index(result, "}")
		if end == -1 {
			break
		}
		result = result[:start] + "%s" + result[end+1:]
	}
	return result
}

// isSimpleOperation returns true if an operation has a single JSON success response type.
// "Simple" operations can have typed wrapper methods in SimpleClient.
func isSimpleOperation(op *OperationDescriptor) bool {
	// Must have responses
	if len(op.Responses) == 0 {
		return false
	}

	// If the operation has a body, it must have a typed body for the simple client
	if op.HasBody && !op.HasTypedBody() {
		return false
	}

	// Count success responses (2xx or default that could be success)
	var successResponses []*ResponseDescriptor
	for _, r := range op.Responses {
		if strings.HasPrefix(r.StatusCode, "2") {
			successResponses = append(successResponses, r)
		}
	}

	// Must have exactly one success response
	if len(successResponses) != 1 {
		return false
	}

	success := successResponses[0]

	// Must have at least one content type and exactly one JSON content type
	// (i.e., if there are multiple content types, we can't have a simple typed response)
	if len(success.Contents) == 0 {
		return false
	}
	if len(success.Contents) != 1 {
		return false
	}

	// The single content type must be JSON
	return success.Contents[0].IsJSON
}

// simpleOperationSuccessResponse returns the single success response for a simple operation.
func simpleOperationSuccessResponse(op *OperationDescriptor) *ResponseDescriptor {
	for _, r := range op.Responses {
		if strings.HasPrefix(r.StatusCode, "2") {
			return r
		}
	}
	return nil
}

// errorResponseForOperation returns the error response (default or 4xx/5xx) if one exists.
func errorResponseForOperation(op *OperationDescriptor) *ResponseDescriptor {
	// First, look for a default response
	for _, r := range op.Responses {
		if r.StatusCode == "default" {
			if len(r.Contents) > 0 && r.Contents[0].IsJSON {
				return r
			}
		}
	}
	// Then look for a 4xx or 5xx response
	for _, r := range op.Responses {
		if strings.HasPrefix(r.StatusCode, "4") || strings.HasPrefix(r.StatusCode, "5") {
			if len(r.Contents) > 0 && r.Contents[0].IsJSON {
				return r
			}
		}
	}
	return nil
}

// goTypeForContent returns the Go type for a response content descriptor.
// If modelsPackage is set, type names are prefixed with the package name.
func goTypeForContent(content *ResponseContentDescriptor, schemaIndex map[string]*SchemaDescriptor, modelsPackage *ModelsPackage, typeMapping TypeMapping) string {
	if content == nil || content.Schema == nil {
		return "any"
	}

	pkgPrefix := modelsPackage.Prefix()

	// If the schema has a reference, look it up
	if content.Schema.Ref != "" {
		if target, ok := schemaIndex[content.Schema.Ref]; ok {
			return pkgPrefix + target.ShortName
		}
	}

	// Check if this is an array schema with items that have a reference
	if content.Schema.Schema != nil && content.Schema.Schema.Items != nil {
		itemProxy := content.Schema.Schema.Items.A
		if itemProxy != nil && itemProxy.IsReference() {
			ref := itemProxy.GetReference()
			if target, ok := schemaIndex[ref]; ok {
				return "[]" + pkgPrefix + target.ShortName
			}
		}
	}

	// If the schema has a short name, use it
	if content.Schema.ShortName != "" {
		return pkgPrefix + content.Schema.ShortName
	}

	// Fall back to the stable name
	if content.Schema.StableName != "" {
		return pkgPrefix + content.Schema.StableName
	}

	// Try to derive from the schema itself using TypeMapping
	if content.Schema.Schema != nil {
		return resolveSchemaType(content.Schema.Schema, typeMapping)
	}

	return "any"
}

// resolveSchemaType converts a schema to a Go type string using the provided TypeMapping.
// This is a standalone helper used as a last-resort fallback in goTypeForContent.
func resolveSchemaType(schema *base.Schema, tm TypeMapping) string {
	if schema == nil {
		return "any"
	}

	// Check for array
	if schema.Items != nil && schema.Items.A != nil {
		itemType := "any"
		if itemSchema := schema.Items.A.Schema(); itemSchema != nil {
			itemType = resolveSchemaType(itemSchema, tm)
		}
		return "[]" + itemType
	}

	// Check explicit type
	for _, t := range schema.Type {
		switch t {
		case "string":
			return resolveFormatType(tm.String, schema.Format)
		case "integer":
			return resolveFormatType(tm.Integer, schema.Format)
		case "number":
			return resolveFormatType(tm.Number, schema.Format)
		case "boolean":
			return tm.Boolean.Default.Type
		case "array":
			return "[]any"
		case "object":
			return "map[string]any"
		}
	}

	return "any"
}

// resolveFormatType looks up a type from a FormatMapping, falling back to its default.
func resolveFormatType(fm FormatMapping, format string) string {
	if format != "" {
		if spec, ok := fm.Formats[format]; ok {
			return spec.Type
		}
	}
	return fm.Default.Type
}

// GenerateBase generates the base client types and helpers.
func (g *ClientGenerator) GenerateBase() (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "base", nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateInterface generates the ClientInterface.
func (g *ClientGenerator) GenerateInterface(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_interface", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateMethods generates the Client methods.
func (g *ClientGenerator) GenerateMethods(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_methods", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateRequestBuilders generates the request builder functions.
func (g *ClientGenerator) GenerateRequestBuilders(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_request_builders", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateSimple generates the SimpleClient with typed responses.
func (g *ClientGenerator) GenerateSimple(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_simple", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateParamTypes generates the parameter struct types.
func (g *ClientGenerator) GenerateParamTypes(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "param_types", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateRequestBodyTypes generates type aliases for request bodies.
func (g *ClientGenerator) GenerateRequestBodyTypes(ops []*OperationDescriptor) string {
	return generateRequestBodyTypes(ops, g.schemaIndex, g.modelsPackage)
}

// generateRequestBodyTypes generates type aliases for request bodies.
// This is shared between ClientGenerator and InitiatorGenerator.
func generateRequestBodyTypes(ops []*OperationDescriptor, schemaIndex map[string]*SchemaDescriptor, modelsPackage *ModelsPackage) string {
	var buf bytes.Buffer
	pkgPrefix := modelsPackage.Prefix()

	for _, op := range ops {
		for _, body := range op.Bodies {
			if !body.GenerateTyped {
				continue
			}
			// Get the underlying type for this request body
			var targetType string
			if body.Schema != nil {
				if body.Schema.Ref != "" {
					// Reference to a component schema
					if target, ok := schemaIndex[body.Schema.Ref]; ok {
						targetType = pkgPrefix + target.ShortName
					}
				} else if body.Schema.ShortName != "" {
					targetType = pkgPrefix + body.Schema.ShortName
				}
			}
			if targetType == "" {
				targetType = "any"
			}

			fmt.Fprintf(&buf, "type %s = %s\n\n", body.GoTypeName, targetType)
		}
	}

	return buf.String()
}

// GenerateClient generates the complete client code.
func (g *ClientGenerator) GenerateClient(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer

	data := SenderTemplateData{
		IsClient:   true,
		Prefix:     "",
		PrefixLower: "",
		TypeName:   "Client",
		Receiver:   "c",
		OptionType: "ClientOption",
		ErrorType:  "ClientHttpError",
		SimpleType: "SimpleClient",
		Operations: ops,
	}

	// Generate request body type aliases first
	bodyTypes := g.GenerateRequestBodyTypes(ops)
	buf.WriteString(bodyTypes)

	// Generate base client
	base, err := g.GenerateBase()
	if err != nil {
		return "", fmt.Errorf("generating base client: %w", err)
	}
	buf.WriteString(base)
	buf.WriteString("\n")

	// Generate interface
	iface, err := g.GenerateInterface(data)
	if err != nil {
		return "", fmt.Errorf("generating client interface: %w", err)
	}
	buf.WriteString(iface)
	buf.WriteString("\n")

	// Generate param types
	paramTypes, err := g.GenerateParamTypes(ops)
	if err != nil {
		return "", fmt.Errorf("generating param types: %w", err)
	}
	buf.WriteString(paramTypes)
	buf.WriteString("\n")

	// Generate methods
	methods, err := g.GenerateMethods(data)
	if err != nil {
		return "", fmt.Errorf("generating client methods: %w", err)
	}
	buf.WriteString(methods)
	buf.WriteString("\n")

	// Generate request builders
	builders, err := g.GenerateRequestBuilders(data)
	if err != nil {
		return "", fmt.Errorf("generating request builders: %w", err)
	}
	buf.WriteString(builders)
	buf.WriteString("\n")

	// Generate simple client if requested
	if g.generateSimple {
		simple, err := g.GenerateSimple(data)
		if err != nil {
			return "", fmt.Errorf("generating simple client: %w", err)
		}
		buf.WriteString(simple)
	}

	return buf.String(), nil
}
