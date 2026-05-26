package execution_test

import (
	"errors"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/execution"
)

func TestValidateInput(t *testing.T) {
	t.Parallel()

	schema := &capability.InputSchema{Fields: []capability.InputField{
		{Name: "id", Type: capability.InputTypeString, Required: true},
		{Name: "count", Type: capability.InputTypeNumber},
		{Name: "active", Type: capability.InputTypeBoolean},
		{Name: "metadata", Type: capability.InputTypeObject},
		{Name: "tags", Type: capability.InputTypeArray},
	}}

	tests := []struct {
		name    string
		input   capability.Input
		wantErr string
	}{
		{
			name: "valid complete object",
			input: capability.Input{
				"id":       "A123",
				"count":    float64(2),
				"active":   true,
				"metadata": map[string]any{"source": "test"},
				"tags":     []any{"priority"},
			},
		},
		{
			name:  "optional fields may be absent",
			input: capability.Input{"id": "A123"},
		},
		{
			name:    "missing required field",
			input:   capability.Input{"count": 1},
			wantErr: execution.ErrorInvalidInput,
		},
		{
			name:    "wrong field type",
			input:   capability.Input{"id": 123},
			wantErr: execution.ErrorInvalidInput,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := execution.ValidateInput(tc.input, schema)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("ValidateInput() error = %v, want nil", err)
				}
				return
			}

			var structured capability.StructuredError
			if !errors.As(err, &structured) {
				t.Fatalf("ValidateInput() error = %T, want StructuredError", err)
			}
			if structured.Code != tc.wantErr {
				t.Fatalf("StructuredError.Code = %q, want %q", structured.Code, tc.wantErr)
			}
		})
	}
}

func TestValidateInputSkipsAbsentSchema(t *testing.T) {
	t.Parallel()

	if err := execution.ValidateInput(nil, nil); err != nil {
		t.Fatalf("ValidateInput() error = %v, want nil", err)
	}
}
