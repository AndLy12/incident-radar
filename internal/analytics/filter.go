package analytics

import "time"

// Filter задаёт необязательные границы периода для записей лога.
type Filter struct {
	From *time.Time
	To   *time.Time
}

// NewFilter создаёт фильтр с необязательными границами периода.
func NewFilter(from, to *time.Time) *Filter {
	return &Filter{
		From: from,
		To:   to,
	}
}

// after проверяет нижнюю границу периода включительно.
func (filter *Filter) after(tm *time.Time) bool {
	if filter.From == nil {
		return true
	}
	return !tm.Before(*filter.From)
}

// before проверяет верхнюю границу периода включительно.
func (filter *Filter) before(tm *time.Time) bool {
	if filter.To == nil {
		return true
	}
	return !tm.After(*filter.To)
}

// Match проверяет, входит ли время в заданный период с учётом границ.
func (filter Filter) Match(tm time.Time) bool {
	return filter.after(&tm) && filter.before(&tm)
}
