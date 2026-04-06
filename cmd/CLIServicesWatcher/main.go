package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mrilki/CLIServicesWatcher/internal/checker"
	"github.com/Mrilki/CLIServicesWatcher/internal/config"
	"github.com/Mrilki/CLIServicesWatcher/internal/logger"
	"github.com/Mrilki/CLIServicesWatcher/internal/models"
	"github.com/Mrilki/CLIServicesWatcher/internal/reporter"
	"github.com/Mrilki/CLIServicesWatcher/internal/worker"
)

func main() {
	numWorkers := flag.Int("workers", 0, "number of workers to use")
	configPath := flag.String("config", "cfg.json", "path to config file")
	outputPath := flag.String("output", "report.json", "path to output file")
	verbose := flag.Bool("verbose", false, "enable debug logging")
	flag.Parse()

	level := slog.LevelInfo
	if *verbose {
		level = slog.LevelDebug
	}
	log := logger.Init(level)

	log.Info("Starting application",
		"config", *configPath,
		"output", *outputPath,
		"workers", *numWorkers,
		"verbose", *verbose,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.Load(*configPath, log)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			log.Warn("Config file not found, using default", "path", *configPath)
			cfg = models.GetDefaultConf()
		} else {
			log.Error("Could not load config", "error", err)
			os.Exit(1)
		}
	}

	factory := checker.NewCheckerFactory(cfg.GetTimeoutDuration())

	workersCount := *numWorkers
	source := "manual"
	if workersCount <= 0 {
		workersCount = min(len(cfg.Targets), worker.DefaultMaxWorkers)
		source = "auto"
	}

	log.Info("Worker pool initialized",
		"workers", workersCount,
		"source", source,
	)

	workersPool := worker.NewPool(ctx, workersCount, factory, log)
	tasksChan := make(chan worker.Task, len(cfg.Targets))

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Warn("Task sender panicked", "error", r)
			}
		}()

		for _, target := range cfg.Targets {
			select {
			case <-ctx.Done():
				log.Info("Stopped queuing new tasks due to context cancellation")
				close(tasksChan)
				return
			case tasksChan <- worker.Task{Target: target}:
			}
		}
		close(tasksChan)
	}()

	done := make(chan struct{})
	go func() {
		workersPool.Run(tasksChan)
		close(done)
	}()

	select {
	case <-done:
		log.Info("All targets checked successfully")
	case sig := <-sigs:
		log.Info("Received signal, shutting down...", "signal", sig)
		cancel()

		select {
		case <-done:
			log.Info("Workers finished gracefully")
		case <-time.After(30 * time.Second):
			log.Warn("Graceful shutdown timed out, forcing exit")
			os.Exit(1)
		}
	}

	results := workersPool.GetResults()

	if len(results) > 0 {
		if err := reporter.SaveToJSON(results, *outputPath, log); err != nil {
			log.Warn("Could not save report", "error", err, "path", *outputPath)
		}
		reporter.PrintStats(results)
	} else {
		log.Warn("No results to save")
	}

	log.Info("Application completed")
}
