CREATE TABLE comments
(
    comment_id text,
    text_display text NOT NULL,
    author_id text NOT NULL,
    author_name text NOT NULL,
    author_url text NOT NULL,
    channel_id text NOT NULL,
    video_id text NOT NULL,
    parent_id text,
    like_count int DEFAULT 0 NOT NULL,
    reply_count int DEFAULT 0 NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (comment_id)
);
