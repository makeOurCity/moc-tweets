package tweets

import (
	"fmt"
	"time"
)

const OrionEntityType = "Tweet"
const TwitterDateTimeFormat = "Mon Jan 02 15:04:05 -0700 2006"

type OrionEntity struct {
	ID                string            `json:"id"`
	Type              string            `json:"type"`
	Body              TextAttribute     `json:"body"`
	Username          TextAttribute     `json:"username"`
	TwitterScreenName TextAttribute     `json:"twitter_screen_name"`
	TwitterID         NumberAttribute   `json:"twitter_id"`
	TweetedAt         DateTimeAttribute `json:"tweeted_at"`
}

type NumberAttribute struct {
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

func NewNumberAttribute(v int64) NumberAttribute {
	return NumberAttribute{
		Type:  "Number",
		Value: v,
	}
}

type TextAttribute struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewTextAttribute(v string) TextAttribute {
	return TextAttribute{
		Type:  "Text",
		Value: v,
	}
}

type DateTimeAttribute struct {
	Type  string    `json:"type"`
	Value time.Time `json:"value"`
}

func NewDateTimeAttribute(v time.Time) DateTimeAttribute {
	return DateTimeAttribute{
		Type:  "DateTime",
		Value: v,
	}
}

func NewDateTimeAttributeFromString(v string) (DateTimeAttribute, error) {
	t, err := time.Parse(TwitterDateTimeFormat, v)
	if err != nil {
		return DateTimeAttribute{}, fmt.Errorf("time.Parse got error: %w", err)
	}
	return NewDateTimeAttribute(t), nil
}
