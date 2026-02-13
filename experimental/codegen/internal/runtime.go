package codegen

import (
	"fmt"
	"sort"

	"github.com/oapi-codegen/oapi-codegen-exp/experimental/codegen/internal/templates"
)

// RuntimeOutput holds the generated Go source code for each runtime sub-package.
type RuntimeOutput struct {
	Params  string // params sub-package (style/bind functions, helpers)
	Types   string // types sub-package (Date, Email, UUID, File, Nullable)
	Helpers string // helpers sub-package (MarshalForm)
}

// GenerateRuntime produces standalone Go source files for each of the three
// runtime sub-packages. baseImportPath is the base import path for the runtime
// module (e.g., "github.com/org/project/runtime"). The params sub-package
// imports the types sub-package for Date references.
func GenerateRuntime(baseImportPath string) (*RuntimeOutput, error) {
	if baseImportPath == "" {
		return nil, fmt.Errorf("base import path is required")
	}

	typesCode, err := generateRuntimeTypes()
	if err != nil {
		return nil, fmt.Errorf("generating runtime types: %w", err)
	}

	paramsCode, err := generateRuntimeParams(baseImportPath + "/types")
	if err != nil {
		return nil, fmt.Errorf("generating runtime params: %w", err)
	}

	helpersCode, err := generateRuntimeHelpers()
	if err != nil {
		return nil, fmt.Errorf("generating runtime helpers: %w", err)
	}

	return &RuntimeOutput{
		Params:  paramsCode,
		Types:   typesCode,
		Helpers: helpersCode,
	}, nil
}

// generateRuntimeTypes produces the types sub-package: Date, Email, UUID, File, Nullable.
func generateRuntimeTypes() (string, error) {
	ctx := NewCodegenContext()

	// Mark ALL custom types as needed.
	for name := range templates.TypeTemplates {
		ctx.NeedCustomType(name)
	}

	output := NewOutput("types")

	// Emit custom type templates in sorted order for deterministic output.
	typeNames := make([]string, 0, len(templates.TypeTemplates))
	for name := range templates.TypeTemplates {
		typeNames = append(typeNames, name)
	}
	sort.Strings(typeNames)

	for _, name := range typeNames {
		def := templates.TypeTemplates[name]
		templateName := def.Template
		parts := splitTemplatePath(templateName)
		typeCode := ctx.loadAndRegisterCustomType(parts)
		if typeCode != "" {
			output.AddType(typeCode)
		}
	}

	output.AddImports(ctx.Imports())
	return output.Format()
}

// generateRuntimeParams produces the params sub-package: style/bind functions
// and shared helpers (ParamLocation, primitiveToString, BindStringToObject, etc.).
// typesImportPath is the import path for the types sub-package so that
// Date/DateFormat references can be qualified with "types.".
func generateRuntimeParams(typesImportPath string) (string, error) {
	ctx := NewCodegenContext()

	// Mark ALL param style/bind combinations as needed.
	for key := range templates.ParamTemplates {
		ctx.params[key] = true
	}

	output := NewOutput("params")

	// Emit param functions with typesPrefix = "types." so Date references
	// are qualified as types.Date, types.DateFormat.
	paramFuncs, err := generateParamFunctionsFromContext(ctx, "types.")
	if err != nil {
		return "", fmt.Errorf("generating param functions: %w", err)
	}
	if paramFuncs != "" {
		output.AddType(paramFuncs)
	}

	output.AddImports(ctx.Imports())
	// Add the types sub-package import for Date/DateFormat references.
	output.AddImport(typesImportPath, "")
	return output.Format()
}

// generateRuntimeHelpers produces the helpers sub-package: MarshalForm.
func generateRuntimeHelpers() (string, error) {
	ctx := NewCodegenContext()

	ctx.NeedHelper("marshal_form")

	output := NewOutput("helpers")

	for _, helperName := range ctx.RequiredHelpers() {
		helperCode, err := generateHelper(helperName, ctx)
		if err != nil {
			return "", fmt.Errorf("generating helper %s: %w", helperName, err)
		}
		if helperCode != "" {
			output.AddType(helperCode)
		}
	}

	output.AddImports(ctx.Imports())
	return output.Format()
}

// splitTemplatePath extracts the filename from a template path like "types/date.tmpl".
func splitTemplatePath(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
