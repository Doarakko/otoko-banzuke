CREATE TABLE videos
(
    video_id text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    thumbnail_url text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    comment_count int DEFAULT 0 NOT NULL,
    published_at TIMESTAMP NOT NULL,
    category_id text NOT NULL,
    category_name text NOT NULL,
    channel_id text NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,

    PRIMARY KEY (video_id)
);
