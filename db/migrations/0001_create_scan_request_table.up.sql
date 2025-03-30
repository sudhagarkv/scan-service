CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- Create the scan_requests table
CREATE TABLE scan_requests
(
    id              uuid                 default uuid_generate_v4(),
    url             text        not null,
    is_private      bool,
    encrypted_token text,
    scan_id         uuid,
    queue_status    text                 default 'NotQueued',
    created_at      timestamptz not null default current_timestamp,
    modified_at     timestamptz,
    deleted_at      timestamptz
);
