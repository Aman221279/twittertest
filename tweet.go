package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dghubble/oauth1"
)

// Twitter API credentials
const (
	consumerKey    = "5w6Bb4MEPLBIaV9UL3iB0uSIs"
	consumerSecret = "cypZxgpmYws8dhQhJalJQf0Pto9VUkqBKvwRfSYdvCpPTQYmaF"
	accessToken    = "1844362662827917316-c3Ebx9wNGoy75w3OsEYAHHeqPnYwh6"
	accessSecret   = "ggwoBImzbNfjhnfCRQCdZJnEtnzuYv8UhZskqBc3PxjK2"
)

func main() {
	tweetID, err := postTweet("This is my First Tweet")
	if err != nil {
		log.Fatalf("Error posting tweet: %v", err)
	}
	fmt.Printf("Posted tweet with ID: %s\n", tweetID)

	fmt.Println("After 10 seconds the file will be automatically deleted...")
	time.Sleep(10 * time.Second) // Delay for 10 seconds

	if err := deleteTweet(tweetID); err != nil {
		log.Fatalf("tweet %v", err)
	}
	fmt.Printf("Deleted tweet with ID: %s\n", tweetID)
}

// postTweet sends a tweet to Twitter
func postTweet(content string) (string, error) {
	client := getClient()
	tweetData := map[string]string{"text": content}

	tweetJSON, _ := json.Marshal(tweetData)
	response, err := client.Post("https://api.twitter.com/2/tweets", "application/json", bytes.NewBuffer(tweetJSON))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var result map[string]map[string]string
	json.NewDecoder(response.Body).Decode(&result)

	if tweetID, ok := result["data"]["id"]; ok {
		return tweetID, nil
	}
	return "", fmt.Errorf("deleted succesful")
}

// deleteTweet deletes a tweet by its ID
func deleteTweet(tweetID string) error {
	client := getClient()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("https://api.twitter.com/2/tweets/%s", tweetID), nil)

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("deleted succesfull")
	}
	return nil
}

// getClient creates an authenticated HTTP client
func getClient() *http.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	return config.Client(context.Background(), token)
}
