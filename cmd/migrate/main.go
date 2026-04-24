package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/alvp01/DanceHub-Backend/internal/database"
	"github.com/joho/godotenv"
)

func main() {
	envPaths := []string{".env", "../../.env"}
	envLoaded := false
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			envLoaded = true
			break
		}
	}
	if !envLoaded {
		log.Println("No .env file found, usando variables de entorno del sistema")
	}

	step := flag.Int("step", 1, "Cantidad de migraciones a revertir con db:rollback")
	name := flag.String("name", "auto", "Nombre para db:diff")
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	cfg := database.ConfigFromEnv()
	dbURL := cfg.AtlasURL()

	if _, err := exec.LookPath("atlas"); err != nil {
		log.Fatalf("❌ Atlas CLI no está instalado o no está en PATH")
	}

	atlasArgs := []string{"migrate"}

	switch command {
	case "db:diff", "diff":
		atlasArgs = append(atlasArgs,
			"diff",
			*name,
			"--env", "gorm",
			"--config", "file://atlas.hcl",
		)
	case "db:inspect", "inspect":
		atlasArgs = []string{
			"schema",
			"inspect",
			"--env", "gorm",
			"--url", "env://src",
			"--config", "file://atlas.hcl",
		}
	case "db:migrate", "migrate":
		atlasArgs = append(atlasArgs,
			"apply",
			"--url", dbURL,
			"--dir", "file://internal/database/migrations/sql",
		)
	case "db:rollback", "rollback":
		atlasArgs = append(atlasArgs,
			"down",
			"--url", dbURL,
			"--dir", "file://internal/database/migrations/sql",
			"--amount", fmt.Sprintf("%d", *step),
		)
	default:
		fmt.Printf("Comando desconocido: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	cmd := exec.Command("atlas", atlasArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("❌ Error ejecutando Atlas: %v", err)
	}
}

func printUsage() {
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/migrate db:diff -name add_students_table")
	fmt.Println("  go run ./cmd/migrate db:inspect")
	fmt.Println("  go run ./cmd/migrate db:migrate")
	fmt.Println("  go run ./cmd/migrate db:rollback")
	fmt.Println("  go run ./cmd/migrate -step=2 db:rollback")
	fmt.Println("\nRequiere Atlas CLI: https://atlasgo.io/getting-started")
}
