package presenter_test

import (
	"bytes"
	"testing"

	"github.com/yargotev/exito-tools/internal/presenter"
)

func TestWriteJSONTerminatesWithNewline(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := presenter.WriteJSON(&output, map[string]bool{"ok": true}); err != nil {
		t.Fatalf("WriteJSON() error = %v", err)
	}

	if got, want := output.String(), "{\"ok\":true}\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
}
