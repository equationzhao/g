package main

import (
	"os"

	"github.com/Equationzhao/g/internal/app"
)

func main() {
	os.Exit(app.Run(os.Args[1:], os.Environ(), app.OSDeps()))
}
