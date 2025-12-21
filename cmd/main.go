package main

import (
	"fmt"
	"os"

	"github.com/tpm1qq/gtm/internal/app"
)

func main() {
	if err := app.RunGTM(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
