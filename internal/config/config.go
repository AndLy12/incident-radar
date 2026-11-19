// Package config хранит и проверяет параметры обработки логов.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var outputFormatExtensions = map[string]string{
	"json":     ".json",
	"markdown": ".md",
}

// LogsFormatExtensions содержит поддерживаемые расширения файлов логов.
var LogsFormatExtensions = []string{".log", ".txt"}

// Config хранит параметры источника логов, отчёта и фильтрации по датам.
type Config struct {
	logsPath         string
	format           string
	output           string
	layoutForLogLine regexp.Regexp
	from             *time.Time
	to               *time.Time
}

// LayoutForLogLine возвращает регулярное выражение для строки лога.
func (c Config) LayoutForLogLine() regexp.Regexp {
	return c.layoutForLogLine
}

// LogsPath возвращает путь или маску для поиска логов.
func (c Config) LogsPath() string {
	return c.logsPath
}

// Format возвращает формат отчёта.
func (c Config) Format() string {
	return c.format
}

// Output возвращает путь к файлу отчёта.
func (c Config) Output() string {
	return c.output
}

// From возвращает нижнюю границу периода и признак её наличия.
func (c Config) From() (time.Time, bool) {
	if c.from == nil {
		return time.Time{}, false
	}
	return *c.from, true
}

// To возвращает верхнюю границу периода и признак её наличия.
func (c Config) To() (time.Time, bool) {
	if c.to == nil {
		return time.Time{}, false
	}
	return *c.to, true
}

// validate проверяет путь к логам, параметры отчёта и порядок дат.
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

	outputFileExtension := filepath.Ext(c.output)
	if outputFileExtension != outputFormatExtensions[c.format] {
		return fmt.Errorf("output file extension does not match format")
	}
	_, err := os.Stat(c.output)
	if err == nil {
		return fmt.Errorf("output file already exists")
	} else if !os.IsNotExist(err) {
		return err
	}

	if c.from != nil && c.to != nil && (c.to.Before(*c.from) || c.to.Equal(*c.from)) {
		return fmt.Errorf("from must be < to")
	}

	return nil
}

// NewConfig создаёт и проверяет конфигурацию обработки логов.
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

// cloneDate копирует значение даты, сохраняя nil для отсутствующей даты.
func cloneDate(date *time.Time) *time.Time {
	if date == nil {
		return nil
	}
	copied := *date
	return &copied
}
