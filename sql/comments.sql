CREATE TABLE comments (
    id text,
    text_display text NOT NULL,
    author_id text NOT NULL,
    author_name text NOT NULL,
    author_url text NOT NULL,
    channel_id text NOT NULL,
    video_id text NOT NULL,
    parent_id text,
    like_count int DEFAULT 0 NOT NULL,
    reply_count int DEFAULT 0 NOT NULL,
    create_date TIMESTAMP DEFAULT now() NOT NULL,
    update_date TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);

