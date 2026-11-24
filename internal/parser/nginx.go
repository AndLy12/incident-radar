// Package parser разбирает записи журнала доступа NGINX.
package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// SplitLine разбирает строку журнала доступа NGINX в LogRecord.
func SplitLine(s string) (LogRecord, error) {
	re := regexp.MustCompile(
		"^(\\S+)\\s-\\s" + // ip
			"(\\S+)\\s\\" + // remote user
			"[([^]]+)\\]\\s" + // time
			"\"([A-Z]+)\\s" + // method
			"([^\"\\s]+)\\s" + // resource
			"([^\"\\s]+)\"\\s" + // protocol
			"([0-9]{3})\\s" + // status
			"([0-9]+)\\s" + // size
			"\"([^\"]*)\"\\s" + // http_referer
			"\"([^\"]*)\"$", // http_user_agent
	)
	matches := re.FindStringSubmatch(s)
	if matches == nil {
		return LogRecord{}, fmt.Errorf("%w: %s", ErrInvalidLogLine, s)
	}
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[3])
	if err != nil {
		return LogRecord{}, fmt.Errorf("%w %s: %w", ErrInvalidTimestamp, matches[3], err)
	}
	size, err := strconv.ParseInt(matches[8], 10, 64)
	if err != nil {
		return LogRecord{}, fmt.Errorf("%w %q: %w", ErrInvalidResponseSize, matches[8], err)
	}

	status, err := strconv.Atoi(matches[7])
	if err != nil {
		return LogRecord{}, fmt.Errorf("%w %s: %w", ErrInvalidResponseStatus, matches[7], err)
	}

	return LogRecord{
		matches[1],
		timestamp,
		matches[4],
		matches[5],
		matches[6],
		status,
		size,
	}, nil
}
