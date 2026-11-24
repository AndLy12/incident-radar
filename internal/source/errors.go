package source

import "errors"

// ErrLogNotFound означает, что входной файл лога отсутствует.
var ErrLogNotFound = errors.New("input log not found")

// ErrUnsupportedLogExtension означает, что путь к логу имеет неподдерживаемое расширение.
var ErrUnsupportedLogExtension = errors.New("path must end with .txt or .log")

// ErrUnsupportedLogFormat означает, что маска нашла файл с неподдерживаемым расширением.
var ErrUnsupportedLogFormat = errors.New("file format not supported")

// ErrNoMatchingLogs означает, что маска не нашла ни одного лога.
var ErrNoMatchingLogs = errors.New("no logs found matching glob")

// ErrLogIsDirectory означает, что путь является директорией а не файлом
var ErrLogIsDirectory = errors.New("log path is a directory")
