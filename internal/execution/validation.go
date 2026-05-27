package execution

import (
	"fmt"
	"math"

	"github.com/yargotev/exito-tools/internal/capability"
)

// ValidateInput checks a complete Capability input object against the neutral schema metadata.
func ValidateInput(input capability.Input, schema *capability.InputSchema) error {
	if schema == nil {
		return nil
	}

	if input == nil {
		input = capability.Input{}
	}

	for _, field := range schema.Fields {
		value, exists := input[field.Name]
		if !exists || value == nil {
			if field.Required {
				return capability.StructuredError{
					Code:    ErrorInvalidInput,
					Message: fmt.Sprintf("Required input field %q is missing.", field.Name),
				}
			}
			continue
		}

		if !inputValueMatchesType(value, field.Type) {
			return capability.StructuredError{
				Code:    ErrorInvalidInput,
				Message: fmt.Sprintf("Input field %q must be %s.", field.Name, field.Type),
			}
		}
	}

	return nil
}

func inputValueMatchesType(value any, inputType capability.InputType) bool {
	switch inputType {
	case capability.InputTypeString:
		_, ok := value.(string)
		return ok
	case capability.InputTypeNumber:
		return isNumber(value)
	case capability.InputTypeBoolean:
		_, ok := value.(bool)
		return ok
	case capability.InputTypeObject:
		_, ok := value.(map[string]any)
		return ok
	case capability.InputTypeArray:
		switch value.(type) {
		case []any, []string:
			return true
		default:
			return false
		}
	default:
		return true
	}
}

func isNumber(value any) bool {
	switch number := value.(type) {
	case float64:
		return !math.IsNaN(number) && !math.IsInf(number, 0)
	case float32:
		return !math.IsNaN(float64(number)) && !math.IsInf(float64(number), 0)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}
