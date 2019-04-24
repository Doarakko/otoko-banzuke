CREATE TABLE channels
(
    channel_id text NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    thumbnail_url text NOT NULL,
    view_count int DEFAULT 0,
    video_count int DEFAULT 0,
    subscriber_count int DEFAULT 0,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (channel_id)
);

-- ALTER TABLE channels DROP COLUMN playlist_id;