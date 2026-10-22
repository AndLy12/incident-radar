package config

import (
	"flag"
)

func Load(args []string) (*Config, error) {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)

	var cfg Config

	fs.StringVar(&cfg.path, "path", "", "Path to the config file")
	fs.StringVar(&cfg.path, "p", "", "Path to the config file")

	fs.StringVar(&cfg.outputFormat, "output-format", "json", "Output format")
	fs.StringVar(&cfg.outputFormat, "f", "json", "Output format")

	fs.StringVar(&cfg.outputPath, "output-path", "output.json", "Output path")
	fs.StringVar(&cfg.outputPath, "o", "output.json", "Output path")

	fs.StringVar(&cfg.fromDate, "from", "", "From date")
	fs.StringVar(&cfg.toDate, "to", "", "To date")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if err := validate(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
