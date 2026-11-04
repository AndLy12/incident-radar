package cli

import (
	"flag"
	"fmt"
	"incidentradar/internal/config"
	"io"
	"time"
)

func Parse(Args []string, diagnostics io.Writer) (config.Config, error) {
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

	if err := fs.Parse(Args); err != nil {
		return config.Config{}, err
	}
	if fs.NArg() != 0 {
		return config.Config{}, fmt.Errorf("unexpected positional arguments: %s", fs.Arg(0))
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

func parseDate(dateStr string, has bool) (*time.Time, error) {
	if !has {
		return nil, nil
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date %s, expended format YYYY-MM-DD", dateStr)
	}
	return &date, nil
}
