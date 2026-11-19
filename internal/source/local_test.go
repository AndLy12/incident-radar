package source

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

// TestOpenFileAcceptsLogAndTextFiles проверяет открытие файлов с расширениями .log и .txt.
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

// TestOpenFileRejectsUnsupportedAndMissingFiles проверяет ошибки для отсутствующего файла и неподдерживаемого расширения.
func TestOpenFileRejectsUnsupportedAndMissingFiles(t *testing.T) {
	if _, err := OpenFile("access.csv"); err == nil {
		t.Fatal("ожидалась ошибка для неподдерживаемого расширения")
	}
	if _, err := OpenFile(filepath.Join(t.TempDir(), "missing.log")); err == nil {
		t.Fatal("ожидалась ошибка для отсутствующего файла")
	}
}

// TestExtendGlob проверяет выбор файлов по маске и обработку ошибок поиска.
func TestExtendGlob(t *testing.T) {
	dir := t.TempDir()

	firstFile := filepath.Join(dir, "access-1.log")
	secondFile := filepath.Join(dir, "access-2.log")
	txtFile := filepath.Join(dir, "other.txt")

	for _, path := range []string{firstFile, secondFile, txtFile} {
		if err := os.WriteFile(path, []byte("test\n"), 0o600); err != nil {
			t.Fatalf("не удалось создать файл %s: %v", path, err)
		}
	}

	tests := []struct {
		name      string
		pattern   string
		wantFiles []string
		wantErr   bool
	}{
		{
			name:    "все log файлы без txt",
			pattern: filepath.Join(dir, "*.log"),
			wantFiles: []string{
				firstFile,
				secondFile,
			},
		},
		{
			name:      "маска выбирает один файл",
			pattern:   filepath.Join(dir, "access-1.*"),
			wantFiles: []string{firstFile},
		},
		{
			name:      "маска выбирает txt",
			pattern:   filepath.Join(dir, "*.txt"),
			wantFiles: []string{txtFile},
		},
		{
			name:      "точный путь существующего файла",
			pattern:   firstFile,
			wantFiles: []string{firstFile},
		},
		{
			name:    "нет совпадений",
			pattern: filepath.Join(dir, "missing-*.log"),
			wantErr: true,
		},
		{
			name:    "некорректная маска",
			pattern: filepath.Join(dir, "["),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, err := ExpandGlob(tt.pattern)

			if tt.wantErr {
				if err == nil {
					t.Fatal("ожидали ошибку, но получили nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			// Порядок путей не влияет на результат проверки.
			got := slices.Clone(gotFiles)
			want := slices.Clone(tt.wantFiles)

			slices.Sort(got)
			slices.Sort(want)

			if !slices.Equal(got, want) {
				t.Errorf(
					"неверный список файлов:\nполучили: %v\nожидали: %v",
					got,
					want,
				)
			}
		})
	}
}
