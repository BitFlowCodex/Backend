package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	Payload struct {
		Commits []struct{} `json:"commits"`
	} `json:"payload"`
}

func check(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: github-activity <github username>")
		return
	}

	username := os.Args[1]
	url := "https://api.github.com/users/" + username + "/events"

	resp, err := http.Get(url)
	check(err,"Error fetching data:")
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("User not found or API error")
		return
	}

	var events []Event
	err = json.NewDecoder(resp.Body).Decode(&events)
	check(err, "Error decoding JSON:")

	for _, e := range events {
		switch e.Type {
			case "PushEvent":
				fmt.Printf("- Pushed %d commits to %s \n", len(e.Payload.Commits), e.Repo.Name)
			case "WatchEvent":
				fmt.Printf("- Starred %s \n", e.Repo.Name)
			default:
				fmt.Printf("- %s in %s \n", e.Type, e.Repo.Name)
		}
	}
}
