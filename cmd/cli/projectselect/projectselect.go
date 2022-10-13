package projectselect

import (
	"fmt"
	"log"

	JAPI "github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kctjohnson/jtbubbles/pkg/jira"
)

type SelectProject JAPI.Board
type ProjectsLoaded []JAPI.Board

func SelectProjectCmd(foundProject JAPI.Board) tea.Cmd {
	return func() tea.Msg {
		return SelectProject(foundProject)
	}
}

func LoadProjectsCmd(client jira.Client) tea.Cmd {
	log.Printf("Creating func")
	return func() tea.Msg {
		log.Printf("In created func")
		// Get the projects associated with the account
		var projects ProjectsLoaded
		projects, err := client.GetBoardList()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Projects api call done")
		return projects
	}
}

type projectListItem struct {
	item JAPI.Board
}

func (p projectListItem) Title() string       { return fmt.Sprintf("JIRA Project ID: %d", p.item.ID) }
func (p projectListItem) Description() string { return p.item.Name }
func (p projectListItem) FilterValue() string { return p.item.Name }

type ProjectSelectModel struct {
	client   jira.Client
	cursor   int
	choice   JAPI.Board
	projects []JAPI.Board
	list     list.Model
}

func NewProjectSelectModel(client jira.Client) ProjectSelectModel {
	log.Printf("In new proj make")
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Projects"

	return ProjectSelectModel{
		client:   client,
		cursor:   0,
		choice:   JAPI.Board{},
		projects: nil,
		list:     l,
	}
}

func (m ProjectSelectModel) Init() tea.Cmd {
	return LoadProjectsCmd(m.client)
}

func (m ProjectSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case ProjectsLoaded:
		for _, project := range msg {
			m.list.InsertItem(99, projectListItem{
				item: project,
			})
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if len(m.list.Items()) > 0 {
				cmds = append(cmds, SelectProjectCmd(m.list.SelectedItem().(projectListItem).item))
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ProjectSelectModel) View() string {
	if len(m.list.Items()) > 0 {
		return m.list.View()
	} else {
		return "Loading..."
	}
}
