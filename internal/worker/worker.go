package worker

import (
	"context"
	"fmt"
	"sync"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/models"
)

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
}

func NewPool(ctx context.Context, workersCount int, checkerFactory *checker.Factory) *Pool {
	return &Pool{
		workersCount: workersCount,
		results:      make([]models.Result, 0),
		factory:      checkerFactory,
		ctx:          ctx,
	}
}

func (p *Pool) Run(tasks <-chan Task) {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.worker(tasks)
	}
	p.wg.Wait()
}

func (p *Pool) worker(tasks <-chan Task) {
	defer p.wg.Done()

	for task := range tasks {
		chkr, err := p.factory.New(task.Target.GetType())
		var res models.Result
		if err != nil {
			res = models.Result{
				Name:    task.Target.Name,
				Address: task.Target.Address,
				Success: false,
				Error:   fmt.Sprintf("failed to create checker: %v", err),
			}
		} else {
			res = chkr.Check(p.ctx, task.Target)
		}

		p.mux.Lock()
		p.results = append(p.results, res)
		p.mux.Unlock()
	}
}

func (p *Pool) GetResults() []models.Result {
	p.mux.Lock()
	defer p.mux.Unlock()
	results := make([]models.Result, len(p.results))
	copy(results, p.results)
	return results
}
