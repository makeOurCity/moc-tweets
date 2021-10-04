package main

import (
	"fmt"
	"log"
	"os"

	tweets "github.com/makeourcity/moc-tweets"
)

var searchText string
var tc *tweets.TwitterClient
var oc *tweets.OrionClient

func init() {
	searchText = os.Getenv("SEARCH_TEXT")

	// Create twitter client
	tc = tweets.NewTwitterClient(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("WTITTER_ACCESS_TOKEN_SECRET"),
		os.Getenv("TWITTER_CONSUMER_KEY"),
		os.Getenv("TWITTER_CONSUMER_SECRET"),
	)

	// Create ORION client
	oc = tweets.NewOrionClient(
		os.Getenv("ORION_ENDPOINT"),
		os.Getenv("APP_CLIENT_ID"),
		os.Getenv("USER_POOL_ID"),
		os.Getenv("FIWARE_SERVICE_NAME"),
	)
}

func main() {
	// Login to Cognito
	if err := oc.Login(os.Getenv("USERNAME"), os.Getenv("PASSWORD")); err != nil {
		panic(fmt.Sprintf("oc.Login got error: %s", err))
	}

	// Search in twitter
	resp, err := tc.Search(searchText)
	if err != nil {
		panic(fmt.Sprintf("tc.Search got error %s", err))
	}

	// Check twitter search result
	if len(resp.Statuses) < 1 {
		fmt.Printf("Searched '%s' but not hit\n", searchText)
		os.Exit(0)
		return
	}

	// Loop twitter search response
	for _, t := range resp.Statuses {
		exists, err := oc.IsExistsEntity(t)
		if err != nil {
			panic(fmt.Sprintf("oc.IsExistsEntity got error: %s", err))
		}

		if exists {
			log.Printf("skip tweet(%d) which is already exists.\n", t.Id)
			continue
		}

		e, err := tweets.Tweet2Entity(t)
		if err != nil {
			panic(fmt.Sprintf("tweets.Tweet2Entity got error: %s", err))
		}
		e.SetSearchText(searchText)

		// send to ORION
		if _, err := oc.Send(*e); err != nil {
			panic(fmt.Sprintf("oc.Send got error: %s", err))
		}
	}
}
