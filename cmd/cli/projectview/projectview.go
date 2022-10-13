package projectview

import (
	JAPI "github.com/andygrunwald/go-jira"
	"github.com/kctjohnson/jtbubbles/pkg/jira"
)

type ProjectViewModel struct {
	client  jira.Client
	project JAPI.Board
	issues  []JAPI.Issue
	epics   []JAPI.Epic
}

func NewProjectViewModel(client jira.Client, project JAPI.Board) ProjectViewModel {
	client.GetIssues(jira.IssueFilter{
		Board: project.Name,
	})
	// boards, err := client.GetBoardList()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// epics, err := client.GetEpicsByBoard("VW")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	return ProjectViewModel{}
}
