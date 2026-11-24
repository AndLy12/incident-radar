package analytics

import "errors"

// ErrNoResourceRecords означает, что список ресурсов пуст.
var ErrNoResourceRecords = errors.New("no resource records found")

// ErrNoResponseSizes означает, что размеры ответов ещё не собраны.
var ErrNoResponseSizes = errors.New("no resources")
