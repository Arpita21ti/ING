package main

// import (
// 	"context"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"server/internal/config"
// 	"server/internal/infrastructure/database/postgres"
// 	"server/internal/worker"
// 	"server/pkg/logger"
// )

// func main() {
// 	// Initialize logger
// 	log := logger.NewLogger()
// 	log.Info("Starting TNP RGPV Background Worker...")

// 	// Load configuration
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatal("Failed to load configuration", "error", err)
// 	}

// 	// Connect to PostgreSQL
// 	db, err := postgres.NewConnection(cfg.Database)
// 	if err != nil {
// 		log.Fatal("Failed to connect to PostgreSQL", "error", err)
// 	}

// 	// Ensure database connection is closed when the application exits
// 	defer func() {
// 		db.Close()
// 		log.Info("PostgreSQL connection pool closed")
// 	}()

// 	// Create context for worker coordination
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// Initialize and start workers
// 	workers := []worker.Worker{
// 		worker.NewEmailWorker(db, cfg, log),
// 		worker.NewNotificationWorker(db, cfg, log),
// 		worker.NewAnalyticsWorker(db, cfg, log),
// 		// Add additional workers as needed
// 	}

// 	// Start all workers
// 	for _, w := range workers {
// 		go func(worker worker.Worker) {
// 			if err := worker.Start(ctx); err != nil {
// 				log.Error("Worker failed", "worker", worker.Name(), "error", err)
// 			}
// 		}(w)
// 		log.Info("Started worker", "worker", w.Name())
// 	}

// 	// Set up signal handling for graceful shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	log.Info("Shutdown signal received, stopping workers...")

// 	// Cancel context to stop all workers
// 	cancel()

// 	// Allow workers time to clean up
// 	time.Sleep(3 * time.Second)

// 	log.Info("All workers stopped successfully")
// }