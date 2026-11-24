// Package cli разбирает аргументы командной строки Incident Radar.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"incidentradar/internal/config"
	"io"
	"time"
)

// Parse разбирает аргументы командной строки и возвращает проверенную конфигурацию.
func Parse(args []string, diagnostics io.Writer) (config.Config, error) {
	var path, format, output string
	var toStr, fromStr string

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(diagnostics)

	fs.StringVar(&path, "path", "", "Path to file or directory")
	fs.StringVar(&path, "p", "", "Path to file or directory")

	fs.StringVar(&format, "format", "json", "Output format")
	fs.StringVar(&format, "f", "", "Output format")

	fs.StringVar(&output, "output", "", "Output file or directory")
	fs.StringVar(&output, "o", "", "Output file or directory")

	fs.StringVar(&fromStr, "from", "", "From file or directory")
	fs.StringVar(&toStr, "to", "", "To file or directory")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return config.Config{}, err
		}
		return config.Config{}, fmt.Errorf("%w: %w", ErrInvalidFlags, err)
	}
	if fs.NArg() != 0 {
		return config.Config{}, fmt.Errorf("%w: %s", ErrUnexpectedPositionalArgument, fs.Arg(0))
	}

	var hasFrom, hasTo bool
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "from":
			hasFrom = true
		case "to":
			hasTo = true
		}
	})

	from, err := parseDate(fromStr, hasFrom)
	if err != nil {
		return config.Config{}, err
	}

	to, err := parseDate(toStr, hasTo)
	if err != nil {
		return config.Config{}, err
	}

	return config.NewConfig(path, format, output, from, to)
}

// parseDate разбирает указанную дату в формате YYYY-MM-DD.
func parseDate(dateStr string, has bool) (*time.Time, error) {
	if !has {
		return nil, nil
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("%w %s, expended format YYYY-MM-DD", ErrInvalidDate, dateStr)
	}
	return &date, nil
}
