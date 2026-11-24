// Package exitcode определяет коды завершения Incident Radar.
package exitcode

import (
	"errors"
	"flag"
	"incidentradar/internal/cli"
	"incidentradar/internal/config"
	"incidentradar/internal/reader"
	"incidentradar/internal/source"
	"os"
	"path/filepath"
)

const (
	// Success означает штатное завершение, включая вывод справки.
	Success int = 0

	// InternalError означает непредвиденную ошибку.
	InternalError int = 1

	// UsageError означает ошибку использования или ожидаемую ошибку источника.
	UsageError int = 2
)

// FromError возвращает код завершения с учётом обёрнутых ошибок.
func FromError(err error) int {
	switch {
	case err == nil, errors.Is(err, flag.ErrHelp):
		return Success
	case errors.Is(err, cli.ErrInvalidFlags),
		errors.Is(err, cli.ErrUnexpectedPositionalArgument),
		errors.Is(err, cli.ErrInvalidDate),
		errors.Is(err, config.ErrMissingLogsPath),
		errors.Is(err, config.ErrUnsupportedFormat),
		errors.Is(err, config.ErrMissingOutput),
		errors.Is(err, config.ErrOutputExtensionMismatch),
		errors.Is(err, config.ErrOutputExists),
		errors.Is(err, config.ErrInvalidDateRange),
		errors.Is(err, source.ErrUnsupportedLogExtension),
		errors.Is(err, source.ErrUnsupportedLogFormat),
		errors.Is(err, source.ErrNoMatchingLogs),
		errors.Is(err, source.ErrLogNotFound),
		errors.Is(err, reader.ErrHTTPSourceUnavailable),
		errors.Is(err, reader.ErrHTTPStatus),
		errors.Is(err, filepath.ErrBadPattern),
		errors.Is(err, source.ErrLogIsDirectory),
		errors.Is(err, os.ErrExist):
		return UsageError
	default:
		return InternalError
	}
}
