CREATE TABLE channels
(
    channel_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    thumbnail_url text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    video_count int DEFAULT 0 NOT NULL,
    subscriber_count int DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (channel_id)
);

-- ALTER TABLE channels DROP COLUMN playlist_id;
ALTER TABLE channels ALTER COLUMN view_count TYPE
bigint; 