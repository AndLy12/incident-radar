// Package main implements the Incident Radar CLI application.
package main

import (
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/cli"
	"incidentradar/internal/exitcode"
	"incidentradar/internal/reader"
	"incidentradar/internal/report"
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
		_, printErr := fmt.Fprintln(os.Stderr, err)
		if printErr != nil {
			log.Println(printErr)
		}
		os.Exit(exitcode.FromError(err))
	}

	logsPath := config.LogsPath()
	logSources := []string{logsPath}
	if !reader.IsHTTPSource(logsPath) {
		logSources, err = source.ExpandGlob(logsPath)
		if err != nil {
			_, printErr := fmt.Fprintln(os.Stderr, err)
			if printErr != nil {
				log.Println(printErr)
			}
			os.Exit(exitcode.FromError(err))
		}
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

	// Читаем логи из всех источников.
	collector := analytics.NewCollector()
	for _, logSource := range logSources {
		err = reader.ReadAndParseFile(*filter, logSource, collector)
		if err != nil {
			_, printErr := fmt.Fprintln(os.Stderr, err)
			if printErr != nil {
				log.Println(printErr)
			}
			os.Exit(exitcode.FromError(err))
		}
	}

	newReport, err := report.NewReport(*collector, logSources)
	if err != nil {
		_, printErr := fmt.Fprintln(os.Stderr, err)
		if printErr != nil {
			log.Println(printErr)
		}
		os.Exit(exitcode.FromError(err))
	}

	if config.Format() == "json" {
		if err = newReport.MakeJSONReport(config.Output()); err != nil {
			_, printErr := fmt.Fprintln(os.Stderr, err)
			if printErr != nil {
				log.Println(printErr)
			}
			os.Exit(exitcode.FromError(err))
		}
	} else if config.Format() == "markdown" {
		if err = newReport.MakeMDReport(config.Output()); err != nil {
			_, printErr := fmt.Fprintln(os.Stderr, err)
			if printErr != nil {
				log.Println(printErr)
			}
			os.Exit(exitcode.FromError(err))
		}
	}
}
