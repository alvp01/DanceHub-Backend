data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas-gorm",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/18/dev?search_path=public"

  migration {
    dir = "file://internal/database/migrations/sql"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}

env "local" {
  url = env("ATLAS_DATABASE_URL")

  migration {
    dir = "file://internal/database/migrations/sql"
  }
}
