package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type configTestCase struct {
	name        string
	config      Config
	extendError bool
}

// TestConfigValidate проверяет допустимые параметры отчёта и ошибки конфигурации.
func TestConfigValidate(t *testing.T) {
	t.Helper()

	// Файловое окружение специально для тестов
	dir := t.TempDir()

	logFile := filepath.Join(dir, "access.log")
	txtFile := filepath.Join(dir, "access.txt")
	unsupportedFile := filepath.Join(dir, "unsupported.csv")

	existingOutput := filepath.Join(dir, "existing.md")

	for _, path := range []string{
		logFile,
		txtFile,
		unsupportedFile,
		existingOutput,
	} {
		if err := os.WriteFile(path, []byte("test"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	validConfig := func() Config {
		return Config{
			logsPath: logFile,
			format:   "markdown",
			output:   filepath.Join(dir, "report.md"),
			from:     nil,
			to:       nil,
		}
	}

	change := func(edit func(*Config)) Config {
		cfg := validConfig()
		edit(&cfg)
		return cfg
	}

	TestCases := []configTestCase{
		{
			name:        "valid",
			config:      validConfig(),
			extendError: false,
		},
		{
			name: "источник txt",
			config: change(func(c *Config) {
				c.logsPath = txtFile
			}),
			extendError: false,
		},
		{
			name: "корректная конфигурация json",
			config: change(func(c *Config) {
				c.format = "json"
				c.output = filepath.Join(dir, "report.json")
			}),
			extendError: false,
		},
		{
			name: "только нижняя граница периода",
			config: change(func(c *Config) {
				date := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.from = &date
			}),
			extendError: false,
		},
		{
			name: "только верхняя граница периода",
			config: change(func(c *Config) {
				date := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
				c.to = &date
			}),
			extendError: false,
		},
		{
			name: "обе границы периода",
			config: change(func(c *Config) {
				dateFrom := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.from = &dateFrom
				dateTo := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
				c.to = &dateTo
			}),
			extendError: false,
		},
		{
			name: "соседние даты",
			config: change(func(c *Config) {
				dateFrom := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.from = &dateFrom
				dateTo := time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC)
				c.to = &dateTo
			}),
			extendError: false,
		},
		{
			name: "период через границу года",
			config: change(func(c *Config) {
				dateFrom := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
				c.from = &dateFrom
				dateTo := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.to = &dateTo
			}),
			extendError: false,
		},
		{
			name: "пустой путь источника",
			config: change(func(c *Config) {
				c.logsPath = ""
			}),
			extendError: true,
		},
		{
			name: "неподдерживаемый формат",
			config: change(func(c *Config) {
				c.format = "xml"
			}),
			extendError: true,
		},
		{
			name: "пустой output",
			config: change(func(c *Config) {
				c.output = ""
			}),
			extendError: true,
		},
		{
			name: "json с расширением md",
			config: change(func(c *Config) {
				c.format = "json"
				c.output = filepath.Join(dir, "report.md")
			}),
			extendError: true,
		},
		{
			name: "markdown с расширением json",
			config: change(func(c *Config) {
				c.format = "markdown"
				c.output = filepath.Join(dir, "report.json")
			}),
			extendError: true,
		},
		{
			name: "output без расширения",
			config: change(func(c *Config) {
				c.output = filepath.Join(dir, "report")
			}),
			extendError: true,
		},
		{
			name: "output с неподдерживаемым расширением",
			config: change(func(c *Config) {
				c.output = filepath.Join(dir, "report.csv")
			}),
			extendError: true,
		},
		{
			name: "одинаковые даты",
			config: change(func(c *Config) {
				dateFrom := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.from = &dateFrom
				dateTo := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.to = &dateTo
			}),
			extendError: true,
		},
		{
			name: "начальная дата позже конечной",
			config: change(func(c *Config) {
				dateFrom := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)
				c.from = &dateFrom
				dateTo := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
				c.to = &dateTo
			}),
			extendError: true,
		},
		{
			name: "output уже существует",
			config: change(func(c *Config) {
				c.output = existingOutput
			}),
			extendError: true,
		},
	}

	for _, testCase := range TestCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.config.validate()

			if testCase.extendError && err == nil {
				t.Fatal("ожидали ошибку, но конфигурация принята")
			}

			if err != nil && !testCase.extendError {
				t.Fatalf("конфигурация корректная, но получили ошибку: %v", err)
			}
		})
	}
}
