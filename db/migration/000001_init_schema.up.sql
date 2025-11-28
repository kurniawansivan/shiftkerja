CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL UNIQUE,
  "password_hash" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "role" varchar NOT NULL DEFAULT 'worker',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

-- Index on email for fast lookups during login
CREATE INDEX ON "users" ("email");