package notifier

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tag string

const (
	TagWarning         Tag = "warning"
	TagTriangularFlag  Tag = "triangular_flag_on_post"
	TagRotatingLight   Tag = "rotating_light"
	TagNoEntry         Tag = "no_entry"
	TagComputer        Tag = "computer"
	TagTada            Tag = "tada"
)
type Push struct {
	Title    string
	Text     string
	Priority string
	Tags     []Tag
}

var validPriorities = map[string]bool{
	"max":    true,
	"urgent": true,
	"high": true,
	"default": true,
	"low": true,
	"min": true,
}

var validTags = map[Tag]bool{
	TagWarning:        true,
	TagTriangularFlag: true,
	TagRotatingLight:  true,
	TagNoEntry:        true,
	TagComputer:       true,
	TagTada:           true,
}
func SendPush(push Push) {
	topic := os.Getenv("NTFY_TOPIC")
	if topic == "" {
		log.Println("NTFY_TOPIC not defined")
		return
	}

	if err := validate(push); err != nil {
		log.Println("invalid push: %w", err)
	}

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://ntfy.sh/%s", topic),
   strings.NewReader(push.Text))

	req.Header.Set("Title", push.Title)
	req.Header.Set("Priority", push.Priority)
	req.Header.Set("Tags", joinTags(push.Tags))

	http.DefaultClient.Do(req)
}

func validate(push Push) error {
	if strings.TrimSpace(push.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(push.Text) == "" {
		return errors.New("text is required")
	}
	if !validPriorities[push.Priority] {
		return fmt.Errorf("invalid priority %q", push.Priority)
	}
	for _, tag := range push.Tags {
		if !validTags[tag] {
			return fmt.Errorf("invalid tag %q", tag)
		}
	}
	return nil
}

func joinTags(tags []Tag) string {
	s := make([]string, len(tags))
	for i, t := range tags {
		s[i] = string(t)
	}
	return strings.Join(s, ",")
}

// usage
	// notifier.SendPush(notifier.Push{
  //   Title:    "Alert",
  //   Text:     "Server is down",
  //   Priority: "urgent",
  //   Tags:     []notifier.Tag{notifier.TagRotatingLight},
	// })
	// return