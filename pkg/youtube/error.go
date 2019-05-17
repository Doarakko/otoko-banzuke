package youtube

import "fmt"

type youtubeError struct {
	content string
	id      string
	message string
}

func (err youtubeError) Error() string {
	return fmt.Sprintf("%v %v %v", err.message, err.content, err.id)
}
