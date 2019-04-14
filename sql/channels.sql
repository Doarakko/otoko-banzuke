CREATE TABLE channels
(
    id text,
    name text NOT NULL,
    description text,
    thumbnail_url text NOT NULL,
    playlist_id text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    video_count int DEFAULT 0 NOT NULL,
    -- comment_count int DEFAULT 0 NOT NULL,
    subscriber_count int DEFAULT 0 NOT NULL,
    create_date TIMESTAMP DEFAULT now() NOT NULL,
    update_date TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);

