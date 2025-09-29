package metrics

import (
	"sync"
	"time"
)

// Metrics represents a metrics collector
type Metrics struct {
	mu   sync.RWMutex
	data map[string]interface{}
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		data: make(map[string]interface{}),
	}
}

// RecordDuration records the duration of an operation
func (m *Metrics) RecordDuration(operation string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[operation+"_duration_ms"] = duration.Milliseconds()
	m.data[operation+"_duration_seconds"] = duration.Seconds()
}

// IncrementCounter increments a counter for an operation
func (m *Metrics) IncrementCounter(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if count, ok := m.data[operation+"_count"].(int); ok {
		m.data[operation+"_count"] = count + 1
	} else {
		m.data[operation+"_count"] = 1
	}
}

// DecrementCounter decrements a counter for an operation
func (m *Metrics) DecrementCounter(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if count, ok := m.data[operation+"_count"].(int); ok && count > 0 {
		m.data[operation+"_count"] = count - 1
	}
}

// SetGauge sets a gauge value for an operation
func (m *Metrics) SetGauge(operation string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[operation+"_gauge"] = value
}

// IncrementGauge increments a gauge value for an operation
func (m *Metrics) IncrementGauge(operation string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if gauge, ok := m.data[operation+"_gauge"].(float64); ok {
		m.data[operation+"_gauge"] = gauge + value
	} else {
		m.data[operation+"_gauge"] = value
	}
}

// DecrementGauge decrements a gauge value for an operation
func (m *Metrics) DecrementGauge(operation string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if gauge, ok := m.data[operation+"_gauge"].(float64); ok {
		m.data[operation+"_gauge"] = gauge - value
	} else {
		m.data[operation+"_gauge"] = -value
	}
}

// RecordHistogram records a histogram value for an operation
func (m *Metrics) RecordHistogram(operation string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Store histogram data as a slice of values
	key := operation + "_histogram"
	if hist, ok := m.data[key].([]float64); ok {
		m.data[key] = append(hist, value)
	} else {
		m.data[key] = []float64{value}
	}
}

// GetCounter returns the counter value for an operation
func (m *Metrics) GetCounter(operation string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if count, ok := m.data[operation+"_count"].(int); ok {
		return count
	}
	return 0
}

// GetGauge returns the gauge value for an operation
func (m *Metrics) GetGauge(operation string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if gauge, ok := m.data[operation+"_gauge"].(float64); ok {
		return gauge
	}
	return 0
}

// GetDuration returns the duration for an operation
func (m *Metrics) GetDuration(operation string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if duration, ok := m.data[operation+"_duration_ms"].(int64); ok {
		return time.Duration(duration) * time.Millisecond
	}
	return 0
}

// GetHistogram returns the histogram data for an operation
func (m *Metrics) GetHistogram(operation string) []float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if hist, ok := m.data[operation+"_histogram"].([]float64); ok {
		return hist
	}
	return nil
}

// GetAllMetrics returns all metrics data
func (m *Metrics) GetAllMetrics() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	// Create a copy of the data to avoid race conditions
	result := make(map[string]interface{})
	for k, v := range m.data {
		result[k] = v
	}
	return result
}

// Reset resets all metrics
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]interface{})
}

// ResetOperation resets metrics for a specific operation
func (m *Metrics) ResetOperation(operation string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Remove all metrics for this operation
	delete(m.data, operation+"_count")
	delete(m.data, operation+"_gauge")
	delete(m.data, operation+"_duration_ms")
	delete(m.data, operation+"_duration_seconds")
	delete(m.data, operation+"_histogram")
}

// GetStats returns statistics for an operation
func (m *Metrics) GetStats(operation string) OperationStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	stats := OperationStats{
		Operation: operation,
	}
	
	if count, ok := m.data[operation+"_count"].(int); ok {
		stats.Count = count
	}
	
	if gauge, ok := m.data[operation+"_gauge"].(float64); ok {
		stats.Gauge = gauge
	}
	
	if duration, ok := m.data[operation+"_duration_ms"].(int64); ok {
		stats.Duration = time.Duration(duration) * time.Millisecond
	}
	
	if hist, ok := m.data[operation+"_histogram"].([]float64); ok {
		stats.Histogram = hist
		stats.HistogramCount = len(hist)
	}
	
	return stats
}

