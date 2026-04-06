package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/models"
	"github.com/Mrilki/CLIServicesWatcher/internal/testutil"
)

type mockChecker struct {
	result            models.Result
	delay             time.Duration
	shouldPanic       bool
	mu                sync.Mutex
	maxConcurrent     int
	currentConcurrent int
}

func (m *mockChecker) Check(ctx context.Context, target models.Target) models.Result {
	m.mu.Lock()
	m.currentConcurrent++
	if m.currentConcurrent > m.maxConcurrent {
		m.maxConcurrent = m.currentConcurrent
	}
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.currentConcurrent--
		m.mu.Unlock()
	}()

	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return models.Result{
				Name:    target.Name,
				Success: false,
				Error:   "context cancelled",
			}
		}

	}
	if m.shouldPanic {
		panic("test shouldPanic")
	}
	return m.result
}

type mockFactory struct {
	checker *mockChecker
}

func (f *mockFactory) New(checkType models.CheckType) (checker.Checker, error) {
	return f.checker, nil
}

func TestWorker_Run_Success(t *testing.T) {
	mockChkr := &mockChecker{
		result: models.Result{
			Name:    "test",
			Address: "http://test.com",
			Type:    models.CheckTypeHTTP,
			Success: true,
			Latency: models.Duration(50 * time.Millisecond),
		},
	}
	factory := &mockFactory{checker: mockChkr}
	ctx := context.Background()

	pool := NewPool(ctx, 3, factory, testutil.DiscardLogger())

	tasks := make(chan Task, 5)
	for i := 0; i < 5; i++ {
		tasks <- Task{
			Target: models.Target{
				Name:    "Test",
				Address: "http://test.com",
				Type:    models.CheckTypeHTTP,
			},
		}
	}
	close(tasks)

	pool.Run(tasks)

	results := pool.GetResults()

	if len(results) != 5 {
		t.Errorf("expected 5 results, got %d", len(results))
	}
	for i, result := range results {
		if !result.Success {
			t.Errorf("result %d expected success, got error %v", i, result.Error)
		}

	}
}

func TestWorker_PanicRecovery(t *testing.T) {
	mockChkr := &mockChecker{shouldPanic: true}
	factory := &mockFactory{checker: mockChkr}
	ctx := context.Background()

	pool := NewPool(ctx, 3, factory, testutil.DiscardLogger())
	tasks := make(chan Task, 2)
	tasks <- Task{Target: models.Target{Name: "Panic1", Type: models.CheckTypeHTTP}}
	tasks <- Task{Target: models.Target{Name: "Panic2", Type: models.CheckTypeHTTP}}
	close(tasks)

	pool.Run(tasks)

	results := pool.GetResults()

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	for i, r := range results {
		if r.Success == true {
			t.Errorf("result %d: expected failure after shouldPanic, got success", i)
		}
		if r.Error == "" {
			t.Errorf("result %d: expected error message, got empty", i)
		}
	}
}

func TestWorker_ContextCancellation(t *testing.T) {
	mockChkr := &mockChecker{
		delay:  1 * time.Second,
		result: models.Result{Name: "Cancel", Success: true},
	}
	factory := &mockFactory{checker: mockChkr}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	pool := NewPool(ctx, 2, factory, testutil.DiscardLogger())

	tasks := make(chan Task, 5)
	for i := 0; i < 5; i++ {
		tasks <- Task{Target: models.Target{Name: "Test", Type: models.CheckTypeHTTP}}
	}
	close(tasks)

	pool.Run(tasks)

	results := pool.GetResults()

	for i, r := range results {
		if r.Success {
			t.Errorf("result %d: expected cancellation, got success", i)
		}
		if r.Error != "context cancelled" {
			t.Errorf("result %d: expected 'context cancelled', got %q", i, r.Error)
		}
	}
}

func TestWorker_ParallelExecution(t *testing.T) {
	mockChkr := &mockChecker{
		delay:  100 * time.Millisecond,
		result: models.Result{Success: true},
	}
	factory := &mockFactory{checker: mockChkr}

	ctx := context.Background()
	pool := NewPool(ctx, 3, factory, testutil.DiscardLogger())

	tasks := make(chan Task, 6)
	for i := 0; i < 6; i++ {
		tasks <- Task{Target: models.Target{Name: "Test", Type: models.CheckTypeHTTP}}
	}
	close(tasks)

	pool.Run(tasks)

	results := pool.GetResults()
	if len(results) != 6 {
		t.Errorf("expected 6 results, got %d", len(results))
	}

	mockChkr.mu.Lock()
	maxConcurrent := mockChkr.maxConcurrent
	mockChkr.mu.Unlock()

	if maxConcurrent < 2 {
		t.Errorf("expected at least 2 concurrent executions, got %d", maxConcurrent)
	}
	if maxConcurrent > 3 {
		t.Errorf("expected at most 3, worker pool limit")
	}
}
