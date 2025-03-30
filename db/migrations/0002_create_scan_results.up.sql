-- Create the scanresults table
CREATE TABLE scan_results
(
    scan_id              uuid        not null default uuid_generate_v4(),
    file                 text        not null,
    url                  text,
    commit_sha           text,
    redacted_secret      text,
    raw                  text,
    detector_name        varchar,
    is_verified          bool,
    scan_completion_time timestamptz,
    created_at           timestamptz not null default current_timestamp,
    updated_at           timestamptz,
    deleted_at           timestamptz
);

CREATE UNIQUE INDEX idx_unique_scan_result ON scan_results (file, url, commit_sha, raw, detector_name);