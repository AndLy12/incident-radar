// Package report формирует отчёты по статистике журнала доступа NGINX.
package report

import (
	"encoding/json"
	"fmt"
	"incidentradar/internal/analytics"
	"log"
	"os"
	"strings"
)

// Title содержит заголовок Markdown-отчёта и начало таблицы источников.
const Title string = `# Отчёт по анализу NGINX-логов

## Источники

| Файл или URL |
|---|`

// files задаёт шаблон строки таблицы источников.
const files string = `
| %s |`

// generalStatistic задаёт шаблон таблицы общей статистики.
const generalStatistic string = `
## Общая статистика

| Показатель | Значение |
|---|---:|
| Общее количество запросов | %d |
| Средний размер ответа, байт | %.2f |
| Максимальный размер ответа, байт | %d |
| 95-й перцентиль размера ответа, байт | %d |`

// topResources задаёт шаблон строки таблицы популярных ресурсов.
const topResources = `
| %s | %d |`

// Codes задаёт шаблон строки таблицы кодов ответа.
const Codes = `
| %d | %d |`

// Report хранит источники и сводную статистику запросов.
type Report struct {
	Files               []string                   `json:"files"`
	TotalRequestsCount  uint64                     `json:"totalRequestsCount"`
	ResponseSizeInBytes responseSizeInBytes        `json:"responseSizeInBytes"`
	Resources           []analytics.ResourceRecord `json:"resources"`
	ResponseCodes       map[int]uint64             `json:"responseCodes"`
}

// responseSizeInBytes хранит средний и максимальный размер ответа и 95-й перцентиль.
type responseSizeInBytes struct {
	AverageSize float64 `json:"average"`
	MaxSize     int64   `json:"max"`
	P95         int64   `json:"p95"`
}

// NewReport формирует отчёт по данным коллектора и списку источников.
func NewReport(collector analytics.Collector, files []string) (Report, error) {
	avg, err := collector.AVGResourceSize()
	if err != nil {
		return Report{}, err
	}
	p95, err := collector.P95ResourcesSize()
	if err != nil {
		return Report{}, err
	}
	top10Resources, err := collector.Top10resources()
	if err != nil {
		return Report{}, err
	}
	return Report{
		Files:              files,
		TotalRequestsCount: collector.CounterRequests,
		ResponseSizeInBytes: responseSizeInBytes{
			AverageSize: avg,
			MaxSize:     collector.MaxSizeBytes,
			P95:         p95,
		},
		Resources:     top10Resources,
		ResponseCodes: collector.StatusCounters,
	}, nil
}

// MakeJSONReport сохраняет отчёт в новый JSON-файл по указанному пути.
func (report Report) MakeJSONReport(outputFileString string) error {
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	outputFile, err := os.OpenFile(outputFileString, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- The local CLI user explicitly selects the report path via --outputStringFile; O_EXCL prevents overwriting an existing file.
	if err != nil {
		return err
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	_, err = outputFile.Write(reportJSON)
	if err != nil {
		return err
	}
	return nil
}

// MakeMDReport сохраняет отчёт в новый Markdown-файл по указанному пути.
func (report Report) MakeMDReport(output string) error {
	var result strings.Builder
	_, err := fmt.Fprint(&result, Title)
	if err != nil {
		return err
	}
	for _, file := range report.Files {
		_, err = fmt.Fprintf(&result, files, file)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(&result, generalStatistic,
		report.TotalRequestsCount,
		report.ResponseSizeInBytes.AverageSize,
		report.ResponseSizeInBytes.MaxSize,
		report.ResponseSizeInBytes.P95,
	)
	if err != nil {
		return err
	}

	result.WriteString("\n## Топ ресурсов\n\n| Ресурс | Количество запросов |\n|---|---:|")
	for _, res := range report.Resources {
		_, err = fmt.Fprintf(&result, topResources, res.Resource, res.Count)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(&result, "## Коды ответа\n\n| Код ответа | Количество запросов |\n|---|---:|")
	if err != nil {
		return err
	}
	for code, count := range report.ResponseCodes {
		_, err = fmt.Fprintf(&result, Codes, code, count)
		if err != nil {
			return err
		}
	}

	outputFile, err := os.OpenFile(output, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- The local CLI user explicitly selects the report path via --output; O_EXCL prevents overwriting an existing file.
	if err != nil {
		return err
	}
	defer func() {
		err = outputFile.Close()
		if err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	_, err = outputFile.Write([]byte(result.String()))
	if err != nil {
		return err
	}
	return nil
}
