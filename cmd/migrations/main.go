package main

import (
	"flag"
	"fmt"
	"os"

	"server/internal/config"
	"server/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// Command line flags
	var (
		up      = flag.Bool("up", false, "Migrate the DB to the most recent version")
		down    = flag.Bool("down", false, "Roll back the last migration")
		version = flag.Int("version", -1, "Migrate to a specific version")
		steps   = flag.Int("step", 0, "Number of migrations to apply (can be negative for rollback)")
	)
	flag.Parse()

	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting TNP RGPV Database Migration Tool...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Construct database URL for migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.UserName,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DatabaseName,
		cfg.Database.SSLMode,
	)

	// Initialize migration instance
	m, err := migrate.New("file://migrations/postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to create migration instance", "error", err)
	}
	defer m.Close()

	// Execute migration based on flags
	if *up {
		log.Info("Applying all pending migrations...")
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration failed", "error", err)
		}
		log.Info("Database migrated successfully")
	} else if *down {
		log.Info("Rolling back the last migration...")
		if err := m.Steps(-1); err != nil {
			log.Fatal("Migration rollback failed", "error", err)
		}
		log.Info("Last migration rolled back successfully")
	} else if *version >= 0 {
		log.Info("Migrating to specific version...", "version", *version)
		if err := m.Migrate(uint(*version)); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration to version failed", "error", err)
		}
		log.Info("Migration to version completed successfully", "version", *version)
	} else if *steps != 0 {
		log.Info("Applying migration steps...", "steps", *steps)
		if err := m.Steps(*steps); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Migration steps failed", "error", err)
		}
		log.Info("Migration steps completed successfully", "steps", *steps)
	} else {
		// No operation specified
		currentVersion, dirty, err := m.Version()
		if err != nil && err != migrate.ErrNilVersion {
			log.Fatal("Failed to get migration version", "error", err)
		}

		if err == migrate.ErrNilVersion {
			fmt.Println("Database is not yet migrated")
		} else {
			dirtyStatus := ""
			if dirty {
				dirtyStatus = " (DIRTY)"
			}
			fmt.Printf("Current migration version: %d%s\n", currentVersion, dirtyStatus)
		}
		
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
}