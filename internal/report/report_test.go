package report

import (
	"encoding/json"
	"errors"
	"incidentradar/internal/analytics"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestNewReport проверяет перенос статистики из коллектора в отчёт.
func TestNewReport(t *testing.T) {
	topResources := []analytics.ResourceRecord{
		{Resource: "/a", Count: 11},
		{Resource: "/b", Count: 10},
		{Resource: "/c", Count: 9},
		{Resource: "/d", Count: 8},
		{Resource: "/e", Count: 7},
		{Resource: "/f", Count: 6},
		{Resource: "/g", Count: 5},
		{Resource: "/h", Count: 4},
		{Resource: "/i", Count: 3},
		{Resource: "/j", Count: 2},
	}
	records := []analytics.ResourceRecord{{Resource: "/k", Count: 1}}
	records = append(records, topResources...)
	collector := analytics.Collector{
		CounterRequests: 66,
		StatusCounters:  map[int]uint64{200: 65, 404: 1},
		ResourceRecords: records,
		MaxSizeBytes:    100,
		Sizes:           []int64{100, 100},
	}

	got, err := NewReport(collector, []string{"first.log", "second.txt"})
	if err != nil {
		t.Fatalf("NewReport: %v", err)
	}

	want := Report{
		Files:               []string{"first.log", "second.txt"},
		TotalRequestsCount:  66,
		ResponseSizeInBytes: responseSizeInBytes{AverageSize: 100, MaxSize: 100, P95: 100},
		Resources:           topResources,
		ResponseCodes:       map[int]uint64{200: 65, 404: 1},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewReport = %+v, ожидалось %+v", got, want)
	}
}

// TestNewReportEmptyCollector проверяет ошибку при создании отчёта без записей.
func TestNewReportEmptyCollector(t *testing.T) {
	got, err := NewReport(*analytics.NewCollector(), nil)
	if !errors.Is(err, analytics.ErrNoResponseSizes) {
		t.Fatalf("NewReport: получено %v, ожидалось %v", err, analytics.ErrNoResponseSizes)
	}
	if !reflect.DeepEqual(got, Report{}) {
		t.Errorf("NewReport = %+v, ожидался пустой отчёт", got)
	}
}

// TestMakeJSONReport проверяет JSON-объект на диске и значения его полей.
func TestMakeJSONReport(t *testing.T) {
	t.Helper()

	report := Report{
		Files:               []string{"access.log"},
		TotalRequestsCount:  3,
		ResponseSizeInBytes: responseSizeInBytes{AverageSize: 150.5, MaxSize: 200, P95: 200},
		Resources:           []analytics.ResourceRecord{{Resource: "/health", Count: 2}, {Resource: "/missing", Count: 1}},
		ResponseCodes:       map[int]uint64{200: 2, 404: 1},
	}
	output := filepath.Join(t.TempDir(), "report.json")
	if err := report.MakeJSONReport(output); err != nil {
		t.Fatalf("MakeJSONReport: %v", err)
	}
	data, err := os.ReadFile(output) // #nosec G304 -- Путь задан тестом внутри t.TempDir.
	if err != nil {
		t.Fatalf("чтение отчёта: %v", err)
	}

	var got map[string]any
	if err = json.Unmarshal(data, &got); err != nil {
		t.Fatalf("ожидался JSON-объект: %v", err)
	}
	want := map[string]any{
		"files":              []any{"access.log"},
		"totalRequestsCount": float64(3),
		"responseSizeInBytes": map[string]any{
			"average": 150.5,
			"max":     float64(200),
			"p95":     float64(200),
		},
		"resources": []any{
			map[string]any{"Resource": "/health", "Count": float64(2)},
			map[string]any{"Resource": "/missing", "Count": float64(1)},
		},
		"responseCodes": map[string]any{"200": float64(2), "404": float64(1)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("JSON = %v, ожидалось %v", got, want)
	}
}

// TestMakeMDReport проверяет содержимое Markdown-отчёта.
func TestMakeMDReport(t *testing.T) {
	t.Helper()

	report := Report{
		Files:               []string{"access.log", "errors.log"},
		TotalRequestsCount:  3,
		ResponseSizeInBytes: responseSizeInBytes{AverageSize: 150.5, MaxSize: 200, P95: 190},
		Resources:           []analytics.ResourceRecord{{Resource: "/health", Count: 2}, {Resource: "/missing", Count: 1}},
		ResponseCodes:       map[int]uint64{200: 2, 404: 1},
	}
	output := filepath.Join(t.TempDir(), "report.md")
	if err := report.MakeMDReport(output); err != nil {
		t.Fatalf("MakeMDReport: %v", err)
	}
	data, err := os.ReadFile(output) // #nosec G304 -- Путь задан тестом внутри t.TempDir.
	if err != nil {
		t.Fatalf("чтение отчёта: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Fatal("Markdown-отчёт пуст")
	}
	for _, line := range []string{
		"# Отчёт по анализу NGINX-логов",
		"## Источники",
		"| access.log |",
		"| errors.log |",
		"## Общая статистика",
		"| Общее количество запросов | 3 |",
		"| Средний размер ответа, байт | 150.50 |",
		"| Максимальный размер ответа, байт | 200 |",
		"| 95-й перцентиль размера ответа, байт | 190 |",
		"## Топ ресурсов",
		"| /health | 2 |",
		"| /missing | 1 |",
		"## Коды ответа",
		"| 200 | 2 |",
		"| 404 | 1 |",
	} {
		if !strings.Contains(content, line) {
			t.Errorf("в Markdown-отчёте отсутствует %q", line)
		}
	}
}

// TestReportPreservesExistingFile проверяет отказ от перезаписи в обоих форматах.
func TestReportPreservesExistingFile(t *testing.T) {
	for _, test := range []struct {
		extension string
		write     func(Report, string) error
	}{
		{extension: "json", write: Report.MakeJSONReport},
		{extension: "md", write: Report.MakeMDReport},
	} {
		t.Run(test.extension, func(t *testing.T) {
			output := filepath.Join(t.TempDir(), "report."+test.extension)
			const original = "existing report"
			if err := os.WriteFile(output, []byte(original), 0o600); err != nil {
				t.Fatalf("создание файла: %v", err)
			}

			if err := test.write(Report{}, output); !errors.Is(err, os.ErrExist) {
				t.Errorf("запись отчёта: %v, ожидалась ошибка os.ErrExist", err)
			}
			data, err := os.ReadFile(output) // #nosec G304 -- Путь задан тестом внутри t.TempDir.
			if err != nil {
				t.Fatalf("чтение существующего файла: %v", err)
			}
			if string(data) != original {
				t.Errorf("содержимое файла = %q, ожидалось %q", data, original)
			}
		})
	}
}

// TestMakeJSONReportMarshalError проверяет ошибку сериализации до создания файла.
func TestMakeJSONReportMarshalError(t *testing.T) {
	report := Report{ResponseSizeInBytes: responseSizeInBytes{AverageSize: math.NaN()}}
	output := filepath.Join(t.TempDir(), "report.json")

	var unsupportedValue *json.UnsupportedValueError
	if err := report.MakeJSONReport(output); !errors.As(err, &unsupportedValue) {
		t.Errorf("MakeJSONReport = %v, ожидалась ошибка json.UnsupportedValueError", err)
	}
	if _, err := os.Stat(output); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("os.Stat = %v, ожидалось отсутствие файла", err)
	}
}
