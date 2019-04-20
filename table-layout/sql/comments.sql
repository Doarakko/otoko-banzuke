CREATE TABLE comments
(
    comment_id text NOT NULL,
    text_display text NOT NULL,
    author_name text NOT NULL,
    author_url text NOT NULL,
    like_count int DEFAULT 0 NOT NULL,
    reply_count int DEFAULT 0 NOT NULL,
    channel_id text NOT NULL,
    video_id text NOT NULL,
    created_at TIMESTAMP DEFAULT now() NOT NULL,
    updated_at TIMESTAMP DEFAULT now() NOT NULL,

    PRIMARY KEY (comment_id)
);
