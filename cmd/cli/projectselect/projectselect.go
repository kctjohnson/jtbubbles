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

func SelectProjectCmd(foundProject JAPI.Board) tea.Cmd {
	return func() tea.Msg {
		return SelectProject(foundProject)
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
	// Get the projects associated with the account
	projects, err := client.GetBoardList()
	if err != nil {
		log.Fatal(err)
	}

	items := make([]list.Item, len(projects))
	for i, project := range projects {
		items[i] = projectListItem{
			item: project,
		}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Projects"

	return ProjectSelectModel{
		client:   client,
		cursor:   0,
		choice:   JAPI.Board{},
		projects: projects,
		list:     l,
	}
}

func (m ProjectSelectModel) Init() tea.Cmd {
	return nil
}

func (m ProjectSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			cmds = append(cmds, SelectProjectCmd(m.list.SelectedItem().(projectListItem).item))
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m ProjectSelectModel) View() string {
	return m.list.View()
}
