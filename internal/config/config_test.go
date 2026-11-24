package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestNewConfig проверяет параметры, которые конструктор возвращает вызывающему коду.
func TestNewConfig(t *testing.T) {
	const path = "logs/*.log"
	output := filepath.Join(t.TempDir(), "report.json")

	t.Run("формат по умолчанию и отсутствие дат", func(t *testing.T) {
		cfg, err := NewConfig(path, "", output, nil, nil)
		if err != nil {
			t.Fatalf("NewConfig: %v", err)
		}
		if got := cfg.Format(); got != "json" {
			t.Errorf("Format() = %q, ожидалось json", got)
		}
		if got := cfg.LogsPath(); got != path {
			t.Errorf("LogsPath() = %q, ожидалось %q", got, path)
		}
		if got := cfg.Output(); got != output {
			t.Errorf("Output() = %q, ожидалось %q", got, output)
		}
		if got, ok := cfg.From(); ok || !got.IsZero() {
			t.Errorf("From() = (%v, %t), ожидалась отсутствующая дата", got, ok)
		}
		if got, ok := cfg.To(); ok || !got.IsZero() {
			t.Errorf("To() = (%v, %t), ожидалась отсутствующая дата", got, ok)
		}
		layout := cfg.LayoutForLogLine()
		const line = `192.0.2.1 - - [17/May/2015:08:05:32 +0000] "GET /health HTTP/1.1" 200 100 "-" "test"`
		if !layout.MatchString(line) {
			t.Error("шаблон не распознал корректную строку NGINX")
		}
	})

	t.Run("копирование границ периода", func(t *testing.T) {
		wantFrom := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
		wantTo := wantFrom.AddDate(0, 1, 0)
		from, to := wantFrom, wantTo
		cfg, err := NewConfig(path, "json", output, &from, &to)
		if err != nil {
			t.Fatalf("NewConfig: %v", err)
		}

		from = from.AddDate(0, 0, 1)
		to = to.AddDate(0, 0, 1)
		if got, ok := cfg.From(); !ok || !got.Equal(wantFrom) {
			t.Errorf("From() = (%v, %t), ожидалось (%v, true)", got, ok, wantFrom)
		}
		if got, ok := cfg.To(); !ok || !got.Equal(wantTo) {
			t.Errorf("To() = (%v, %t), ожидалось (%v, true)", got, ok, wantTo)
		}
	})

	t.Run("причина ошибки валидации", func(t *testing.T) {
		_, err := NewConfig(path, "markdown", output, nil, nil)
		if !errors.Is(err, ErrInvalidConfig) || !errors.Is(err, ErrOutputExtensionMismatch) {
			t.Errorf("NewConfig: получено %v, ожидалась ошибка конфигурации с причиной %v", err, ErrOutputExtensionMismatch)
		}
	})
}

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
