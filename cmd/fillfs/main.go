// Package main provides the fillfs CLI entrypoint.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/thorstenkramm/fillfs/internal/app"
	"github.com/thorstenkramm/fillfs/internal/options"
	"github.com/thorstenkramm/fillfs/internal/runerr"
)

func main() {
	cfg, err := options.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := app.Run(context.Background(), cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(runerr.Code(err, 1))
	}
}
