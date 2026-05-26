package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yargotev/exito-tools/internal/app"
	clisurface "github.com/yargotev/exito-tools/internal/surface/cli"
)

func main() {
	root := clisurface.NewRoot(app.New)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	if err := root.ExecuteContext(context.Background()); err != nil {
		if !clisurface.IsExitError(err) {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(clisurface.ExitCode(err))
	}
}
