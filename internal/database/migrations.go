// internal/database/migrations.go
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"
)

// runMigrations ejecuta los archivos .sql que aún no han sido aplicados
func runMigrations(db *sql.DB, migrationFiles embed.FS) error {
	// 1. Crear tabla de control de migraciones si no existe
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	// 2. Leer archivos SQL embebidos
	files, err := getSQLFiles(migrationFiles)
	if err != nil {
		return err
	}

	// 3. Aplicar las que faltan
	applied := 0
	for _, filename := range files {
		ok, err := isMigrationApplied(db, filename)
		if err != nil {
			return err
		}
		if ok {
			log.Printf("   ⏭️  migración ya aplicada: %s", filename)
			continue
		}

		if err := applyMigration(db, migrationFiles, filename); err != nil {
			return fmt.Errorf("error aplicando migración '%s': %w", filename, err)
		}
		applied++
	}

	if applied == 0 {
		log.Println("✅ No hay migraciones pendientes")
	} else {
		log.Printf("✅ %d migración(es) aplicada(s)", applied)
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            filename   VARCHAR(255) PRIMARY KEY,
            applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
    `)
	return err
}

func getSQLFiles(migrationFiles embed.FS) ([]string, error) {
	var files []string

	err := fs.WalkDir(migrationFiles, "sql", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error leyendo archivos de migración: %w", err)
	}

	// Ordenar por nombre garantiza el orden 001, 002, 003...
	sort.Strings(files)
	return files, nil
}

func isMigrationApplied(db *sql.DB, filename string) (bool, error) {
	var exists bool
	err := db.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE filename = $1)`,
		filename,
	).Scan(&exists)
	return exists, err
}

func applyMigration(db *sql.DB, migrationFiles embed.FS, filename string) error {
	// Leer el contenido del archivo SQL
	content, err := migrationFiles.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error leyendo '%s': %w", filename, err)
	}

	// Ejecutar en una transacción: si falla, no queda a medias
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() // no-op si ya se hizo Commit

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("error ejecutando SQL: %w", err)
	}

	// Registrar la migración como aplicada
	if _, err := tx.Exec(
		`INSERT INTO schema_migrations (filename) VALUES ($1)`,
		filename,
	); err != nil {
		return err
	}

	return tx.Commit()
}
