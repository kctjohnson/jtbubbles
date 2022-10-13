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
	return func() tea.Msg {
		// Get the projects associated with the account
		var projects ProjectsLoaded
		projects, err := client.GetBoardList()
		if err != nil {
			log.Fatal(err)
		}
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
	width    int
	height   int
}

func NewProjectSelectModel(client jira.Client) ProjectSelectModel {
	return ProjectSelectModel{
		client:   client,
		cursor:   0,
		choice:   JAPI.Board{},
		projects: nil,
		list:     list.Model{},
	}
}

func (m ProjectSelectModel) Init() tea.Cmd {
	return LoadProjectsCmd(m.client)
}

func (m ProjectSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case ProjectsLoaded:
		l := list.New([]list.Item{}, list.NewDefaultDelegate(), m.width, m.height)
		l.Title = "Projects"
		m.list = l
		for _, project := range msg {
			m.list.InsertItem(0, projectListItem{
				item: project,
			})
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		if len(m.list.Items()) > 0 {
			m.list.SetSize(msg.Width, msg.Height)
		}
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

	if len(m.list.Items()) > 0 {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ProjectSelectModel) View() string {
	if len(m.list.Items()) > 0 {
		return m.list.View()
	} else {
		return "Loading..."
	}
}
