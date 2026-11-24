// Package analytics собирает статистику запросов и фильтрует записи по времени.
package analytics

import (
	"cmp"
	"incidentradar/internal/parser"
	"math"
	"slices"
)

// Collector накапливает статистику по записям журнала доступа.
type Collector struct {
	CounterRequests   uint64
	StatusCounters    map[int]uint64
	ResourcesCounters map[string]uint64
	ResourceRecords   []ResourceRecord
	MaxSizeBytes      int64
	Sizes             []int64
}

// ResourceRecord хранит путь ресурса и счётчик обращений.
type ResourceRecord struct {
	Resource string
	Count    uint64
}

// NewCollector создаёт пустой сборщик статистики.
func NewCollector() *Collector {
	return &Collector{
		CounterRequests:   0,
		StatusCounters:    make(map[int]uint64),
		ResourcesCounters: make(map[string]uint64),
		ResourceRecords:   make([]ResourceRecord, 0),
		MaxSizeBytes:      0,
		Sizes:             []int64{},
	}
}

// Add учитывает запись лога в статистике запросов и размеров ответов.
func (c *Collector) Add(lr parser.LogRecord) {
	c.CounterRequests++
	if _, ok := c.StatusCounters[lr.Status]; !ok {
		c.StatusCounters[lr.Status] = 0
	}
	c.StatusCounters[lr.Status]++
	if _, ok := c.ResourcesCounters[lr.Resource]; !ok {
		c.ResourcesCounters[lr.Resource] = 0
		c.ResourceRecords = append(c.ResourceRecords, ResourceRecord{lr.Resource, 0})
	} else {
		c.ResourcesCounters[lr.Resource]++
		for _, res := range c.ResourceRecords {
			if res.Resource == lr.Resource {
				res.Count++
			}
		}
	}
	if lr.ResponseSize > c.MaxSizeBytes {
		c.MaxSizeBytes = lr.ResponseSize
	}
	c.Sizes = append(c.Sizes, lr.ResponseSize)
}

// SortResources сортирует ресурсы по убыванию счётчика и по пути при равенстве.
func (c *Collector) SortResources() {
	slices.SortFunc(c.ResourceRecords, func(r1, r2 ResourceRecord) int {
		if res1 := cmp.Compare(r2.Count, r1.Count); res1 != 0 {
			return res1
		}
		return cmp.Compare(r1.Resource, r2.Resource)
	})
}

// reverseSortSizes сортирует размеры ответов по возрастанию.
func (c *Collector) reverseSortSizes() {
	slices.SortFunc(c.Sizes, cmp.Compare)
}

// Top10resources возвращает первые десять ресурсов после сортировки по счётчику.
func (c *Collector) Top10resources() ([]ResourceRecord, error) {
	if len(c.ResourceRecords) == 0 {
		return nil, ErrNoResourceRecords
	}
	result := make([]ResourceRecord, 0)
	c.SortResources()
	for i := range min(10, len(c.ResourceRecords)) {
		result = append(result, c.ResourceRecords[i])
	}
	return result, nil
}

// P95ResourcesSize возвращает размер с индексом floor(0.05 × N) в порядке возрастания.
func (c *Collector) P95ResourcesSize() (int64, error) {
	if len(c.Sizes) == 0 {
		return 0, ErrNoResponseSizes
	}
	c.reverseSortSizes()
	quantity := len(c.Sizes)
	indexFor95percent := int(math.Floor(float64(quantity) * 0.05))
	return c.Sizes[indexFor95percent], nil
}

// AVGResourceSize вычисляет средний размер ответа в байтах.
func (c *Collector) AVGResourceSize() (float64, error) {
	if len(c.Sizes) == 0 {
		return 0, ErrNoResponseSizes
	}
	sum := int64(0)
	for _, size := range c.Sizes {
		sum += size
	}
	return float64(sum) / float64(len(c.Sizes)), nil
}

// MaxResourcesSize возвращает максимальный размер ответа в байтах.
func (c *Collector) MaxResourcesSize() int64 {
	return c.MaxSizeBytes
}
