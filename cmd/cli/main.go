package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/kctjohnson/jtbubbles/pkg/jira"
)

func Execute() {
	username := os.Getenv("USERNAME")
	url := os.Getenv("JIRA_URL")
	apiToken := os.Getenv("API_TOKEN")

	var client jira.Client
	client.Connect(apiToken, username, url)
	boards, err := client.GetBoardList()
	if err != nil {
		log.Fatal(err)
	}

	for _, board := range boards {
		fmt.Printf("Board: %s\n", board.Name)
	}
}
