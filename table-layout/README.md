# Table Layout
## channels
|Name|Description|Type|Null|Deault|etc|
|:--|:--|:--|:--|:--|:--|
|channel_id|channel id|text|NG||`PRIMARY KEY`|
|name|channel name|text|NG|||
|description|channel description|text|NG|||
|thumbnail_url|thumbnail URL|text|NG||
|playlist_id|playlist ID|text|OK||
|view_count|total view count|int|NG|0||
|video_count|total video count|int|NG|0||
|subscriber_count|total subscriber count|int|NG|0||
|created_at||TIMESTAMP|NG|now()||
|updated_at||TIMESTAMP|NG|now()||

## videos
|Name|Description|Type|Null|Deault|etc|
|:--|:--|:--|:--|:--|:--|
|video_id|video id|text|NG||`PRIMARY KEY`|
|title|video name|text|NG|||
|description|video description|text|NG|||
|thumbnail_url|thumbnail URL|text|NG||
|view_count|total view count|int|NG|0||
|comment_count|total comment count|int|NG|0||
|category_id|category id|text|NG|||
|category_name|category id|text|NG|||
|channel_id|channel name|text|NG||`channels.channel_id`|
|published_at|video published date|TIMESTAMP|NG|||
|created_at||TIMESTAMP|NG|now()||
|updated_at||TIMESTAMP|NG|now()||

## comments
|Name|Description|Type|Null|Deault|etc|
|:--|:--|:--|:--|:--|:--|
|comment_id|comment id|text|NG||`PRIMARY KEY`|
|text_display|comment text|text|NG||plaintext|
|author_name|author name|text|NG|||
|author_url|author url|text|NG||not channel id|
|like_count|total like count|int|NG|0||
|reply count|total reply count|int|NG|0||
|channel_id|channel id|text|NG||`channels.channel_id`|
|video_id|video id|text|NG||`videos.video_id`|
|published_at|video published date|TIMESTAMP|NG|||
|created_at||TIMESTAMP|NG|now()||
|updated_at||TIMESTAMP|NG|now()||
