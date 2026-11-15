package analytics

import (
	"incidentradar/internal/parser"
	"reflect"
	"testing"
)

func TestCollectorAddAndMetrics(t *testing.T) {
	collector := NewCollector()

	collector.Add(parser.LogRecord{Status: 200, Resource: "/health", ResponseSize: 10})
	collector.Add(parser.LogRecord{Status: 200, Resource: "/health", ResponseSize: 30})
	collector.Add(parser.LogRecord{Status: 500, Resource: "/error", ResponseSize: 20})

	if collector.CounterRequests != 3 {
		t.Fatalf("CounterRequests = %d, ожидалось 3", collector.CounterRequests)
	}
	if !reflect.DeepEqual(collector.StatusCounters, map[int]uint64{200: 2, 500: 1}) {
		t.Fatalf("StatusCounters = %#v, счётчики статусов отличаются", collector.StatusCounters)
	}
	if !reflect.DeepEqual(collector.ResourcesCounters, map[string]uint64{"/health": 1, "/error": 0}) {
		t.Fatalf("ResourcesCounters = %#v, счётчики ресурсов отличаются", collector.ResourcesCounters)
	}
	if collector.MaxResourcesSize() != 30 {
		t.Fatalf("MaxResourcesSize = %d, ожидалось 30", collector.MaxResourcesSize())
	}
}

func TestCollectorAggregatesAndSortsResources(t *testing.T) {
	collector := NewCollector()
	collector.ResourceRecords = []ResourceRecord{
		{Resource: "/z", Count: 1},
		{Resource: "/a", Count: 2},
		{Resource: "/b", Count: 2},
	}

	collector.SortResources()
	want := []ResourceRecord{
		{Resource: "/a", Count: 2},
		{Resource: "/b", Count: 2},
		{Resource: "/z", Count: 1},
	}
	if !reflect.DeepEqual(collector.ResourceRecords, want) {
		t.Fatalf("ResourceRecords = %#v, ожидалось %#v", collector.ResourceRecords, want)
	}
}

func TestCollectorTop10AndSizeStatistics(t *testing.T) {
	collector := NewCollector()
	for i := 0; i < 10; i++ {
		collector.ResourceRecords = append(collector.ResourceRecords, ResourceRecord{
			Resource: string(rune('a' + i)),
			Count:    uint64(i),
		})
	}
	collector.Sizes = []int64{100, 20, 40, 80, 60}

	top, err := collector.Top10resources()
	if err != nil {
		t.Fatalf("Top10resources() завершился с ошибкой: %v", err)
	}
	if len(top) != 10 || top[0].Count != 9 || top[9].Count != 0 {
		t.Fatalf("неожиданный список ресурсов: %#v", top)
	}

	average, err := collector.AVGResourceSize()
	if err != nil || average != 60 {
		t.Fatalf("AVGResourceSize() = %v, %v; ожидалось 60, nil", average, err)
	}
	p95, err := collector.P95ResourcesSize()
	if err != nil || p95 != 20 {
		t.Fatalf("P95ResourcesSize() = %d, %v; ожидалось 20, nil", p95, err)
	}
	if !reflect.DeepEqual(collector.Sizes, []int64{20, 40, 60, 80, 100}) {
		t.Fatalf("размеры отсортированы неверно: %#v", collector.Sizes)
	}
}

func TestCollectorEmptyStatisticsReturnErrors(t *testing.T) {
	collector := NewCollector()

	if _, err := collector.Top10resources(); err == nil {
		t.Fatal("ожидалась ошибка Top10resources()")
	}
	if _, err := collector.P95ResourcesSize(); err == nil {
		t.Fatal("ожидалась ошибка P95ResourcesSize()")
	}
	if _, err := collector.AVGResourceSize(); err == nil {
		t.Fatal("ожидалась ошибка AVGResourceSize()")
	}
}
