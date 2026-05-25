package presenter

import (
	"encoding/json"
	"io"
)

// WriteJSON writes a deterministic, newline-terminated JSON value to the output stream.
func WriteJSON(w io.Writer, value any) error {
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	return encoder.Encode(value)
}
