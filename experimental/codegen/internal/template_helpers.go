package codegen

import (
	"strings"
	"text/template"
)

// senderFuncs returns template functions for the shared sender templates.
// These eliminate complex inline conditionals in templates by moving
// name composition, signature building, and comment generation into Go.
func senderFuncs() template.FuncMap {
	return template.FuncMap{
		"methodName":              senderMethodName,
		"typedMethodName":         senderTypedMethodName,
		"requestBuilderName":      senderRequestBuilderName,
		"typedRequestBuilderName": senderTypedRequestBuilderName,
		"methodParams":            senderMethodParams,
		"methodArgs":              senderMethodArgs,
		"methodCallArgs":          senderMethodCallArgs,
		"requestBuilderParams":    senderRequestBuilderParams,
		"requestBuilderArgs":      senderRequestBuilderArgs,
		"methodComment":           senderMethodComment,
		"typedMethodComment":      senderTypedMethodComment,
		"requestBuilderComment":   senderRequestBuilderComment,
	}
}

// --- Name composition ---

// senderMethodName returns the Go method name for an operation.
//
//	"FindPets" or "FindPetsWithBody"
func senderMethodName(op *OperationDescriptor) string {
	if op.HasBody {
		return op.GoOperationID + "WithBody"
	}
	return op.GoOperationID
}

// senderTypedMethodName returns the Go method name for a typed-body variant.
//
//	"FindPets" or "FindPetsWithFormBody"
func senderTypedMethodName(op *OperationDescriptor, body *RequestBodyDescriptor) string {
	return op.GoOperationID + body.FuncSuffix
}

// senderRequestBuilderName returns the free function name for building an HTTP request.
//
//	"NewFindPetsRequestWithBody" or "NewFindPetsWebhookRequestWithBody"
func senderRequestBuilderName(data SenderTemplateData, op *OperationDescriptor) string {
	name := "New" + op.GoOperationID + data.Prefix + "Request"
	if op.HasBody {
		name += "WithBody"
	}
	return name
}

// senderTypedRequestBuilderName returns the free function name for building a typed-body request.
//
//	"NewFindPetsWebhookRequestWithFormBody"
func senderTypedRequestBuilderName(data SenderTemplateData, op *OperationDescriptor, body *RequestBodyDescriptor) string {
	return "New" + op.GoOperationID + data.Prefix + "Request" + body.FuncSuffix
}

// --- Signature fragments ---
//
// These eliminate {{if $.IsClient}}{{range}}…{{else}}…{{end}} nesting
// by encoding the client-vs-initiator difference (path params vs targetURL)
// in Go rather than template conditionals.

// senderMethodParams returns the parameter declarations after "ctx context.Context"
// in a method signature.
//
//	Client:    ", id string, name string"
//	Initiator: ", targetURL string"
func senderMethodParams(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		var buf strings.Builder
		for _, p := range op.PathParams {
			buf.WriteString(", ")
			buf.WriteString(p.GoVariableName())
			buf.WriteString(" ")
			buf.WriteString(p.TypeDecl)
		}
		return buf.String()
	}
	return ", targetURL string"
}

// senderMethodArgs returns the forwarding arguments from a method to a request builder.
//
//	Client:    "c.Server, id, name"
//	Initiator: "targetURL"
func senderMethodArgs(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		var buf strings.Builder
		buf.WriteString(data.Receiver)
		buf.WriteString(".Server")
		for _, p := range op.PathParams {
			buf.WriteString(", ")
			buf.WriteString(p.GoVariableName())
		}
		return buf.String()
	}
	return "targetURL"
}

// senderMethodCallArgs returns the call-site arguments for delegating to a peer method
// (used by simple.go.tmpl where the receiver is already implicit).
//
//	Client:    ", id, name"
//	Initiator: ", targetURL"
func senderMethodCallArgs(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		var buf strings.Builder
		for _, p := range op.PathParams {
			buf.WriteString(", ")
			buf.WriteString(p.GoVariableName())
		}
		return buf.String()
	}
	return ", targetURL"
}

// senderRequestBuilderParams returns the parameter declarations for a NewXxxRequest function.
//
//	Client:    "server string, id string, name string"
//	Initiator: "targetURL string"
func senderRequestBuilderParams(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		var buf strings.Builder
		buf.WriteString("server string")
		for _, p := range op.PathParams {
			buf.WriteString(", ")
			buf.WriteString(p.GoVariableName())
			buf.WriteString(" ")
			buf.WriteString(p.TypeDecl)
		}
		return buf.String()
	}
	return "targetURL string"
}

// senderRequestBuilderArgs returns the forwarding arguments between request builders.
//
//	Client:    "server, id, name"
//	Initiator: "targetURL"
func senderRequestBuilderArgs(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		var buf strings.Builder
		buf.WriteString("server")
		for _, p := range op.PathParams {
			buf.WriteString(", ")
			buf.WriteString(p.GoVariableName())
		}
		return buf.String()
	}
	return "targetURL"
}

// --- Comment helpers ---

// senderMethodComment returns the doc comment suffix for a method.
//
//	Client:    " makes a GET request to /pets"
//	Initiator: " sends a GET webhook request"
func senderMethodComment(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		return " makes a " + op.Method + " request to " + op.Path
	}
	return " sends a " + op.Method + " " + data.PrefixLower + " request"
}

// senderTypedMethodComment returns the doc comment suffix for a typed-body method.
//
//	Client:    " makes a POST request to /pets with application/json body"
//	Initiator: " sends a POST webhook request with application/json body"
func senderTypedMethodComment(data SenderTemplateData, op *OperationDescriptor, body *RequestBodyDescriptor) string {
	if data.IsClient {
		return " makes a " + op.Method + " request to " + op.Path + " with " + body.ContentType + " body"
	}
	return " sends a " + op.Method + " " + data.PrefixLower + " request with " + body.ContentType + " body"
}

// senderRequestBuilderComment returns the doc comment for a request builder function.
//
//	Client:    "creates a GET request for /pets"
//	Initiator: "creates a GET request for the webhook"
func senderRequestBuilderComment(data SenderTemplateData, op *OperationDescriptor) string {
	if data.IsClient {
		return "creates a " + op.Method + " request for " + op.Path
	}
	return "creates a " + op.Method + " request for the " + data.PrefixLower
}
