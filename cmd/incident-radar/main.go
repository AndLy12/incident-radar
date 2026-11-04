// Package main implements the Incident Radar CLI application.
package main

import (
	"fmt"
	"incidentradar/internal/app"
	"os"
)

func main() {
	args := os.Args[1:]
	if err := app.Run(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
