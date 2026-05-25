package main

import (
	"context"
	"log"
	"os"

	"github.com/yargotev/exito-tools/internal/app"
	clisurface "github.com/yargotev/exito-tools/internal/surface/cli"
)

func main() {
	root := clisurface.NewRoot(app.New)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	if err := root.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
