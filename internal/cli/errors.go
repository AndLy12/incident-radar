package cli

import "errors"

// ErrInvalidFlags означает, что аргументы командной строки не удалось разобрать.
var ErrInvalidFlags = errors.New("invalid flags")

// ErrUnexpectedPositionalArgument означает, что после флагов передан позиционный аргумент.
var ErrUnexpectedPositionalArgument = errors.New("unexpected positional arguments")

// ErrInvalidDate означает, что значение флага даты не соответствует формату YYYY-MM-DD.
var ErrInvalidDate = errors.New("invalid date")
