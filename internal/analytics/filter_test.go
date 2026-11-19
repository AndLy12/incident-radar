package analytics

import (
	"testing"
	"time"
)

// TestFilterMatchesInclusiveBounds проверяет включение обеих границ периода.
func TestFilterMatchesInclusiveBounds(t *testing.T) {
	from := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	to := from.Add(24 * time.Hour)
	filter := NewFilter(&from, &to)

	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		{"before", from.Add(-time.Nanosecond), false},
		{"at from", from, true},
		{"inside", from.Add(time.Hour), true},
		{"at to", to, true},
		{"after", to.Add(time.Nanosecond), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := filter.Match(test.time); got != test.want {
				t.Fatalf("Match() = %v, ожидалось %v", got, test.want)
			}
		})
	}
}

// TestFilterOptionalBounds проверяет фильтр с одной границей или без границ.
func TestFilterOptionalBounds(t *testing.T) {
	moment := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	for name, filter := range map[string]Filter{
		"no bounds": {},
		"from only": *NewFilter(&moment, nil),
		"to only":   *NewFilter(nil, &moment),
	} {
		t.Run(name, func(t *testing.T) {
			if !filter.Match(moment) {
				t.Fatal("граничный момент должен соответствовать фильтру")
			}
		})
	}
}
