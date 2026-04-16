# backend

## Migraciones (estilo Rails)

Este proyecto incluye un CLI en Go para ejecutar y revertir migraciones de base de datos.

Comandos disponibles:

- Ejecutar migraciones pendientes:

	go run ./cmd/migrate db:migrate

- Revertir la última migración aplicada:

	go run ./cmd/migrate db:rollback

- Revertir N migraciones:

	go run ./cmd/migrate -step=2 db:rollback

Notas:

- Los scripts de subida viven en `internal/database/migrations/sql`.
- Los scripts de rollback viven en `internal/database/migrations/sql` y terminan en `.down.sql`.
- Cada rollback debe compartir el mismo prefijo del archivo de subida (por ejemplo `001_x.sql` y `001_x.down.sql`).
