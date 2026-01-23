package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/vondi-global/lib/logger"
	libmigrator "github.com/vondi-global/lib/migrator"

	"github.com/vondi-global/listings/internal/app/migrator"
	"github.com/vondi-global/listings/internal/config"
)

type flags struct {
	command        string
	withFixtures   bool
	onlyFixtures   bool
	forceVersion   int
	forceVersionSet bool // tracks if -version flag was explicitly provided
	createName     string
	isFixture      bool
	showVersion    bool
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	_ = godotenv.Load()

	f := parseFlags()

	if f.showVersion {
		fmt.Println("listings-migrator")
		return nil
	}

	if f.command == "" {
		printUsage()
		return fmt.Errorf("command is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logger.Init(cfg.App.Env, cfg.App.LogLevel, "", false, false)

	db, err := connectDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	m := migrator.New(db, cfg)

	return executeCommand(m, f)
}

func parseFlags() flags {
	var f flags
	flag.StringVar(&f.command, "command", "", "Command: up, down, down-all, status, force, create")
	flag.BoolVar(&f.withFixtures, "with-fixtures", false, "Run fixtures after migrations (for up)")
	flag.BoolVar(&f.onlyFixtures, "only-fixtures", false, "Run only fixtures (for up, down, down-all)")
	flag.IntVar(&f.forceVersion, "version", -1, "Version number (for force command)")
	flag.StringVar(&f.createName, "name", "", "Migration name (for create command)")
	flag.BoolVar(&f.isFixture, "fixture", false, "Create fixture instead of migration (for create)")
	flag.BoolVar(&f.showVersion, "v", false, "Show version")
	flag.Parse()

	// Track if -version was explicitly provided (changed from default -1)
	f.forceVersionSet = f.forceVersion != -1

	if f.command == "" && flag.NArg() > 0 {
		f.command = flag.Arg(0)
	}

	return f
}

func connectDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func executeCommand(m *migrator.Migrator, f flags) error {
	switch f.command {
	case "up":
		return executeUp(m, f)
	case "down":
		return executeDown(m, f)
	case "down-all":
		return executeDownAll(m, f)
	case "status":
		return executeStatus(m)
	case "force":
		return executeForce(m, f)
	case "create":
		return executeCreate(m, f)
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", f.command)
	}
}

func executeUp(m *migrator.Migrator, f flags) error {
	switch {
	case f.withFixtures:
		if err := m.UpWithFixtures(); err != nil {
			return err
		}
		fmt.Println("All migrations and fixtures applied successfully")
	case f.onlyFixtures:
		if err := m.UpOnlyFixtures(); err != nil {
			return err
		}
		fmt.Println("All fixtures applied successfully")
	default:
		if err := m.Up(); err != nil {
			return err
		}
		fmt.Println("All migrations applied successfully")
	}
	return nil
}

func executeDown(m *migrator.Migrator, f flags) error {
	if f.onlyFixtures {
		if err := m.DownOnlyFixtures(); err != nil {
			return err
		}
		fmt.Println("Last fixture rolled back successfully")
	} else {
		if err := m.Down(); err != nil {
			return err
		}
		fmt.Println("Last migration rolled back successfully")
	}
	return nil
}

func executeDownAll(m *migrator.Migrator, f flags) error {
	if f.onlyFixtures {
		if err := m.DownAllOnlyFixtures(); err != nil {
			return err
		}
		fmt.Println("All fixtures rolled back successfully")
	} else {
		if err := m.DownAll(); err != nil {
			return err
		}
		fmt.Println("All migrations rolled back successfully")
	}
	return nil
}

func executeStatus(m *migrator.Migrator) error {
	status, err := m.Status()
	if err != nil {
		return err
	}
	fmt.Println("Migrations:")
	printStatus(status)

	fStatus, err := m.FixturesStatus()
	if err == nil && fStatus.TotalCount > 0 {
		fmt.Println()
		fmt.Println("Fixtures:")
		printStatus(fStatus)
	}
	return nil
}

func executeForce(m *migrator.Migrator, f flags) error {
	if !f.forceVersionSet {
		return fmt.Errorf("-version is required for force command. Usage: migrator force -version N")
	}
	if f.forceVersion < 0 {
		return fmt.Errorf("-version must be >= 0")
	}
	if err := m.Force(f.forceVersion); err != nil {
		return err
	}
	fmt.Printf("Version forced to %d\n", f.forceVersion)
	return nil
}

func executeCreate(m *migrator.Migrator, f flags) error {
	if f.createName == "" {
		return fmt.Errorf("-name is required for create command")
	}
	if f.isFixture {
		if err := m.CreateFixture(f.createName); err != nil {
			return err
		}
		fmt.Printf("Fixture created: %s\n", f.createName)
	} else {
		if err := m.Create(f.createName); err != nil {
			return err
		}
		fmt.Printf("Migration created: %s\n", f.createName)
	}
	return nil
}

func printUsage() {
	fmt.Println("Listings Service Migrator")
	fmt.Println()
	fmt.Println("Usage: migrator <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up           Apply all pending migrations")
	fmt.Println("  down         Rollback last migration")
	fmt.Println("  down-all     Rollback all migrations")
	fmt.Println("  status       Show migration status")
	fmt.Println("  force        Force set version (for dirty state fix)")
	fmt.Println("  create       Create new migration files")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -with-fixtures   Run fixtures after migrations (for up)")
	fmt.Println("  -only-fixtures   Run only fixtures (for up, down, down-all)")
	fmt.Println("  -version N       Version number (for force)")
	fmt.Println("  -name NAME       Migration name (for create)")
	fmt.Println("  -fixture         Create fixture instead of migration (for create)")
	fmt.Println("  -v               Show version")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  migrator up")
	fmt.Println("  migrator up -with-fixtures")
	fmt.Println("  migrator up -only-fixtures")
	fmt.Println("  migrator down")
	fmt.Println("  migrator down -only-fixtures")
	fmt.Println("  migrator down-all")
	fmt.Println("  migrator status")
	fmt.Println("  migrator force -version 5")
	fmt.Println("  migrator create -name add_users_table")
	fmt.Println("  migrator create -name seed_test_data -fixture")
	fmt.Println()
	fmt.Println("Environment variables (via .env or VONDILISTINGS_ prefix):")
	fmt.Println("  VONDILISTINGS_DB_HOST, VONDILISTINGS_DB_PORT, etc.")
}

func printStatus(status libmigrator.MigrationStatus) {
	fmt.Printf("  Current version: %d\n", status.CurrentVersion)
	if status.Dirty {
		fmt.Printf("  Status:          DIRTY (fix with: migrator force -version N)\n")
	} else {
		fmt.Printf("  Status:          clean\n")
	}
	fmt.Printf("  Applied:         %d\n", status.AppliedCount)
	fmt.Printf("  Pending:         %d\n", status.PendingCount)
	fmt.Printf("  Total:           %d\n", status.TotalCount)
}
