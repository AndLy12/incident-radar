package reader

import (
	"fmt"
	"incidentradar/internal/analytics"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

// TestReadAndParseFileReturnsOpenError проверяет ошибку при открытии отсутствующего лога.
func TestReadAndParseFileReturnsOpenError(t *testing.T) {
	err := ReadAndParseFile(analytics.Filter{}, filepath.Join(t.TempDir(), "missing.log"), analytics.NewCollector())
	if err == nil {
		t.Fatal("ожидалась ошибка открытия файла")
	}
}

// TestReadAndParseHTTP проверяет чтение лога по HTTP.
func TestReadAndParseHTTP(t *testing.T) {
	t.Helper()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprintln(w, "192.0.2.1 - - [01/Mar/2025:10:00:00 +0000] \"GET /before-error HTTP/1.1\" 200 100 \"-\" \"curl/8.0\"")
		if err != nil {
			t.Error(err)
		}
	}))
	defer testServer.Close()

	err := ReadAndParseFile(analytics.Filter{}, testServer.URL+"/access.txt", analytics.NewCollector())
	if err != nil {
		t.Fatal(err)
	}
}
