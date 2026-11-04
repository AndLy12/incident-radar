package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

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
		return LogRecord{}, fmt.Errorf("invalid log line: %s", s)
	}
	timestamp, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[3])
	if err != nil {
		return LogRecord{}, fmt.Errorf("parse timestamp %s: %w", matches[3], err)
	}
	size, err := strconv.ParseInt(matches[8], 10, 64)
	if err != nil {
		return LogRecord{}, fmt.Errorf(
			"parse response size %q: %w",
			matches[8],
			err,
		)
	}

	status, err := strconv.Atoi(matches[7])
	if err != nil {
		return LogRecord{}, fmt.Errorf("parse response status %s: %w", matches[7], err)
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
