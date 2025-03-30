-- Create the scanresults table
CREATE TABLE github_scan_results
(
    number                bigserial,
    file                  text,
    url                   text,
    commit_sha            varchar,
    raw                   text,
    detector_name         varchar,
    detector_display_name varchar,
    found_at              timestamp,
    created_at            timestamptz not null default current_timestamp,
    updated_at            timestamptz,
    deleted_at            timestamptz
);

CREATE UNIQUE INDEX idx_unique_github_scan_result ON github_scan_results (file, url, raw, detector_name, detector_display_name, number);