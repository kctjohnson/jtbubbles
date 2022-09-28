package jira

import (
	"log"

	jiraAPI "github.com/andygrunwald/go-jira"
)

type Client struct {
	client *jiraAPI.Client
}

func (c *Client) Connect(apiToken string, username string, url string) {
	// Set up client
	authTransport := jiraAPI.BasicAuthTransport{
		Username: username,
		Password: apiToken,
	}

	client, err := jiraAPI.NewClient(authTransport.Client(), url)

	if err != nil {
		log.Print("Error connecting to jira!")
		log.Fatal(err)
	} else {
		c.client = client
	}
}
