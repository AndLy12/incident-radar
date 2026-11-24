// Package reader читает и обрабатывает записи из источников логов.
package reader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/parser"
	"incidentradar/internal/source"
	"io"
	"log"
	"net/http"
	"strings"
)

// IsHTTPSource сообщает, что источник нужно читать по HTTP или HTTPS.
func IsHTTPSource(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// ReadAndParseFile читает лог с диска или по HTTP и передаёт подходящие записи сборщику.
func ReadAndParseFile(filter analytics.Filter, file string, collector *analytics.Collector) error {
	var resource io.ReadCloser
	if IsHTTPSource(file) {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, file, nil)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrHTTPSourceUnavailable, err)
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrHTTPSourceUnavailable, err)
		}
		resource = resp.Body
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return fmt.Errorf("%w: %s returned %s", ErrHTTPStatus, file, resp.Status)
		}
	} else {
		var err error
		resource, err = source.OpenFile(file)
		if err != nil {
			return err
		}
		defer func() {
			if err = resource.Close(); err != nil {
				log.Println(err)
			}
		}()
	}

	// Читаем логи из файла
	scanner := bufio.NewScanner(resource)
	var lineNumber int
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Парсинг строки логов
		lr, err := parser.SplitLine(line)
		if err != nil {
			log.Printf("%s:%d:%s", file, lineNumber, err.Error())
			continue
		}

		// Фильтруем по дате
		if filter.Match(lr.Time) {
			collector.Add(lr)
			fmt.Println(lr)
		}
	}
	if err := scanner.Err(); !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}
