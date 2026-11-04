package source

import (
	"fmt"
	"os"
	"path/filepath"
)

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
