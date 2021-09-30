package tweets

import (
	"fmt"
	"time"
)

const OrionEntityType = "Tweet"

type OrionEntity struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Body      TextAttribute     `json:"body"`
	Username  TextAttribute     `json:"username"`
	TwitterID TextAttribute     `json:"twitter_id"`
	TweetedAt DateTimeAttribute `json:"tweeted_at"`
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
	t, err := time.Parse("Wed Nov 04 12:25:42 +0000 2020", v)
	if err != nil {
		return DateTimeAttribute{}, fmt.Errorf("time.Parse got error: %w", err)
	}
	return NewDateTimeAttribute(t), nil
}
