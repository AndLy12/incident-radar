package app

import (
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/cli"
	"incidentradar/internal/reader"
	"os"
	"path/filepath"
	"time"
)

func Run(args []string) error {
	// Создаем новый конфиг, распарсив команду и ее флаги
	config, err := cli.Parse(args, os.Stderr)
	if err != nil {
		return err
	}

	filePath, err := filepath.Glob(config.LogsPath())
	if err != nil {
		return err
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
	for _, file := range filePath {
		validLogs, err := reader.ReadAndParseFile(*filter, file)
		if err != nil {
			return err
		}
		fmt.Println(validLogs)
	}
	return nil
}
