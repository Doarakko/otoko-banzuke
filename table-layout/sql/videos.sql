CREATE TABLE videos
(
    video_id text,
    title text NOT NULL,
    description text,
    thumbnail_url text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    like_count int DEFAULT 0 NOT NULL,
    dislike_count int DEFAULT 0 NOT NULL,
    comment_count int DEFAULT 0 NOT NULL,
    channel_id text NOT NULL,
    published_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (video_id)
);
