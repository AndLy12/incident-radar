package cli

import (
	"os"
	"testing"
)

// TestCase описывает входные данные и ожидаемый признак ошибки для теста.
type TestCase struct {
	name string
	args []string
	err  bool
}

// TestParse проверяет обработку некорректных аргументов командной строки.
func TestParse(t *testing.T) {
	t.Helper()

	tests := []TestCase{
		{
			name: "неизвестный флаг",
			args: []string{"--unknown", "value"},
			err:  true,
		},
		{
			name: "опечатка в названии флага",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--formatt", "json",
			},
			err: true,
		},
		{
			name: "path без значения",
			args: []string{"--path"},
			err:  true,
		},
		{
			name: "output без значения",
			args: []string{"--path", "access.log", "--output"},
			err:  true,
		},
		{
			name: "format без значения",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--format",
			},
			err: true,
		},
		{
			name: "from без значения",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--from",
			},
			err: true,
		},
		{
			name: "to без значения",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--to",
			},
			err: true,
		},
		{
			name: "текст вместо даты",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--from", "abc",
			},
			err: true,
		},
		{
			name: "неверный формат даты",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--from", "15.11.2025",
			},
			err: true,
		},
		{
			name: "несуществующий месяц",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--to", "2025-13-01",
			},
			err: true,
		},
		{
			name: "несуществующий день",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--from", "2025-02-30",
			},
			err: true,
		},
		{
			name: "дата со временем",
			args: []string{
				"--path", "access.log",
				"--output", "report.md",
				"--to", "2025-11-15T12:00:00",
			},
			err: true,
		},
		{
			name: "не указан источник",
			args: []string{"--output", "report.md"},
			err:  true,
		},
		{
			name: "не указан output",
			args: []string{"--path", "access.log"},
			err:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.args, os.Stderr)
			if test.err && err == nil {
				t.Fatal("ожидалась ошибка, но конфиг собрался")
			}

			if !test.err && err != nil {
				t.Fatalf("конфигурация верная, но получили ошибку: %v", err)
			}
		})
	}
}
