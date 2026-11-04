package analytics

import "time"

type Filter struct {
	From *time.Time
	To   *time.Time
}

func NewFilter(from, to *time.Time) *Filter {
	return &Filter{
		From: from,
		To:   to,
	}
}

func (filter *Filter) after(tm *time.Time) bool {
	if filter.From == nil {
		return true
	}
	return !tm.Before(*filter.From)
}

func (filter *Filter) before(tm *time.Time) bool {
	if filter.To == nil {
		return true
	}
	return !tm.After(*filter.To)
}

func (filter Filter) Match(tm time.Time) bool {
	return filter.after(&tm) && filter.before(&tm)
}
