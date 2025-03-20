// TODO: Update ad implement
package main

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"strings"

// 	"server/internal/config"
// 	"server/internal/infrastructure/database/postgres"
// 	"server/pkg/logger"
// 	"server/seeds"
// 	"server/seeds/development"
// 	"server/seeds/production"
// 	"server/seeds/testing"
// )

// func main() {
// 	// Command line flags
// 	var (
// 		env     = flag.String("env", "", "Environment to seed (development, testing, production)")
// 		list    = flag.Bool("list", false, "List available seeders")
// 		seeders = flag.String("seeders", "", "Comma-separated list of specific seeders to run")
// 	)
// 	flag.Parse()

// 	// Initialize logger
// 	log := logger.NewLogger()
// 	log.Info("Starting TNP RGPV Database Seeding Tool...")

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

// 	// Get all available seeders
// 	allSeeders := map[string]seeds.Seeder{
// 		// Development seeders
// 		"dev_roles":        development.NewRolesSeeder(),
// 		"dev_admin_users":  development.NewAdminUsersSeeder(),
// 		"dev_sample_data":  development.NewSampleDataSeeder(),

// 		// Testing seeders
// 		"test_users":      testing.NewTestUsersSeeder(),
// 		"test_scenarios":  testing.NewTestScenariosSeeder(),

// 		// Production seeders
// 		"prod_roles":      production.NewInitialRolesSeeder(),
// 		"prod_settings":   production.NewDefaultSettingsSeeder(),
// 	}

// 	// List seeders if requested
// 	if *list {
// 		fmt.Println("Available seeders:")
// 		fmt.Println("Development seeders:")
// 		for name := range allSeeders {
// 			if strings.HasPrefix(name, "dev_") {
// 				fmt.Printf("  - %s\n", name)
// 			}
// 		}
// 		fmt.Println("Testing seeders:")
// 		for name := range allSeeders {
// 			if strings.HasPrefix(name, "test_") {
// 				fmt.Printf("  - %s\n", name)
// 			}
// 		}
// 		fmt.Println("Production seeders:")
// 		for name := range allSeeders {
// 			if strings.HasPrefix(name, "prod_") {
// 				fmt.Printf("  - %s\n", name)
// 			}
// 		}
// 		os.Exit(0)
// 	}

// 	// Determine which seeders to run
// 	var seedersToRun []seeds.Seeder

// 	if *seeders != "" {
// 		// Run specific seeders
// 		seederNames := strings.Split(*seeders, ",")
// 		for _, name := range seederNames {
// 			name = strings.TrimSpace(name)
// 			seeder, exists := allSeeders[name]
// 			if !exists {
// 				log.Fatal("Seeder not found", "seeder", name)
// 			}
// 			seedersToRun = append(seedersToRun, seeder)
// 		}
// 	} else if *env != "" {
// 		// Run all seeders for a specific environment
// 		prefix := ""
// 		switch *env {
// 		case "development":
// 			prefix = "dev_"
// 		case "testing":
// 			prefix = "test_"
// 		case "production":
// 			prefix = "prod_"
// 		default:
// 			log.Fatal("Invalid environment specified", "env", *env)
// 		}

// 		for name, seeder := range allSeeders {
// 			if strings.HasPrefix(name, prefix) {
// 				seedersToRun = append(seedersToRun, seeder)
// 			}
// 		}
// 	} else {
// 		fmt.Println("Please specify an environment (-env) or specific seeders (-seeders)")
// 		fmt.Println("\nUsage:")
// 		flag.PrintDefaults()
// 		os.Exit(1)
// 	}

// 	// Run the selected seeders
// 	log.Info("Running seeders...", "count", len(seedersToRun))
// 	runner := seeds.NewRunner(db, log)

// 	if err := runner.Run(seedersToRun); err != nil {
// 		log.Fatal("Failed to run seeders", "error", err)
// 	}

// 	log.Info("Database seeding completed successfully")
// }