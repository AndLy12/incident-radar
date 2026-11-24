package exitcode_test

import (
	"errors"
	"flag"
	"fmt"
	"incidentradar/internal/analytics"
	"incidentradar/internal/cli"
	"incidentradar/internal/config"
	"incidentradar/internal/exitcode"
	"incidentradar/internal/reader"
	"incidentradar/internal/source"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestFromError проверяет коды завершения для прямых и обёрнутых ошибок.
func TestFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{name: "success", want: 0},
		{name: "help", err: flag.ErrHelp, want: 0},
		{name: "wrapped help", err: fmt.Errorf("cli: %w", flag.ErrHelp), want: 0},
		{name: "unexpected positional argument", err: cli.ErrUnexpectedPositionalArgument, want: 2},
		{name: "invalid date", err: cli.ErrInvalidDate, want: 2},
		{name: "missing logs path", err: config.ErrMissingLogsPath, want: 2},
		{name: "unsupported report format", err: config.ErrUnsupportedFormat, want: 2},
		{name: "missing output", err: config.ErrMissingOutput, want: 2},
		{name: "wrong output extension", err: config.ErrOutputExtensionMismatch, want: 2},
		{name: "existing output found during validation", err: config.ErrOutputExists, want: 2},
		{name: "invalid date range", err: config.ErrInvalidDateRange, want: 2},
		{name: "unsupported log extension", err: source.ErrUnsupportedLogExtension, want: 2},
		{name: "unsupported log format", err: source.ErrUnsupportedLogFormat, want: 2},
		{name: "no matching logs", err: source.ErrNoMatchingLogs, want: 2},
		{name: "invalid glob", err: filepath.ErrBadPattern, want: 2},
		{
			name: "existing output found during creation",
			err:  &os.PathError{Op: "open", Path: "report.json", Err: os.ErrExist},
			want: 2,
		},
		{
			name: "wrapped config validation error",
			err:  fmt.Errorf("%w: %w", config.ErrInvalidConfig, config.ErrMissingLogsPath),
			want: 2,
		},
		{
			name: "wrapped config filesystem error",
			err:  fmt.Errorf("%w: %w", config.ErrInvalidConfig, os.ErrPermission),
			want: 1,
		},
		{
			name: "missing output directory",
			err:  &os.PathError{Op: "open", Path: "missing/report.json", Err: os.ErrNotExist},
			want: 1,
		},
		{name: "unexpected error", err: errors.New("unexpected failure"), want: 1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := exitcode.FromError(test.err); got != test.want {
				t.Errorf("FromError(%v) = %d, ожидалось %d", test.err, got, test.want)
			}
		})
	}
}

// TestFromErrorUsageFailures проверяет классификацию ошибок, возвращаемых пакетами.
func TestFromErrorUsageFailures(t *testing.T) {
	server := httptest.NewServer(http.NotFoundHandler())
	t.Cleanup(server.Close)

	_, unknownFlag := cli.Parse([]string{"--unknown"}, io.Discard)
	_, missingFlagValue := cli.Parse([]string{"--path"}, io.Discard)
	_, missingLog := source.OpenFile(filepath.Join(t.TempDir(), "missing.log"))
	invalidURL := reader.ReadAndParseFile(analytics.Filter{}, "http://[", analytics.NewCollector())
	notFound := reader.ReadAndParseFile(analytics.Filter{}, server.URL, analytics.NewCollector())

	for _, test := range []struct {
		name string
		err  error
		want error
	}{
		{name: "unknown flag", err: unknownFlag, want: cli.ErrInvalidFlags},
		{name: "missing flag value", err: missingFlagValue, want: cli.ErrInvalidFlags},
		{name: "missing input log", err: missingLog, want: source.ErrLogNotFound},
		{name: "invalid source URL", err: invalidURL, want: reader.ErrHTTPSourceUnavailable},
		{name: "HTTP 404", err: notFound, want: reader.ErrHTTPStatus},
	} {
		t.Run(test.name, func(t *testing.T) {
			if !errors.Is(test.err, test.want) {
				t.Fatalf("получено %v, ожидалась причина %v", test.err, test.want)
			}
			if got := exitcode.FromError(test.err); got != 2 {
				t.Errorf("FromError(%v) = %d, ожидалось 2", test.err, got)
			}
		})
	}
}
