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

// RunMigrations ejecuta los archivos .sql de la carpeta sql/ que aún no han sido aplicados.
func RunMigrations(db *sql.DB, migrationFiles embed.FS) error {
	// 1. Crear tabla de control de migraciones si no existe
	if err := createMigrationsTable(db); err != nil {
		return err
	}

	// 2. Leer archivos SQL embebidos
	files, err := getSQLFiles(migrationFiles, "sql")
	if err != nil {
		return err
	}

	// 3. Aplicar las que faltan
	applied := 0
	for _, filename := range files {
		if strings.HasSuffix(filename, ".down.sql") {
			continue
		}

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

// RollbackMigrations revierte las últimas N migraciones aplicadas.
func RollbackMigrations(db *sql.DB, migrationFiles embed.FS, steps int) error {
	if steps <= 0 {
		return fmt.Errorf("steps debe ser mayor a 0")
	}

	if err := createMigrationsTable(db); err != nil {
		return err
	}

	rows, err := db.Query(
		`SELECT filename FROM schema_migrations ORDER BY filename DESC LIMIT $1`,
		steps,
	)
	if err != nil {
		return fmt.Errorf("error obteniendo migraciones aplicadas: %w", err)
	}
	defer rows.Close()

	var applied []string
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return fmt.Errorf("error leyendo migración aplicada: %w", err)
		}
		applied = append(applied, filename)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterando migraciones aplicadas: %w", err)
	}

	if len(applied) == 0 {
		log.Println("✅ No hay migraciones para revertir")
		return nil
	}

	rolledBack := 0
	for _, filename := range applied {
		downFilename := downMigrationFileFromUp(filename)
		if err := rollbackMigration(db, migrationFiles, filename, downFilename); err != nil {
			return fmt.Errorf("error revirtiendo migración '%s': %w", filename, err)
		}
		rolledBack++
	}

	log.Printf("✅ %d migración(es) revertida(s)", rolledBack)
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

func getSQLFiles(migrationFiles embed.FS, dir string) ([]string, error) {
	var files []string

	err := fs.WalkDir(migrationFiles, dir, func(path string, d fs.DirEntry, err error) error {
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

func rollbackMigration(db *sql.DB, migrationFiles embed.FS, upFilename, downFilename string) error {
	content, err := migrationFiles.ReadFile(downFilename)
	if err != nil {
		return fmt.Errorf("no se encontró la migración de rollback '%s': %w", downFilename, err)
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("error ejecutando SQL de rollback: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM schema_migrations WHERE filename = $1`, upFilename); err != nil {
		return fmt.Errorf("error eliminando registro de migración aplicada: %w", err)
	}

	return tx.Commit()
}

func downMigrationFileFromUp(upFilename string) string {
	if strings.HasSuffix(upFilename, ".down.sql") {
		return upFilename
	}

	if strings.HasSuffix(upFilename, ".sql") {
		return strings.TrimSuffix(upFilename, ".sql") + ".down.sql"
	}

	return upFilename + ".down.sql"
}
