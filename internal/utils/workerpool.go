package utils

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Job represents a job that can be executed by a worker
type Job interface {
	Execute() error
	GetID() string
}

// WorkerPool manages a pool of workers for concurrent processing
type WorkerPool struct {
	workers    int
	jobs       chan Job
	results    chan JobResult
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.RWMutex
	isRunning  bool
}

// JobResult represents the result of a job execution
type JobResult struct {
	JobID string
	Error error
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Job, workers*2),
		results: make(chan JobResult, workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	
	if wp.isRunning {
		return
	}
	
	wp.isRunning = true
	
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	
	if !wp.isRunning {
		return
	}
	
	wp.cancel()
	wp.wg.Wait()
	close(wp.jobs)
	close(wp.results)
	wp.isRunning = false
}

// Submit submits a job to the worker pool
func (wp *WorkerPool) Submit(job Job) error {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	
	if !wp.isRunning {
		return fmt.Errorf("worker pool is not running")
	}
	
	select {
	case wp.jobs <- job:
		return nil
	case <-wp.ctx.Done():
		return wp.ctx.Err()
	default:
		return fmt.Errorf("worker pool is full")
	}
}

// GetResults returns the results channel
func (wp *WorkerPool) GetResults() <-chan JobResult {
	return wp.results
}

// worker is the worker goroutine
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	
	for {
		select {
		case job := <-wp.jobs:
			if job == nil {
				return
			}
			
			// Execute the job
			err := job.Execute()
			
			// Send result
			select {
			case wp.results <- JobResult{JobID: job.GetID(), Error: err}:
			case <-wp.ctx.Done():
				return
			}
			
		case <-wp.ctx.Done():
			return
		}
	}
}

// IsRunning returns whether the worker pool is running
func (wp *WorkerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.isRunning
}

// GetWorkerCount returns the number of workers
func (wp *WorkerPool) GetWorkerCount() int {
	return wp.workers
}

// WaitForCompletion waits for all jobs to complete
func (wp *WorkerPool) WaitForCompletion(timeout time.Duration) error {
	done := make(chan struct{})
	
	go func() {
		wp.wg.Wait()
		close(done)
	}()
	
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("worker pool did not complete within timeout")
	}
}

// GetStats returns statistics about the worker pool
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	
	return WorkerPoolStats{
		Workers:   wp.workers,
		IsRunning: wp.isRunning,
		JobsQueued: len(wp.jobs),
		ResultsQueued: len(wp.results),
	}
}

// WorkerPoolStats represents statistics about the worker pool
type WorkerPoolStats struct {
	Workers       int
	IsRunning     bool
	JobsQueued    int
	ResultsQueued int
}

// DefaultWorkerPool creates a worker pool with default settings
func DefaultWorkerPool() *WorkerPool {
	return NewWorkerPool(10) // Default to 10 workers
}

// NewWorkerPoolWithContext creates a worker pool with a custom context
func NewWorkerPoolWithContext(ctx context.Context, workers int) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		workers: workers,
		jobs:    make(chan Job, workers*2),
		results: make(chan JobResult, workers*2),
		ctx:     ctx,
		cancel:  cancel,
	}
}
