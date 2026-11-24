package main

import (
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/reader"
	"incidentradar/internal/source"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// TestGlobProcessingIntegration сравнивает статистику для отдельных логов по маске и одного объединённого лога.
func TestGlobProcessingIntegration(t *testing.T) {
	t.Helper()
	dir := t.TempDir()

	firstLine := `192.168.1.1 - - [17/May/2015:08:00:00 +0000] "GET /products HTTP/1.1" 200 100 "-" "Test"`
	secondLine := `192.168.1.2 - - [18/May/2015:09:00:00 +0000] "GET /orders HTTP/1.1" 404 200 "-" "Test"`

	files := map[string]string{
		"first.log":    firstLine + "\n",
		"second.log":   secondLine + "\n",
		"combined.txt": firstLine + "\n" + secondLine + "\n",
	}

	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Error(err)
		}
	}

	// Маска должна захватить только два .log-файла.
	paths, err := source.ExpandGlob(filepath.Join(dir, "*.log"))
	if err != nil {
		t.Error(err)
	}

	wantPaths := []string{
		filepath.Join(dir, "first.log"),
		filepath.Join(dir, "second.log"),
	}

	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("ожидали файлы %v, получили %v", wantPaths, paths)
	}

	filter := analytics.NewFilter(nil, nil)

	// Получаем результат обработки двух найденных файлов.
	got := analytics.NewCollector()
	for _, path := range paths {
		if err := reader.ReadAndParseFile(*filter, path, got); err != nil {
			t.Error(err)
		}
	}

	// Получаем результат обработки тех же записей одним файлом.
	want := analytics.NewCollector()
	combinedPath := filepath.Join(dir, "combined.txt")
	if err := reader.ReadAndParseFile(*filter, combinedPath, want); err != nil {
		t.Error(err)
	}

	// Дополнительно убеждаемся, что записи попали в коллектор.
	if reflect.DeepEqual(want, analytics.NewCollector()) {
		t.Fatal("после обработки коллектор остался пустым")
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"результаты обработки нескольких файлов и одного файла различаются:\nполучили: %+v\nожидали: %+v",
			got,
			want,
		)
	}
}

// TestCases проверяет, что обработка пропускает некорректные строки и применяет фильтр по датам.
func TestCases(t *testing.T) {
	t.Helper()

	dir := t.TempDir()

	gotFile := filepath.Join(dir, "got.log")
	wantFile := filepath.Join(dir, "want.log")

	dateFrom := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(2025, 3, 31, 0, 0, 0, 0, time.UTC)

	err := os.WriteFile(gotFile, []byte("192.0.2.1 - - [01/01/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"\n"+
		"192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"\n"+
		"this line is intentionally malformed\n"+
		"192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"\n"+
		"192.0.2.2 - - [01/12/2025:10:01:00 +0000] \"GET /after-error HTTP/1.1\" 500 500 \"-\" \"curl/8.0\"\n"), 0o600)
	if err != nil {
		t.Error(err)
	}

	err = os.WriteFile(wantFile, []byte("192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"\n"+
		"192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"\n"), 0o600)
	if err != nil {
		t.Error(err)
	}

	got := analytics.NewCollector()
	filter := analytics.NewFilter(&dateFrom, &dateTo)
	err = reader.ReadAndParseFile(*filter, gotFile, got)
	if err != nil {
		t.Error(err)
	}

	want := analytics.NewCollector()
	err = reader.ReadAndParseFile(*filter, wantFile, want)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf(
			"результаты обработки нескольких файлов и одного файла различаются:\nполучили: %+v\nожидали: %+v",
			got,
			want,
		)
	}
}

// TestHTTP проверяет обработку лога с тестового HTTP-сервера.
func TestHTTP(t *testing.T) {
	t.Helper()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprintln(w, "192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"")
		if err != nil {
			t.Error(err)
		}
	}))
	defer testServer.Close()

	filter := analytics.NewFilter(nil, nil)
	testCollector := analytics.NewCollector()

	err := reader.ReadAndParseFile(*filter, testServer.URL+"/access.log", testCollector)
	if err != nil {
		t.Error(err)
	}
}
