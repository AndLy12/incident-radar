package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFileAcceptsLogAndTextFiles(t *testing.T) {
	dir := t.TempDir()
	for _, extension := range []string{".log", ".txt"} {
		path := filepath.Join(dir, "access"+extension)
		if err := os.WriteFile(path, []byte("line"), 0o600); err != nil {
			t.Fatalf("не удалось создать файл: %v", err)
		}
		file, err := OpenFile(path)
		if err != nil {
			t.Fatalf("OpenFile(%q) завершился с ошибкой: %v", path, err)
		}
		if err := file.Close(); err != nil {
			t.Fatalf("не удалось закрыть файл: %v", err)
		}
	}
}

func TestOpenFileRejectsUnsupportedAndMissingFiles(t *testing.T) {
	if _, err := OpenFile("access.csv"); err == nil {
		t.Fatal("ожидалась ошибка для неподдерживаемого расширения")
	}
	if _, err := OpenFile(filepath.Join(t.TempDir(), "missing.log")); err == nil {
		t.Fatal("ожидалась ошибка для отсутствующего файла")
	}
}
