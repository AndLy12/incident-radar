package reader

import (
	"bufio"
	"errors"
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/parser"
	"incidentradar/internal/source"
	"io"
	"log"
)

func ReadAndParseFile(filter analytics.Filter, file string) ([]parser.LogRecord, error) {
	f, err := source.OpenFile(file)
	if err != nil {
		return []parser.LogRecord{}, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	logs := make([]parser.LogRecord, 0)
	// Читаем логи из файла
	scanner := bufio.NewScanner(f)
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
			logs = append(logs, lr)
			fmt.Println(lr)
		}
	}
	if err := scanner.Err(); !errors.Is(err, io.EOF) {
		return []parser.LogRecord{}, err
	}
	return logs, nil
}
