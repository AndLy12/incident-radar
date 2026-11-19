// Package main implements the Incident Radar CLI application.
package main

import (
	"incidentradar/internal/analytics"
	"incidentradar/internal/cli"
	"incidentradar/internal/reader"
	"incidentradar/internal/source"
	"log"
	"os"
	"time"
)

// main разбирает аргументы и собирает статистику по найденным логам.
func main() {
	args := os.Args[1:]
	// Создаем новый конфиг, распарсив команду и ее флаги
	config, err := cli.Parse(args, os.Stderr)
	if err != nil {
		log.Fatal(err)
	}

	filePaths, err := source.ExpandGlob(config.LogsPath())
	if err != nil {
		log.Fatal(err)
	}

	// Подготавливаем фильтрацию по дате
	from, has := config.From()
	var fromForFilter *time.Time
	if has {
		fromForFilter = &from
	}
	to, has := config.To()
	var toForFilter *time.Time
	if has {
		toForFilter = &to
	}
	filter := analytics.NewFilter(fromForFilter, toForFilter)

	// Проходим по всем файлам
	collector := analytics.NewCollector()
	for _, filePath := range filePaths {
		err := reader.ReadAndParseFile(*filter, filePath, collector)
		if err != nil {
			log.Fatal(err)
		}
	}
}