// OperationStats represents statistics for an operation
type OperationStats struct {
	Operation      string        `json:"operation"`
	Count          int           `json:"count"`
	Gauge          float64       `json:"gauge"`
	Duration       time.Duration `json:"duration"`
	Histogram      []float64     `json:"histogram"`
	HistogramCount int           `json:"histogram_count"`
}

// RecordDatabaseOperation records a database operation
func (m *Metrics) RecordDatabaseOperation(operation string, table string, duration time.Duration, success bool) {
	m.RecordDuration("db_"+operation, duration)
	m.IncrementCounter("db_"+operation)
	
	if success {
		m.IncrementCounter("db_success")
	} else {
		m.IncrementCounter("db_error")
	}
	
	// Record table-specific metrics
	m.IncrementCounter("db_table_" + table)
}

// RecordNetworkOperation records a network operation
func (m *Metrics) RecordNetworkOperation(operation string, url string, duration time.Duration, statusCode int, success bool) {
	m.RecordDuration("network_"+operation, duration)
	m.IncrementCounter("network_"+operation)
	
	if success {
		m.IncrementCounter("network_success")
	} else {
		m.IncrementCounter("network_error")
	}
	
	// Record status code metrics
	m.IncrementCounter("network_status_" + string(rune(statusCode)))
}

// RecordFileOperation records a file operation
func (m *Metrics) RecordFileOperation(operation string, filePath string, duration time.Duration, success bool) {
	m.RecordDuration("file_"+operation, duration)
	m.IncrementCounter("file_"+operation)
	
	if success {
		m.IncrementCounter("file_success")
	} else {
		m.IncrementCounter("file_error")
	}
}

// RecordCacheOperation records a cache operation
func (m *Metrics) RecordCacheOperation(operation string, key string, duration time.Duration, hit bool) {
	m.RecordDuration("cache_"+operation, duration)
	m.IncrementCounter("cache_"+operation)
	
	if hit {
		m.IncrementCounter("cache_hit")
	} else {
		m.IncrementCounter("cache_miss")
	}
}

// RecordWorkerPoolOperation records a worker pool operation
func (m *Metrics) RecordWorkerPoolOperation(operation string, workerID int, duration time.Duration, success bool) {
	m.RecordDuration("worker_"+operation, duration)
	m.IncrementCounter("worker_"+operation)
	
	if success {
		m.IncrementCounter("worker_success")
	} else {
		m.IncrementCounter("worker_error")
	}
}

// GetSystemMetrics returns system-level metrics
func (m *Metrics) GetSystemMetrics() SystemMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	metrics := SystemMetrics{}
	
	// Calculate totals
	metrics.TotalOperations = m.getTotalCount("_count")
	metrics.TotalErrors = m.getTotalCount("_error")
	metrics.TotalSuccess = m.getTotalCount("_success")
	
	// Calculate averages
	metrics.AverageDuration = m.getAverageDuration()
	
	return metrics
}

// SystemMetrics represents system-level metrics
type SystemMetrics struct {
	TotalOperations  int           `json:"total_operations"`
	TotalErrors      int           `json:"total_errors"`
	TotalSuccess     int           `json:"total_success"`
	AverageDuration  time.Duration `json:"average_duration"`
}

// getTotalCount calculates the total count for a suffix
func (m *Metrics) getTotalCount(suffix string) int {
	total := 0
	for key, value := range m.data {
		if len(key) > len(suffix) && key[len(key)-len(suffix):] == suffix {
			if count, ok := value.(int); ok {
				total += count
			}
		}
	}
	return total
}

// getAverageDuration calculates the average duration
func (m *Metrics) getAverageDuration() time.Duration {
	var totalDuration time.Duration
	var count int
	
	for key, value := range m.data {
		if len(key) > 11 && key[len(key)-11:] == "_duration_ms" {
			if duration, ok := value.(int64); ok {
				totalDuration += time.Duration(duration) * time.Millisecond
				count++
			}
		}
	}
	
	if count == 0 {
		return 0
	}
	
	return totalDuration / time.Duration(count)
}
