package capability_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
)

func TestDefinitionInputSchemaJSON(t *testing.T) {
	t.Parallel()

	definition := capability.Definition{
		ID:          "orders.get",
		Title:       "Get Order",
		Description: "Gets an order by ID.",
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "id", Type: capability.InputTypeString, Required: true, Description: "Order identifier."},
		}},
	}

	content, err := json.Marshal(definition)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	rendered := string(content)
	for _, fragment := range []string{`"inputSchema"`, `"fields"`, `"name":"id"`, `"type":"string"`, `"required":true`} {
		if !strings.Contains(rendered, fragment) {
			t.Fatalf("definition JSON missing %s in %s", fragment, rendered)
		}
	}
}

func TestDefinitionOmitsAbsentInputSchema(t *testing.T) {
	t.Parallel()

	content, err := json.Marshal(capability.Definition{ID: "foundation.example"})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(content), "inputSchema") {
		t.Fatalf("definition JSON should omit absent input schema: %s", string(content))
	}
}
