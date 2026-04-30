-- Create "academies" table
CREATE TABLE "academies" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "primary_phone" text NOT NULL,
  "secondary_phone" text NULL,
  "password_hash" text NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_academies_email" to table: "academies"
CREATE UNIQUE INDEX "idx_academies_email" ON "academies" ("email");
-- Create index "idx_academies_name" to table: "academies"
CREATE UNIQUE INDEX "idx_academies_name" ON "academies" ("name");
-- Create index "idx_academies_primary_phone" to table: "academies"
CREATE UNIQUE INDEX "idx_academies_primary_phone" ON "academies" ("primary_phone");
-- Create index "idx_academies_secondary_phone" to table: "academies"
CREATE UNIQUE INDEX "idx_academies_secondary_phone" ON "academies" ("secondary_phone");
-- Create "refresh_tokens" table
CREATE TABLE "refresh_tokens" (
  "id" uuid NOT NULL,
  "academy_id" uuid NOT NULL,
  "token_hash" text NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  "revoked_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_refresh_tokens_academy_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_academy_id" ON "refresh_tokens" ("academy_id");
-- Create index "idx_refresh_tokens_token_hash" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token_hash" ON "refresh_tokens" ("token_hash");
-- Create "students" table
CREATE TABLE "students" (
  "id" uuid NOT NULL,
  "academy_id" uuid NOT NULL,
  "name" text NOT NULL,
  "last_name" text NOT NULL,
  "email" text NOT NULL,
  "phone" text NOT NULL,
  "id_document" text NOT NULL,
  "birth_date" text NOT NULL,
  "address" text NOT NULL,
  "allergies" text NULL,
  "pathologies" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_students_academy_id" to table: "students"
CREATE INDEX "idx_students_academy_id" ON "students" ("academy_id");
-- Create index "idx_students_email" to table: "students"
CREATE UNIQUE INDEX "idx_students_email" ON "students" ("email");
-- Create index "idx_students_id_document" to table: "students"
CREATE UNIQUE INDEX "idx_students_id_document" ON "students" ("id_document");
-- Create index "idx_students_phone" to table: "students"
CREATE UNIQUE INDEX "idx_students_phone" ON "students" ("phone");
