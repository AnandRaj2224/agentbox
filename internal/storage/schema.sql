CREATE TABLE executions(
    execution_id UUID PRIMARY KEY,
    runtime TEXT NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    exitcode INTEGER NOT NULL,
    execution_duration BIGINT NOT NULL,
    source_code TEXT NOT NULL,
    program_output TEXT NOT NULL
)
