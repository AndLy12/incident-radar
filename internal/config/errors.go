package config

import "errors"

// ErrInvalidConfig означает, что конфигурация не прошла проверку.
var ErrInvalidConfig = errors.New("invalid config")

// ErrMissingLogsPath означает, что путь к логам не указан.
var ErrMissingLogsPath = errors.New("logsPath required")

// ErrUnsupportedFormat означает, что указан неподдерживаемый формат отчёта.
var ErrUnsupportedFormat = errors.New("format must be json or md")

// ErrMissingOutput означает, что путь к отчёту не указан.
var ErrMissingOutput = errors.New("output must be set")

// ErrOutputExtensionMismatch означает, что расширение файла не соответствует формату отчёта.
var ErrOutputExtensionMismatch = errors.New("output file extension does not match format")

// ErrOutputExists означает, что файл отчёта уже существует.
var ErrOutputExists = errors.New("output file already exists")

// ErrInvalidDateRange означает, что начальная дата не предшествует конечной.
var ErrInvalidDateRange = errors.New("from must be < to")
