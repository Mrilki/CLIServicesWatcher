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
	result       []models.Result
	checker      checker.Checker
	mux          sync.Mutex
	wg           sync.WaitGroup
	ctx          context.Context
}

func NewPool(workersCount int, checker checker.Checker) *Pool {
	return &Pool{
		workersCount: workersCount,
		result:       make([]models.Result, 0),
		checker:      checker,
		ctx:          context.Background(),
	}
}

func (p *Pool) SetContext(ctx context.Context) {
	p.ctx = ctx
}

func (p *Pool) Run(tasks <-chan Task) {
	for i := 0; i < p.workersCount; i++ {
		p.wg.Add(1)
		go p.worker(i, tasks)
	}
	p.wg.Wait()
}

func (p *Pool) worker(id int, tasks <-chan Task) {
	defer p.wg.Done()

	for task := range tasks {
		select {
		case <-p.ctx.Done():
			return
		default:

		}
		res := p.checker.Check(p.ctx, task.Target)

		p.mux.Lock()
		p.result = append(p.result, res)
		p.mux.Unlock()

		fmt.Printf("[Worker %d] %s\n", id, task.Target.Name)
	}
}

func (p *Pool) GetResults() []models.Result {
	p.mux.Lock()
	defer p.mux.Unlock()
	result := make([]models.Result, len(p.result))
	copy(result, p.result)
	return result
}
