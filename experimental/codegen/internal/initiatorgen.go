package codegen

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// InitiatorGenerator generates initiator (sender) code from operation descriptors.
// It is parameterized by prefix to support both webhooks and callbacks.
type InitiatorGenerator struct {
	tmpl           *template.Template
	prefix         string // "Webhook" or "Callback"
	schemaIndex    map[string]*SchemaDescriptor
	generateSimple bool
	modelsPackage  *ModelsPackage
}

// NewInitiatorGenerator creates a new initiator generator.
// rp holds the package prefixes for runtime sub-packages; all empty when embedded.
func NewInitiatorGenerator(prefix string, schemaIndex map[string]*SchemaDescriptor, generateSimple bool, modelsPackage *ModelsPackage, rp RuntimePrefixes, typeMapping TypeMapping) (*InitiatorGenerator, error) {
	tmpl := template.New("initiator").Funcs(templates.Funcs()).Funcs(senderFuncs()).Funcs(clientFuncs(schemaIndex, modelsPackage, typeMapping)).Funcs(rp.FuncMap())

	if err := loadTemplates(tmpl, initiatorTemplateEntries(), senderTemplateEntries(), sharedServerTemplateEntries()); err != nil {
		return nil, err
	}

	return &InitiatorGenerator{
		tmpl:           tmpl,
		prefix:         prefix,
		schemaIndex:    schemaIndex,
		generateSimple: generateSimple,
		modelsPackage:  modelsPackage,
	}, nil
}

func (g *InitiatorGenerator) senderData(ops []*OperationDescriptor) SenderTemplateData {
	return SenderTemplateData{
		IsClient:    false,
		Prefix:      g.prefix,
		PrefixLower: strings.ToLower(g.prefix),
		TypeName:    g.prefix + "Initiator",
		Receiver:    "p",
		OptionType:  g.prefix + "InitiatorOption",
		ErrorType:   g.prefix + "HttpError",
		SimpleType:  "Simple" + g.prefix + "Initiator",
		Operations:  ops,
	}
}

// baseTemplateData builds the data for the initiator base template, which still
// uses the old shape (Prefix, PrefixLower, Operations).
type initiatorBaseData struct {
	Prefix      string
	PrefixLower string
	Operations  []*OperationDescriptor
}

// GenerateBase generates the base initiator types and helpers.
func (g *InitiatorGenerator) GenerateBase(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	data := initiatorBaseData{
		Prefix:      g.prefix,
		PrefixLower: strings.ToLower(g.prefix),
		Operations:  ops,
	}
	if err := g.tmpl.ExecuteTemplate(&buf, "initiator_base", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateInterface generates the InitiatorInterface.
func (g *InitiatorGenerator) GenerateInterface(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_interface", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateMethods generates the Initiator methods.
func (g *InitiatorGenerator) GenerateMethods(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_methods", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateRequestBuilders generates the request builder functions.
func (g *InitiatorGenerator) GenerateRequestBuilders(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_request_builders", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateSimple generates the SimpleInitiator with typed responses.
func (g *InitiatorGenerator) GenerateSimple(data SenderTemplateData) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "sender_simple", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateParamTypes generates the parameter struct types.
func (g *InitiatorGenerator) GenerateParamTypes(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "param_types", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateRequestBodyTypes generates type aliases for request bodies.
func (g *InitiatorGenerator) GenerateRequestBodyTypes(ops []*OperationDescriptor) string {
	return generateRequestBodyTypes(ops, g.schemaIndex, g.modelsPackage)
}

// GenerateInitiator generates the complete initiator code.
func (g *InitiatorGenerator) GenerateInitiator(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer

	data := g.senderData(ops)

	// Generate request body type aliases first
	bodyTypes := g.GenerateRequestBodyTypes(ops)
	buf.WriteString(bodyTypes)

	// Generate base initiator
	base, err := g.GenerateBase(ops)
	if err != nil {
		return "", fmt.Errorf("generating base initiator: %w", err)
	}
	buf.WriteString(base)
	buf.WriteString("\n")

	// Generate interface
	iface, err := g.GenerateInterface(data)
	if err != nil {
		return "", fmt.Errorf("generating initiator interface: %w", err)
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
		return "", fmt.Errorf("generating initiator methods: %w", err)
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

	// Generate simple initiator if requested
	if g.generateSimple {
		simple, err := g.GenerateSimple(data)
		if err != nil {
			return "", fmt.Errorf("generating simple initiator: %w", err)
		}
		buf.WriteString(simple)
	}

	return buf.String(), nil
}
