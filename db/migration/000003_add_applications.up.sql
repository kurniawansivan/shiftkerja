CREATE TABLE "applications" (
  "id" bigserial PRIMARY KEY,
  "shift_id" bigint NOT NULL,
  "worker_id" bigint NOT NULL,
  "status" varchar NOT NULL DEFAULT 'PENDING', -- PENDING, ACCEPTED, REJECTED
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "applications" ADD FOREIGN KEY ("shift_id") REFERENCES "shifts" ("id");
ALTER TABLE "applications" ADD FOREIGN KEY ("worker_id") REFERENCES "users" ("id");

-- Prevent applying to the same shift twice
CREATE UNIQUE INDEX ON "applications" ("shift_id", "worker_id");