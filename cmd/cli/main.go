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

	fmt.Print("\n\n")

	epics, err := client.GetEpicsByBoard("VW")
	if err != nil {
		log.Fatal(err)
	}

	for _, issue := range epics {
		fmt.Printf("Issue: %s - %s\n", issue.Key, issue.Fields.Summary)
	}
	// issues, err := client.GetIssues(jira.IssueFilter{
	// 	Board:  "VW",
	// 	Parent: "",
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// for _, issue := range issues {
	// 	fmt.Printf("Issue: %s - %s\n", issue.Key, issue.Fields.Summary)
	// }

}
