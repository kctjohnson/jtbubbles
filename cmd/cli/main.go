package cli

import (
	"fmt"
	"log"
	"os"

	JAPI "github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kctjohnson/jtbubbles/cmd/cli/projectselect"
	"github.com/kctjohnson/jtbubbles/pkg/jira"
)

type ModelState int

const (
	PROJECT_SELECT ModelState = iota
	ISSUE_VIEW
)

type FoundProject JAPI.Board
func FoundProjectCmd(foundProject JAPI.Board) tea.Cmd {
	return func() tea.Msg {
		return FoundProject(foundProject)
	}
}

type MainModel struct {
	Projects          []JAPI.Board
	Epics             []JAPI.Issue
	Client            jira.Client
	cursor            int
	choice            JAPI.Board
	state             ModelState
	projectSelectView projectselect.ProjectSelectModel
	//issueView         IssueModel
}

func InitialModel() MainModel {
	username := os.Getenv("USERNAME")
	url := os.Getenv("JIRA_URL")
	apiToken := os.Getenv("API_TOKEN")

	var client jira.Client
	client.Connect(apiToken, username, url)
	boards, err := client.GetBoardList()
	if err != nil {
		log.Fatal(err)
	}

	epics, err := client.GetEpicsByBoard("VW")
	if err != nil {
		log.Fatal(err)
	}

	return MainModel{
		Projects:          boards,
		Epics:             epics,
		Client:            client,
		cursor:            0,
		choice:            JAPI.Board{},
		state:             PROJECT_SELECT,
		projectSelectView: projectselect.NewProjectSelectModel(client),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case projectselect.SelectProject:
		log.Printf("Selected Project: %s\n", msg.Name)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	switch m.state {
	case PROJECT_SELECT:
		newProjectSelect, cmd := m.projectSelectView.Update(msg)
		m.projectSelectView = newProjectSelect.(projectselect.ProjectSelectModel)
		return m, cmd
	case ISSUE_VIEW:
		return m, nil
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case PROJECT_SELECT:
		return m.projectSelectView.View()
	case ISSUE_VIEW:
		return "Issue view!"
	}
	return ""
}

func Execute() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(InitialModel())
	if m, err := p.StartReturningModel(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	} else {
		if mm, ok := m.(MainModel); ok && mm.choice.Name != "" {
			fmt.Printf("Chosen: %s\n", mm.choice.Name)
		}
	}
}
