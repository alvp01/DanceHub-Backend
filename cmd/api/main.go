// cmd/api/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/alvp01/DanceHub-Backend/internal/academy"
	"github.com/alvp01/DanceHub-Backend/internal/database"
	"github.com/alvp01/DanceHub-Backend/internal/database/migrations"
	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, usando variables de entorno del sistema")
	}

	cfg := database.ConfigFromEnv()
	db, err := database.Init(cfg, migrations.Files)
	if err != nil {
		log.Fatalf("❌ Error inicializando base de datos: %v", err)
	}
	defer db.Close()

	jwtManager, err := jwtpkg.NewManager()
	if err != nil {
		log.Fatalf("❌ Error inicializando JWT: %v", err)
	}

	academyRepo := academy.NewRepository(db)
	academyService := academy.NewService(academyRepo, jwtManager)
	academyHandler := academy.NewHandler(academyService)
	authMiddleware := middleware.Auth(jwtManager)

	mux := http.NewServeMux()
	academyHandler.RegisterRoutes(mux, authMiddleware)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
