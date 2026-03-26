package worker

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"sync"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

const DefaultMaxWorkers = 10

type Task struct {
	Target models.Target
}

type Pool struct {
	workersCount int
	results      []models.Result
	factory      *checker.Factory
	mux          sync.Mutex
	wg           sync.WaitGroup
	ctx          context.Context
	log          *slog.Logger
}

func NewPool(ctx context.Context, workersCount int, checkerFactory *checker.Factory, log *slog.Logger) *Pool {
	return &Pool{
		workersCount: workersCount,
		results:      make([]models.Result, 0),
		factory:      checkerFactory,
		ctx:          ctx,
		log:          log,
	}
}

func (p *Pool) Run(tasks <-chan Task) {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.worker(i, tasks)
	}
	p.wg.Wait()

	p.log.Info("Worker pool completed", "results", len(p.results))
}

func (p *Pool) worker(id int, tasks <-chan Task) {
	defer p.wg.Done()

	defer func() {
		if r := recover(); r != nil {
			p.log.Warn("Worker panicked",
				"worker", id,
				"error", r,
				"stack", string(debug.Stack()))
		}
	}()

	p.log.Debug("Worker started", "worker", id)

	for task := range tasks {
		p.log.Debug("Processing task",
			"worker", id,
			"target", task.Target.Name,
			"type", task.Target.Type)

		chkr, err := p.factory.New(task.Target.GetType())
		var res models.Result

		if err != nil {
			p.log.Error("Failed to create checker",
				"worker", id,
				"target", task.Target.Name,
				"type", task.Target.Type,
				"error", err,
			)

			res = models.Result{
				Name:    task.Target.Name,
				Address: task.Target.Address,
				Type:    task.Target.Type,
				Success: false,
				Error:   fmt.Sprintf("failed to create checker: %v", err),
			}
		} else {
			func() {
				defer func() {
					if r := recover(); r != nil {
						p.log.Warn("Checker panicked",
							"worker", id,
							"target", task.Target.Name,
							"error", r,
						)
						res = models.Result{
							Name:    task.Target.Name,
							Address: task.Target.Address,
							Type:    task.Target.Type,
							Success: false,
							Error:   fmt.Sprintf("checker panicked: %v", r),
						}
					}
				}()
				p.log.Debug("Running checker",
					"worker", id,
					"target", task.Target.Name,
					"type", task.Target.Type,
				)

				res = chkr.Check(p.ctx, task.Target)
			}()
		}

		p.mux.Lock()
		p.results = append(p.results, res)
		p.mux.Unlock()

		p.log.Debug("Task completed",
			"worker", id,
			"target", task.Target.Name,
			"success", res.Success,
		)
	}
	p.log.Debug("Worker finished", "worker", id)
}

func (p *Pool) GetResults() []models.Result {
	p.mux.Lock()
	defer p.mux.Unlock()
	results := make([]models.Result, len(p.results))
	copy(results, p.results)
	return results
}
