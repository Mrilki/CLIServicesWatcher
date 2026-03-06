package worker

import (
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
}

func NewPool(WorkersCount int, checker checker.Checker) *Pool {
	return &Pool{
		workersCount: WorkersCount,
		result:       make([]models.Result, 0),
		checker:      checker,
	}
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
		res := p.checker.Check(task.Target)

		p.mux.Lock()
		p.result = append(p.result, res)
		p.mux.Unlock()

		fmt.Printf("[Worker %d] %s\n", id, task.Target.Name)
	}
}

func (p *Pool) GetResults() []models.Result {
	return p.result
}
