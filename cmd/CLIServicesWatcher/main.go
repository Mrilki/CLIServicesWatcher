package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/config"
	"github.com/Mrilki/CLIServicesWatcher/internal/reporter"
	"github.com/Mrilki/CLIServicesWatcher/internal/worker"
)

func main() {
	fmt.Println("Starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("\n Interrupt signal received. Shutting down gracefully...")
		cancel()
	}()

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
	workersPool.SetContext(ctx)

	tasksChan := make(chan worker.Task, len(cfg.Targets))

	go func() {
		for _, target := range cfg.Targets {
			select {
			case <-ctx.Done():
				fmt.Println("Stopped queuing new tasks.")
				close(tasksChan)
				return
			default:
			}
			tasksChan <- worker.Task{Target: target}
		}
		close(tasksChan)
	}()

	workersPool.Run(tasksChan)

	results := workersPool.GetResults()

	if len(results) > 0 {
		if err := reporter.SaveToJSON(results, "report.json"); err != nil {
			log.Fatalf("Could not save results: %v\n", err)
		}
		reporter.PrintStats(results)
	} else {
		fmt.Println("No results saved.")
	}

}
