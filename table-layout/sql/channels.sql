CREATE TABLE channels
(
    channel_id text,
    name text NOT NULL,
    description text,
    thumbnail_url text NOT NULL,
    playlist_id text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    video_count int DEFAULT 0 NOT NULL,
    subscriber_count int DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (channel_id)
);

