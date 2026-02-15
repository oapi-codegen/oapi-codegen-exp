package codegen

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// ServerGenerator generates server code from operation descriptors.
type ServerGenerator struct {
	tmpl       *template.Template
	serverType string
}

// NewServerGenerator creates a new server generator for the specified server type.
// rp holds the package prefixes for runtime sub-packages; all empty when embedded.
func NewServerGenerator(serverType string, rp RuntimePrefixes) (*ServerGenerator, error) {
	if serverType == "" {
		// No server generation requested
		return &ServerGenerator{serverType: ""}, nil
	}

	tmpl := template.New("server").Funcs(templates.Funcs()).Funcs(rp.FuncMap())

	// Get templates for the specified server type
	serverTemplates, err := getServerTemplates(serverType)
	if err != nil {
		return nil, err
	}

	// Convert server-specific templates to entries
	serverEntries := make([]templateEntry, 0, len(serverTemplates))
	for _, st := range serverTemplates {
		serverEntries = append(serverEntries, templateEntry{Name: st.Name, Template: st.Template})
	}

	if err := loadTemplates(tmpl, serverEntries, sharedServerTemplateEntries()); err != nil {
		return nil, err
	}

	return &ServerGenerator{tmpl: tmpl, serverType: serverType}, nil
}

// getServerTemplates returns the templates for the specified server type.
func getServerTemplates(serverType string) (map[string]templates.ServerTemplate, error) {
	switch serverType {
	case ServerTypeStdHTTP:
		return templates.StdHTTPServerTemplates, nil
	case ServerTypeChi:
		return templates.ChiServerTemplates, nil
	case ServerTypeEcho:
		return templates.EchoServerTemplates, nil
	case ServerTypeEchoV4:
		return templates.EchoV4ServerTemplates, nil
	case ServerTypeGin:
		return templates.GinServerTemplates, nil
	case ServerTypeGorilla:
		return templates.GorillaServerTemplates, nil
	case ServerTypeFiber:
		return templates.FiberServerTemplates, nil
	case ServerTypeIris:
		return templates.IrisServerTemplates, nil
	default:
		return nil, fmt.Errorf("unsupported server type: %q (supported: %q, %q, %q, %q, %q, %q, %q, %q)",
			serverType,
			ServerTypeStdHTTP, ServerTypeChi, ServerTypeEcho, ServerTypeEchoV4, ServerTypeGin,
			ServerTypeGorilla, ServerTypeFiber, ServerTypeIris)
	}
}

// GenerateInterface generates the ServerInterface.
func (g *ServerGenerator) GenerateInterface(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "interface", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateHandler generates the HTTP handler and routing code.
func (g *ServerGenerator) GenerateHandler(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "handler", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateWrapper generates the ServerInterfaceWrapper.
func (g *ServerGenerator) GenerateWrapper(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "wrapper", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateErrors generates the error types.
func (g *ServerGenerator) GenerateErrors() (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "errors", nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateParamTypes generates the parameter struct types.
func (g *ServerGenerator) GenerateParamTypes(ops []*OperationDescriptor) (string, error) {
	var buf bytes.Buffer
	if err := g.tmpl.ExecuteTemplate(&buf, "param_types", ops); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GenerateServer generates all server code components.
// Returns empty string if no server type was configured.
func (g *ServerGenerator) GenerateServer(ops []*OperationDescriptor) (string, error) {
	if g.serverType == "" || g.tmpl == nil {
		return "", nil
	}

	var buf bytes.Buffer

	// Generate interface
	iface, err := g.GenerateInterface(ops)
	if err != nil {
		return "", err
	}
	buf.WriteString(iface)
	buf.WriteString("\n")

	// Generate param types
	paramTypes, err := g.GenerateParamTypes(ops)
	if err != nil {
		return "", err
	}
	buf.WriteString(paramTypes)
	buf.WriteString("\n")

	// Generate wrapper
	wrapper, err := g.GenerateWrapper(ops)
	if err != nil {
		return "", err
	}
	buf.WriteString(wrapper)
	buf.WriteString("\n")

	// Generate handler
	handler, err := g.GenerateHandler(ops)
	if err != nil {
		return "", err
	}
	buf.WriteString(handler)
	buf.WriteString("\n")

	// Generate errors
	errors, err := g.GenerateErrors()
	if err != nil {
		return "", err
	}
	buf.WriteString(errors)

	return buf.String(), nil
}
