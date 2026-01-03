-- CreateExtension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CreateTable
CREATE TABLE "developers" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "external_id" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "team" VARCHAR(255) NOT NULL,
    "seniority" VARCHAR(50),
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "developers_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "usage_events" (
    "id" UUID NOT NULL DEFAULT gen_random_uuid(),
    "external_id" VARCHAR(255) NOT NULL,
    "developer_id" UUID NOT NULL,
    "event_type" VARCHAR(100) NOT NULL,
    "event_timestamp" TIMESTAMPTZ(6) NOT NULL,
    "lines_added" INTEGER NOT NULL DEFAULT 0,
    "lines_deleted" INTEGER NOT NULL DEFAULT 0,
    "model_used" VARCHAR(100),
    "accepted" BOOLEAN,
    "tokens_input" INTEGER NOT NULL DEFAULT 0,
    "tokens_output" INTEGER NOT NULL DEFAULT 0,
    "created_at" TIMESTAMPTZ(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "usage_events_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "developers_external_id_key" ON "developers"("external_id");

-- CreateIndex
CREATE UNIQUE INDEX "developers_email_key" ON "developers"("email");

-- CreateIndex
CREATE INDEX "idx_developers_team" ON "developers"("team");

-- CreateIndex
CREATE INDEX "idx_developers_external_id" ON "developers"("external_id");

-- CreateIndex
CREATE UNIQUE INDEX "usage_events_external_id_key" ON "usage_events"("external_id");

-- CreateIndex
CREATE INDEX "idx_events_developer" ON "usage_events"("developer_id");

-- CreateIndex
CREATE INDEX "idx_events_timestamp" ON "usage_events"("event_timestamp");

-- CreateIndex
CREATE INDEX "idx_events_type" ON "usage_events"("event_type");

-- CreateIndex
CREATE INDEX "idx_events_developer_timestamp" ON "usage_events"("developer_id", "event_timestamp");

-- AddForeignKey
ALTER TABLE "usage_events" ADD CONSTRAINT "usage_events_developer_id_fkey" FOREIGN KEY ("developer_id") REFERENCES "developers"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- CreateMaterializedView for daily_stats
CREATE MATERIALIZED VIEW daily_stats AS
SELECT
    developer_id,
    DATE(event_timestamp) as stat_date,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_shown') as suggestions_shown,
    COUNT(*) FILTER (WHERE event_type = 'cpp_suggestion_accepted') as suggestions_accepted,
    COUNT(*) FILTER (WHERE event_type = 'chat_message') as chat_interactions,
    COUNT(*) FILTER (WHERE event_type = 'cmd_k_prompt') as cmd_k_usages,
    SUM(lines_added) as total_lines_added,
    SUM(lines_deleted) as total_lines_deleted,
    SUM(lines_added) FILTER (WHERE accepted = true) as ai_lines_added
FROM usage_events
GROUP BY developer_id, DATE(event_timestamp);

-- CreateIndex for materialized view
CREATE UNIQUE INDEX idx_daily_stats_pk ON daily_stats(developer_id, stat_date);
CREATE INDEX idx_daily_stats_date ON daily_stats(stat_date);

-- Create function to refresh materialized view
CREATE OR REPLACE FUNCTION refresh_daily_stats()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY daily_stats;
END;
$$;

-- Create trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_developers_updated_at
    BEFORE UPDATE ON developers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
