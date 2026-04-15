package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/alvp01/DanceHub-Backend/internal/database"
	"github.com/alvp01/DanceHub-Backend/internal/database/migrations"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, usando variables de entorno del sistema")
	}

	step := flag.Int("step", 1, "Cantidad de migraciones a revertir con db:rollback")
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	cfg := database.ConfigFromEnv()
	db, err := database.OpenDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Error conectando a la base de datos: %v", err)
	}
	defer db.Close()

	switch command {
	case "db:migrate", "migrate":
		if err := database.RunMigrations(db, migrations.Files); err != nil {
			log.Fatalf("❌ Error ejecutando migraciones: %v", err)
		}
	case "db:rollback", "rollback":
		if err := database.RollbackMigrations(db, migrations.Files, *step); err != nil {
			log.Fatalf("❌ Error revirtiendo migraciones: %v", err)
		}
	default:
		fmt.Printf("Comando desconocido: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Uso:")
	fmt.Println("  go run ./cmd/migrate db:migrate")
	fmt.Println("  go run ./cmd/migrate db:rollback")
	fmt.Println("  go run ./cmd/migrate -step=2 db:rollback")
}
