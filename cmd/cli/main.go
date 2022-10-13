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
	INITIAL_LOAD ModelState = iota
	FAILED_TO_CONNECT
	PROJECT_SELECT
	ISSUE_VIEW
)

type ConnectClient jira.Client

func ConnectClientCmd() tea.Msg {
	username := os.Getenv("USERNAME")
	url := os.Getenv("JIRA_URL")
	apiToken := os.Getenv("API_TOKEN")

	var client jira.Client
	client.Connect(apiToken, username, url)
	return ConnectClient(client)
}

type FoundProject JAPI.Board

func FoundProjectCmd(foundProject JAPI.Board) tea.Cmd {
	return func() tea.Msg {
		return FoundProject(foundProject)
	}
}

type MainModel struct {
	Client            jira.Client
	cursor            int
	choice            JAPI.Board
	state             ModelState
	projectSelectView projectselect.ProjectSelectModel
}

func InitialModel() MainModel {
	return MainModel{
		Client: jira.Client{},
		cursor: 0,
		choice: JAPI.Board{},
		state:  INITIAL_LOAD,
		projectSelectView: projectselect.ProjectSelectModel{},
	}
}

func (m MainModel) Init() tea.Cmd {
	return ConnectClientCmd
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectClient:
		m.Client = jira.Client(msg)
		if m.Client.ClientValid() {
			m.state = PROJECT_SELECT
			m.projectSelectView = projectselect.NewProjectSelectModel(m.Client)
			return m, m.projectSelectView.Init()
		} else {
			m.state = FAILED_TO_CONNECT
			return m, nil
		}

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
	case INITIAL_LOAD:
		return "Connecting to Jira..."
	case FAILED_TO_CONNECT:
		return "Failed to connect to Jira! Check config settings!"
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

	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if m, err := p.StartReturningModel(); err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	} else {
		if mm, ok := m.(MainModel); ok && mm.choice.Name != "" {
			fmt.Printf("Chosen: %s\n", mm.choice.Name)
		}
	}
}
