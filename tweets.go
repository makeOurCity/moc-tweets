package tweets

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
)

type TwitterClient struct {
	api *anaconda.TwitterApi
}

func NewTwitterClient(accessToken, accessTokenSecret, consumerKey, consumerSecret string) *TwitterClient {
	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)

	c := TwitterClient{
		api: api,
	}

	return &c
}

func (t *TwitterClient) Search(word string) (anaconda.SearchResponse, error) {
	v := url.Values{}
	res, err := t.api.GetSearch(word, v)
	if err != nil {
		return anaconda.SearchResponse{}, fmt.Errorf("api.GetSearch got error: %w", err)
	}

	return res, nil
}

func Tweet2Entity(t anaconda.Tweet) (*OrionEntity, error) {
	e := OrionEntity{
		Type: OrionEntityType,
	}

	e.ID = fmt.Sprintf("urn:ngsi-ld:%s:%d", OrionEntityType, t.Id)
	e.Body = NewTextAttribute(t.FullText)
	e.Username = NewTextAttribute(t.User.Name)
	e.TwitterID = NewNumberAttribute(t.User.Id)
	e.TwitterScreenName = NewTextAttribute(t.User.ScreenName)

	attr, err := NewDateTimeAttributeFromString(t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("NewDateTimeAttributeFromString got error: %w", err)
	}
	e.TweetedAt = attr

	return &e, nil
}
