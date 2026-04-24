# backend

## Base de datos

Este backend usa:

- GORM para operaciones de base de datos.
- Atlas para aplicar y revertir migraciones SQL.
- ariga.io/atlas-provider-gorm para generar migraciones desde los modelos GORM.

## Migraciones con Atlas

Instala Atlas CLI:

	https://atlasgo.io/getting-started

Define la variable con tu conexión:

	export ATLAS_DATABASE_URL="postgres://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable"

Inspeccionar el esquema deseado cargado desde modelos GORM:

	go run ./cmd/migrate db:inspect

Generar una nueva migración por diff (desde modelos GORM):

	go run ./cmd/migrate db:diff -name add_students_table

Aplicar migraciones pendientes:

	atlas migrate apply --env local --config file://atlas.hcl

Revertir la última migración:

	atlas migrate down --env local --config file://atlas.hcl --amount 1

También puedes usar el wrapper en Go:

	go run ./cmd/migrate db:diff -name add_students_table
	go run ./cmd/migrate db:inspect
	go run ./cmd/migrate db:migrate
	go run ./cmd/migrate db:rollback
	go run ./cmd/migrate -step=2 db:rollback

Notas:

- Los scripts de subida viven en `internal/database/migrations/sql`.
- Los scripts de rollback viven en `internal/database/migrations/sql` y terminan en `.down.sql`.
- Cada rollback debe compartir el mismo prefijo del archivo de subida (por ejemplo `001_x.sql` y `001_x.down.sql`).
