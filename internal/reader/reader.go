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

// ReadAndParseFile читает лог с диска или по HTTP и передаёт подходящие записи сборщику.
func ReadAndParseFile(filter analytics.Filter, file string, collector *analytics.Collector) error {
	var resource io.ReadCloser
	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, file, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		resource = resp.Body
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				log.Println(err)
			}
		}()
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
