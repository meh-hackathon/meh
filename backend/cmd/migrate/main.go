package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/meh-hackathon/meh/config"
	"github.com/meh-hackathon/meh/db"
	"github.com/meh-hackathon/meh/logger"
)

func main() {
	logger.New()
	logger.SetLevel(slog.LevelInfo)

	upCommand := flag.NewFlagSet("up", flag.ExitOnError)
	downCommand := flag.NewFlagSet("down", flag.ExitOnError)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	err = db.Init(config.DatabaseHost, config.DatabasePort, config.DatabaseUser, config.DatabasePassword)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "up":
		upCommand.Parse(os.Args[2:])
		handleUpCommand()
	case "down":
		downCommand.Parse(os.Args[2:])
		handleDownCommand()
	case "reset":
		handleResetCommand()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func handleResetCommand() {
	fmt.Println("Resetting database to initial state...")
	_, err := db.Exec(`DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;`)
	if err != nil {
		logger.Error("Failed to reset database", "error", err)
		os.Exit(1)
	}
	fmt.Println("Database reset successful.")
	fmt.Println("You can now run 'migrate up' to apply migrations.")
}

func handleUpCommand() {
	fmt.Println("Migrating database to latest version...")
	err := db.MigrateUp()
	if err != nil {
		logger.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	fmt.Println("Migration successful.")
}

func handleDownCommand() {
	if len(os.Args) < 3 {
		fmt.Println("Error: No target version specified")
		fmt.Println("Usage: migrate down <version>")
		fmt.Println("       Use 'migrate down 0' to revert all migrations")
		os.Exit(1)
	}

	targetVersion, err := strconv.ParseUint(os.Args[2], 10, 64)
	if err != nil {
		logger.Error("Invalid version number", "error", err)
		fmt.Println("Version must be a non-negative integer")
		os.Exit(1)
	}

	fmt.Printf("Migrating database down to version %d...\n", targetVersion)
	err = db.MigrateDown(uint(targetVersion))
	if err != nil {
		logger.Error("Migration failed", "error", err)
		os.Exit(1)
	}
	fmt.Println("Migration successful.")
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  migrate up                   - Migrate database to the latest version")
	fmt.Println("  migrate down <version>       - Downgrade database to specified version")
	fmt.Println("  migrate down 0               - Revert all migrations")
	fmt.Println("  migrate help                 - Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  migrate up")
	fmt.Println("  migrate down 3")
	fmt.Println("  migrate down 0")
}
