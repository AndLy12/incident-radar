package parser

import (
	"time"
)

type LogRecord struct {
	IP           string
	Time         time.Time
	Method       string
	Resource     string
	Protocol     string
	Status       int
	ResponseSize int64
}
