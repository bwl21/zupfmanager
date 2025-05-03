package ui

import (
	"context"
	"fmt"

	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectItem represents a project in the list
type ProjectItem struct {
	project *ent.Project
}

func (i ProjectItem) Title() string       { return i.project.Title }
func (i ProjectItem) Description() string { return fmt.Sprintf("ID: %d", i.project.ID) }
func (i ProjectItem) FilterValue() string { return i.project.Title }

// ProjectList model
type ProjectList struct {
	BaseModel
	list     list.Model
	projects []*ent.Project
}

// Create a new project list
func NewProjectList() *ProjectList {
	m := &ProjectList{
		BaseModel: BaseModel{
			ActiveView: ProjectListView,
			PrevView:   MainMenuView,
		},
	}

	// Setup list
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("yellow"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("white"))
	delegate.ShowDescription = true
	m.list = list.New([]list.Item{}, delegate, 0, 0)
	m.list.Title = "Projects"
	m.list.Styles.Title = TitleStyle
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeyMap.Add,
			DefaultKeyMap.Back,
			DefaultKeyMap.Quit,
		}
	}

	return m
}

func (m *ProjectList) Init() tea.Cmd {
	return tea.Batch(
		m.InitDB(),
		m.loadProjects(),
	)
}

// Load projects from database
func (m *ProjectList) loadProjects() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil {
			return nil
		}

		projects, err := m.Client.Project.Query().All(context.Background())
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}

		m.projects = projects
		items := make([]list.Item, len(projects))
		for i, project := range projects {
			items[i] = ProjectItem{project: project}
		}

		return projectsLoadedMsg{items: items}
	}
}

// Message sent when projects are loaded
type projectsLoadedMsg struct {
	items []list.Item
}

func (m *ProjectList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h-2, msg.Height-v-2)

	case tea.KeyMsg:
		// Handle keys
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Back):
			return m, SwitchView(MainMenuView)
		case key.Matches(msg, DefaultKeyMap.Add):
			// TODO: Implement project creation
			// For now, just display an error
			m.Error = "Project creation not yet implemented"
		case key.Matches(msg, DefaultKeyMap.Select) || msg.String() == "enter":
			if len(m.projects) > 0 {
				if i, ok := m.list.SelectedItem().(ProjectItem); ok {
					return m, tea.Batch(
						SelectProject(i.project),
						SwitchView(ProjectDetailView),
					)
				}
			}
		}

	case DBConnectedMsg:
		m.Client = msg.Client
		cmds = append(cmds, m.loadProjects())

	case ErrorMsg:
		m.Error = msg.Error

	case projectsLoadedMsg:
		m.list.SetItems(msg.items)
	}

	// Update the list
	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProjectList) View() string {
	if m.Error != "" {
		return ErrorStyle.Render(fmt.Sprintf("Error: %s", m.Error))
	}

	helpView := "\nj/k: navigate • enter: select • a: add • esc: back • q: quit"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()),
		lipgloss.NewStyle().Margin(0, 2).Foreground(subtleColor).Render(helpView),
	)
}
