package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunReadsMatchingFiles(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "access.log")
	outputPath := filepath.Join(dir, "report.md")
	line := `93.180.71.3 - - [17/May/2015:08:05:32 +0000] "GET /downloads HTTP/1.1" 200 10 "-" "agent"`
	if err := os.WriteFile(logPath, []byte(line+"\n"), 0o600); err != nil {
		t.Fatalf("не удалось создать лог-файл: %v", err)
	}
	if err := Run([]string{"--path", logPath, "--format", "markdown", "--output", outputPath}); err != nil {
		t.Fatalf("Run() завершился с ошибкой: %v", err)
	}
}

func TestRunReturnsArgumentError(t *testing.T) {
	if err := Run([]string{"--unknown"}); err == nil {
		t.Fatal("ожидалась ошибка неизвестного аргумента")
	}
}
