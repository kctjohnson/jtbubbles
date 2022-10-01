package jira

import (
	"fmt"
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

func (c *Client) GetBoardList() ([]jiraAPI.Board, error) {
	boards, _, err := c.client.Board.GetAllBoards(nil)
	if err != nil {
		return nil, err
	}
	return boards.Values, nil
}

func (c *Client) GetEpicsByBoard(boardName string) ([]jiraAPI.Issue, error) {
	epics, _, err := c.client.Issue.Search(fmt.Sprintf("issuetype=epic AND project=%s", boardName), nil)
	if err != nil {
		return nil, err
	}

	return epics, nil
}

//TODO: Maybe we want to use a pointer for the Parent filter option?
type IssueFilter struct {
	Board  string
	Parent string
}

func (c *Client) GetIssues(filter IssueFilter) ([]jiraAPI.Issue, error) {
	jql := fmt.Sprintf("project=%s", filter.Board)
	if filter.Parent != "" {
		jql = jql + fmt.Sprintf(" AND parent=%s", filter.Parent)
	}

	issues, _, err := c.client.Issue.Search(jql, nil)
	if err != nil {
		return nil, err
	}
	return issues, nil
}
