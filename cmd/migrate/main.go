package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	up := flag.Bool("up", false, "Apply all up migrations")
	down := flag.Bool("down", false, "Apply down migrations (drops everything in last migration)")
	steps := flag.Int("steps", 0, "Number of steps to migrate (positive to up, negative to down)")
	flag.Parse()

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("unable to get cwd:", err)
		os.Exit(1)
	}
	migrationsPath := filepath.ToSlash(filepath.Join(cwd, "migrations", "sql"))
	sourceURL := "file://" + migrationsPath

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("DATABASE_URL not set")
		os.Exit(1)
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		fmt.Println("failed to create migrate instance:", err)
		os.Exit(1)
	}

	if *up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("migrate up error:", err)
			os.Exit(1)
		}
		fmt.Println("migrations applied (up)")
		return
	}

	if *down {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("migrate down error:", err)
			os.Exit(1)
		}
		fmt.Println("migrations rolled back (down)")
		return
	}

	if *steps != 0 {
		if *steps > 0 {
			if err := m.Steps(*steps); err != nil {
				fmt.Println("migrate steps up error:", err)
				os.Exit(1)
			}
			fmt.Println("migrations applied (steps up)")
			return
		}
		if err := m.Steps(*steps); err != nil {
			fmt.Println("migrate steps error:", err)
			os.Exit(1)
		}
		fmt.Println("migrations applied (steps)")
		return
	}

	fmt.Println("No action specified. Use -up, -down, or -steps")
}
