// Package source находит и открывает локальные файлы логов.
package source

import (
	"errors"
	"fmt"
	"incidentradar/internal/config"
	"os"
	"path/filepath"
)

// OpenFile открывает локальный файл лога с расширением .log или .txt.
func OpenFile(path string) (*os.File, error) {
	if filepath.Ext(path) != ".log" && filepath.Ext(path) != ".txt" {
		return nil, ErrUnsupportedLogExtension
	}

	f, err := os.Open(path) // #nosec G304 -- Пользователь может намеренно открыть файл не из этой директории
	if errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("%w: %w", ErrLogNotFound, err)
	}
	if err != nil {
		return nil, err
	}

	fileStat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if fileStat != nil && fileStat.IsDir() {
		return nil, ErrLogIsDirectory
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
			return nil, ErrUnsupportedLogFormat
		}
	}
	if c == 0 {
		return nil, ErrNoMatchingLogs
	}
	return files, nil
}
