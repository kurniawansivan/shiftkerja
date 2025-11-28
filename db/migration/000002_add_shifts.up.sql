CREATE TABLE "shifts" (
  "id" bigserial PRIMARY KEY,
  "owner_id" bigint NOT NULL, -- Who posted this?
  "title" varchar NOT NULL,
  "description" text,
  "pay_rate" decimal(10, 2) NOT NULL,
  "lat" float8 NOT NULL,
  "lng" float8 NOT NULL,
  "status" varchar NOT NULL DEFAULT 'OPEN', -- OPEN, FILLED, CANCELLED
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "shifts" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

-- Index for faster lookups by owner
CREATE INDEX ON "shifts" ("owner_id");