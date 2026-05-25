package main

import (
	"context"
	"log"
	"os"

	"github.com/yargotev/exito-tools/internal/app"
	clisurface "github.com/yargotev/exito-tools/internal/surface/cli"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	root := clisurface.NewRoot(application)
	root.SetOut(os.Stdout)
	root.SetErr(os.Stderr)

	if err := root.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
}
