package reader

import "errors"

// ErrHTTPSourceUnavailable означает, что HTTP-запрос к источнику не удалось выполнить.
var ErrHTTPSourceUnavailable = errors.New("HTTP log source unavailable")

// ErrHTTPStatus означает, что источник вернул HTTP-статус за пределами 2xx.
var ErrHTTPStatus = errors.New("unexpected HTTP status from log source")
