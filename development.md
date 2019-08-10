# Development Guide
Run on local using Heroku PostgreSQL.

## Requirements
- Golang
- Youtube Data API
- [Heroku](https://www.heroku.com) account
- [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli)
- Credit card
    - It does not take money, to use Heroku addons.

## Setup
1. Clone code
```
$ git clone https://github.com/Doarakko/otoko-banzuke
```
2. Remove comment out
- `main.go` and `routine/main.go`
```
package main

import (
	"log"
	"github.com/joho/godotenv"

	// omission
)

// omission

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// omission
```

3. Enter your Youtube Data API key and Heroku postgres URL
```
YOUTUBE_API_KEY = xyz
DATABASE_URL = postgres://abcde
```
```
$ mv .env.example .env
```

4. Create table on Heroku.
```
$ heroku pg:psql --app <enter your heroku app name> < table/layout/sql/channels.sql
$ heroku pg:psql --app <enter your heroku app name> < table/layout/sql/videos.sql
$ heroku pg:psql --app <enter your heroku app name> < table/layout/sql/comments.sql
```

5. Build
```
$ go build -o otoko-banzuke
```

6. Run on local
```
$ ./otoko-banzuke
```
Access to `http://localhost:8080`, and commend Otoko.

7. Search Otoko and insert comment
```
$ cd routine
$ go run routine.go
```
If you run on Heroku, set this program to scheduler.


## Hints
- `routine/main.go`

if commended channels have many video, exceeds YouTube API limit.
So default is comment out `searchAllComments()`.
```
func main() {
	// omission
	// searchAllComments()
	searchNewComments()
}
```

- `pkg/youtube/comment.go`

If you want to search other comments, change this regular expression.
```
var re = regexp.MustCompile("^.+(漢|漢達|男|男達|男性|男子|おとこ|オトコ|女|女達|女性|女子|おんな|オンナ)(。|\\.|~|〜|!|！|\\*|＊|w|W|♂|♀){0,1}$")

// CheckComment if otoko comment return true
func (c *Comment) CheckComment() bool {
	return re.MatchString(c.TextDisplay) && c.LikeCount >= 5
}
```