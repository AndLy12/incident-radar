package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"time"
)

var logsFormatExtensions []string = []string{".md", ".txt"}

type Config struct {
	logsPath         string
	format           string
	output           string
	layoutForLogLine regexp.Regexp
	from             *time.Time
	to               *time.Time
}

func (c Config) LayoutForLogLine() regexp.Regexp {
	return c.layoutForLogLine
}

func (c Config) LogsPath() string {
	return c.logsPath
}

func (c Config) Format() string {
	return c.format
}

func (c Config) Output() string {
	return c.output
}

func (c Config) From() (time.Time, bool) {
	if c.from == nil {
		return time.Time{}, false
	}
	return *c.from, true
}

func (c Config) To() (time.Time, bool) {
	if c.to == nil {
		return time.Time{}, false
	}
	return *c.to, true
}

func (c Config) validate() error {
	if c.logsPath == "" {
		return fmt.Errorf("logsPath required")
	}
	if c.format != "json" && c.format != "markdown" {
		return fmt.Errorf("format must be json or md")
	}
	if c.output == "" {
		return fmt.Errorf("output must be set")
	}
	if c.from != nil && c.to != nil && c.from.Before(*c.to) {
		return fmt.Errorf("from must be <= to")
	}

	logsFileExtension := filepath.Ext(c.logsPath)
	var flagExtNotInLogsFormatExtensions bool = true
	for _, ext := range logsFormatExtensions {
		if logsFileExtension == ext {
			flagExtNotInLogsFormatExtensions = false
		}
	}
	if flagExtNotInLogsFormatExtensions {
		return fmt.Errorf("file format not supported")
	}
	return nil
}

func NewConfig(path, format, output string, from, to *time.Time) (Config, error) {
	config := Config{
		logsPath: path,
		format:   format,
		output:   output,
		layoutForLogLine: *regexp.MustCompile(
			"^(\\S+)\\s-\\s" + // ip
				"(\\S+)\\s\\" + // remote user
				"[([^]]+)\\]\\s" + // time
				"\"([A-Z]+)\\s" + // method
				"([^\"\\s]+)\\s" + // resource
				"([^\"\\s]+)\"\\s" + // protocol
				"([0-9]{3})\\s" + // status
				"([0-9]+)\\s" + // size
				"\"([^\"]*)\"\\s" + // http_referer
				"\"([^\"]*)\"$", // http_user_agent
		),
		from: cloneDate(from),
		to:   cloneDate(to),
	}
	if err := config.validate(); err != nil {
		return Config{}, errors.New("invalid config")
	}
	return config, nil
}

func cloneDate(date *time.Time) *time.Time {
	if date == nil {
		return nil
	}
	copied := *date
	return &copied
}
