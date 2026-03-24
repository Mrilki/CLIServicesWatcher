package main

import (
	"context"
	"flag"
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

const defaultMaxWorkers = 10

func main() {
	numWorkers := flag.Int("workers", 0, "number of workers to use")
	configPath := flag.String("config", "cfg.json", "path to config file")
	outputPath := flag.String("output", "report.json", "path to output file")

	flag.Parse()

	fmt.Println("Starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("signal handler panicked: %v", r)
			}
		}()

		<-sigs
		fmt.Println("\n Interrupt signal received. Shutting down gracefully...")
		cancel()
	}()

	fmt.Println("Loading config...")
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}
	fmt.Printf("Default timeout seconds: %d\n", cfg.Timeout)

	factory := checker.NewCheckerFactory(cfg.GetTimeoutDuration())

	workersCount := *numWorkers
	if workersCount <= 0 {
		workersCount = min(len(cfg.Targets), defaultMaxWorkers)
		fmt.Printf("Workers: %d (auto)\n", workersCount)
	} else {
		fmt.Printf("Workers: %d (manual)\n", workersCount)
	}

	workersPool := worker.NewPool(ctx, workersCount, factory)
	tasksChan := make(chan worker.Task, len(cfg.Targets))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("task sender panicked: %v", r)
			}
		}()

		for _, target := range cfg.Targets {
			select {
			case <-ctx.Done():
				fmt.Println("Stopped queuing new tasks.")
				close(tasksChan)
				return
			case tasksChan <- worker.Task{Target: target}:
			}
		}
		close(tasksChan)
	}()

	workersPool.Run(tasksChan)

	results := workersPool.GetResults()

	if len(results) > 0 {
		if err := reporter.SaveToJSON(results, *outputPath); err != nil {
			log.Printf("Warning: could not save report: %v", err)
		}
		reporter.PrintStats(results)
	} else {
		fmt.Println("No results saved.")
	}

}
