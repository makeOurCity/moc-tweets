package main

import (
	"fmt"
	"os"

	tweets "github.com/makeourcity/moc-tweets"
)

var tc *tweets.TwitterClient
var oc *tweets.OrionClient

func init() {
	tc = tweets.NewTwitterClient(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("WTITTER_ACCESS_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)
}

func main() {
	resp, err := tc.Search("SEARCH_TEXT")
	if err != nil {
		panic(fmt.Sprintf("tc.Search got error %s", err))
	}

	list := resp.Statuses
	for _, t := range list {
		e, err := tweets.Tweet2Entity(t)
		if err != nil {
			panic(fmt.Sprintf("tweets.Tweet2Entity got error: %s", err))
		}

		fmt.Println(e)
	}
}
