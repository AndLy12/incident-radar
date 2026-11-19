// Package source находит и открывает локальные файлы логов.
package source

import (
	"fmt"
	"incidentradar/internal/config"
	"os"
	"path/filepath"
)

// OpenFile открывает локальный файл лога с расширением .log или .txt.
func OpenFile(path string) (*os.File, error) {
	if filepath.Ext(path) != ".log" && filepath.Ext(path) != ".txt" {
		return nil, fmt.Errorf("path must end with .txt or .log")
	}

	f, err := os.Open(path) // #nosec G304 -- Пользователь может намеренно открыть файл не из этой директории
	if err != nil {
		return nil, err
	}
	return f, nil
}

// ExpandGlob возвращает совпавшие с маской пути и проверяет расширения файлов.
func ExpandGlob(path string) ([]string, error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return nil, err
	}

	c := 0
	for _, file := range files {
		c++
		flagExtNotInLogsFormatExtensions := true
		for _, ext := range config.LogsFormatExtensions {
			if filepath.Ext(file) == ext {
				flagExtNotInLogsFormatExtensions = false
			}
		}
		if flagExtNotInLogsFormatExtensions {
			return nil, fmt.Errorf("file format not supported")
		}
	}
	if c == 0 {
		return nil, fmt.Errorf("no logs found matching glob")
	}
	return files, nil
}
