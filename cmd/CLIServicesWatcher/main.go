package main

import (
	"fmt"
	"log"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/config"
	"github.com/Mrilki/CLIServicesWatcher/internal/reporter"
	"github.com/Mrilki/CLIServicesWatcher/internal/worker"
)

func main() {
	fmt.Println("Starting...")
	fmt.Println("Loading config...")
	cfg, err := config.Load("cfg.json")
	if err != nil {
		log.Fatalf("Could not load config: %v\n", err)
	}
	fmt.Printf("Default timeout seconds: %d\n", cfg.Timeout)

	var monitor checker.Checker
	monitor = checker.NewHttpChecker(cfg.GetTimeoutDuration())

	workersCount := len(cfg.Targets)
	if workersCount > 10 {
		workersCount = 10
	}
	workersPool := worker.NewPool(workersCount, monitor)

	tasksChan := make(chan worker.Task, len(cfg.Targets))

	go func() {
		for _, target := range cfg.Targets {
			tasksChan <- worker.Task{Target: target}
		}
		close(tasksChan)
	}()

	workersPool.Run(tasksChan)

	results := workersPool.GetResults()

	if err := reporter.SaveToJSON(results, "report.json"); err != nil {
		log.Fatalf("Could not save results: %v\n", err)
	}
	reporter.PrintStats(results)

}
