// internal/database/postgres.go
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Config agrupa toda la configuración de la DB
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func ConfigFromEnv() Config {
	return Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
}

// dsn construye el string de conexión
// si dbName está vacío, conecta al servidor sin especificar DB (para crearla)
func (c Config) dsn(dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, dbName, c.SSLMode,
	)
}

// Init es el punto de entrada: garantiza que la DB y las tablas existen
func Init(cfg Config, migrationFiles embed.FS) (*sql.DB, error) {
	// 1. Esperar a que Postgres esté listo (crítico en Docker)
	if err := waitForPostgres(cfg); err != nil {
		return nil, err
	}

	// 2. Crear la DB si no existe
	if err := ensureDatabaseExists(cfg); err != nil {
		return nil, err
	}

	// 3. Conectar a la DB correcta
	db, err := connect(cfg.dsn(cfg.Name))
	if err != nil {
		return nil, fmt.Errorf("error conectando a %s: %w", cfg.Name, err)
	}

	// 4. Ejecutar migraciones pendientes
	if err := runMigrations(db, migrationFiles); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// waitForPostgres reintenta la conexión hasta que Postgres acepte conexiones
// Necesario porque Docker puede tardar en levantar el servicio
func waitForPostgres(cfg Config) error {
	const (
		maxRetries = 10
		retryDelay = 2 * time.Second
	)

	// Conectar a la DB de sistema "postgres" para verificar disponibilidad
	dsnSystem := cfg.dsn("postgres")

	log.Println("⏳ Esperando a que PostgreSQL esté disponible...")

	for i := 1; i <= maxRetries; i++ {
		db, err := sql.Open("postgres", dsnSystem)
		if err == nil {
			err = db.Ping()
			db.Close()
		}

		if err == nil {
			log.Println("✅ PostgreSQL disponible")
			return nil
		}

		log.Printf("   intento %d/%d fallido: %v", i, maxRetries, err)

		if i < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	return fmt.Errorf("postgresql no disponible después de %d intentos", maxRetries)
}

// ensureDatabaseExists crea la DB si no existe
func ensureDatabaseExists(cfg Config) error {
	// Conectar al servidor sin especificar la DB objetivo
	db, err := connect(cfg.dsn("postgres"))
	if err != nil {
		return fmt.Errorf("error conectando al servidor postgres: %w", err)
	}
	defer db.Close()

	// Verificar si la DB existe
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)`
	if err := db.QueryRow(query, cfg.Name).Scan(&exists); err != nil {
		return fmt.Errorf("error verificando existencia de DB: %w", err)
	}

	if exists {
		log.Printf("✅ Base de datos '%s' ya existe", cfg.Name)
		return nil
	}

	// Crear la DB (no se puede parametrizar el nombre con $1 en DDL)
	// cfg.Name viene de variables de entorno, no de input de usuario → seguro
	_, err = db.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, cfg.Name))
	if err != nil {
		return fmt.Errorf("error creando base de datos '%s': %w", cfg.Name, err)
	}

	log.Printf("🎉 Base de datos '%s' creada exitosamente", cfg.Name)
	return nil
}

// connect abre y verifica una conexión con pool configurado
func connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
