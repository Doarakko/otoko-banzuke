# Table Layout
## channels
|Name|Description|Type|Length|Deault|Foreign Key|
|:--|:--|:--|:--|:--|:--|
|id|channel id|varchar|256|||
|name|channel name|varchar|256|||
|thumbnail_url|thumbnail URL|varchar|256|||
|create_date||string|256|||
|update_date||string|256|||

## videos
|Name|Description|Type|Length|Deault|Foreign Key|
|:--|:--|:--|:--|:--|:--|
|id|video id|varchar|256|||
|channel_id|channel id|varchar|256|||
|create_date||string|256|||
|update_date||string|256|||

## comments
|Name|Description|Type|Length|Deault|Foreign Key|
|:--|:--|:--|:--|:--|:--|
|id|comment id|varchar|256|||
|author_id|author channel id|varchar|256|||
|author_url|author channel url|varchar|256|||
|author_image_url|author profile image url|varchar|256|||
|channel_id|channel id|varchar|256|||
|video_id|video id|varchar|256|||
|parent_id|parent comment id|varchar|256|||
|like_count||int|256|||
|reply_count||int|256|||
|create_date||date|256|||
|update_date||date|256|||
|active_flag||boolean|256|||
