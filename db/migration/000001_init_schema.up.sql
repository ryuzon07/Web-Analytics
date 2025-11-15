CREATE TABLE "events" (
  "id" SERIAL PRIMARY KEY,
  "site_id" varchar(100) NOT NULL,
  "event_type" varchar(50) NOT NULL,
  "path" text,
  "user_id" varchar(100),
  "timestamp" timestamptz NOT NULL
);

CREATE INDEX "idx_events_site_id_timestamp" ON "events" ("site_id", "timestamp" DESC);
CREATE INDEX "idx_events_user_id" ON "events" ("user_id");
CREATE INDEX "idx_events_site_id_event_type_timestamp" ON "events" ("site_id", "event_type", "timestamp");