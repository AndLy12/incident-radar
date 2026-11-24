package parser

import "errors"

// ErrInvalidLogLine означает, что строка не соответствует формату журнала.
var ErrInvalidLogLine = errors.New("invalid log line")

// ErrInvalidTimestamp означает, что время запроса не удалось разобрать.
var ErrInvalidTimestamp = errors.New("parse timestamp")

// ErrInvalidResponseSize означает, что размер ответа не удалось разобрать.
var ErrInvalidResponseSize = errors.New("parse response size")

// ErrInvalidResponseStatus означает, что код ответа не удалось разобрать.
var ErrInvalidResponseStatus = errors.New("parse response status")
