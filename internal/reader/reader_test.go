package reader

import (
	"incidentradar/internal/analytics"
	"path/filepath"
	"testing"
)

func TestReadAndParseFileReturnsOpenError(t *testing.T) {
	_, err := ReadAndParseFile(analytics.Filter{}, filepath.Join(t.TempDir(), "missing.log"))
	if err == nil {
		t.Fatal("ожидалась ошибка открытия файла")
	}
}
