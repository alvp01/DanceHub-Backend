// cmd/api/main.go
package main

import (
	"log"
	"os"

	"github.com/alvp01/DanceHub-Backend/internal/academy"
	"github.com/alvp01/DanceHub-Backend/internal/database"
	jwtpkg "github.com/alvp01/DanceHub-Backend/internal/jwt"
	"github.com/alvp01/DanceHub-Backend/internal/middleware"
	"github.com/alvp01/DanceHub-Backend/internal/student"
	"github.com/gin-gonic/gin"
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

	cfg := database.ConfigFromEnv()
	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("❌ Error inicializando base de datos: %v", err)
	}
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	jwtManager, err := jwtpkg.NewManager()
	if err != nil {
		log.Fatalf("❌ Error inicializando JWT: %v", err)
	}

	academyRepo := academy.NewRepository(db)
	academyService := academy.NewService(academyRepo, jwtManager)
	academyHandler := academy.NewHandler(academyService)
	studentRepo := student.NewRepository(db)
	studentService := student.NewService(studentRepo)
	studentHandler := student.NewHandler(studentService)
	authMiddleware := middleware.Auth(jwtManager)

	router := gin.Default()
	academyHandler.RegisterRoutes(router, authMiddleware)
	studentHandler.RegisterRoutes(router, authMiddleware)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor corriendo en :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
