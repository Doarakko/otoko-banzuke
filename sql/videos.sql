CREATE TABLE videos
(
    id text,
    name text NOT NULL,
    description text,
    thumbnail_url text NOT NULL,
    view_count int DEFAULT 0 NOT NULL,
    like_count int DEFAULT 0 NOT NULL,
    dislike_count int DEFAULT 0 NOT NULL,
    comment_count int DEFAULT 0 NOT NULL,
    channel_id text NOT NULL,
    upload_date TIMESTAMP NOT NULL,
    create_date TIMESTAMP DEFAULT now() NOT NULL,
    update_date TIMESTAMP DEFAULT now() NOT NULL,
    PRIMARY KEY (id)
);

