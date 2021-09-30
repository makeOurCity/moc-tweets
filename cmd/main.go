package main

import (
	"encoding/json"
	"fmt"
	"os"

	tweets "github.com/makeourcity/moc-tweets"
)

var searchText string
var tc *tweets.TwitterClient
var oc *tweets.OrionClient

func init() {
	searchText = os.Getenv("SEARCH_TEXT")

	tc = tweets.NewTwitterClient(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("WTITTER_ACCESS_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)
}

func main() {
	resp, err := tc.Search(searchText)
	if err != nil {
		panic(fmt.Sprintf("tc.Search got error %s", err))
	}

	if len(resp.Statuses) < 1 {
		fmt.Printf("Searched '%s' but not hit\n", searchText)
		os.Exit(0)
		return
	}
	for _, t := range resp.Statuses {
		e, err := tweets.Tweet2Entity(t)
		if err != nil {
			panic(fmt.Sprintf("tweets.Tweet2Entity got error: %s", err))
		}

		b, _ := json.Marshal(e)
		fmt.Println(string(b))
	}
}
