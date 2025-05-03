package ui

import (
	"context"
	"fmt"

	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectSongItem represents a song associated with a project
type ProjectSongItem struct {
	projectSong *ent.ProjectSong
	song        *ent.Song
}

func (i ProjectSongItem) Title() string { return i.song.Title }
func (i ProjectSongItem) Description() string {
	return fmt.Sprintf("ID: %d | Priority: %d | Difficulty: %s | Filename: %s",
		i.song.ID, i.projectSong.Priority, i.projectSong.Difficulty, i.song.Filename)
}
func (i ProjectSongItem) FilterValue() string { return i.song.Title }

// ProjectView model for displaying project details and associated songs
type ProjectView struct {
	BaseModel
	project      *ent.Project
	list         list.Model
	projectSongs []*ent.ProjectSong
	songs        []*ent.Song
}

// Create a new project view
func NewProjectView() *ProjectView {
	m := &ProjectView{
		BaseModel: BaseModel{
			ActiveView: ProjectDetailView,
			PrevView:   ProjectListView,
		},
	}

	delegate := list.NewDefaultDelegate()
	m.list = list.New([]list.Item{}, delegate, 0, 0)
	m.list.Styles.Title = TitleStyle
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeyMap.Add,
			DefaultKeyMap.Edit,
			DefaultKeyMap.Delete,
			DefaultKeyMap.Back,
			DefaultKeyMap.Quit,
		}
	}

	return m
}

func (m *ProjectView) Init() tea.Cmd {
	return m.InitDB()
}

// Load project songs from database
func (m *ProjectView) loadProjectSongs() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil || m.project == nil {
			return nil
		}

		// Load the project with associated songs
		projectSongs, err := m.Client.ProjectSong.Query().
			Where(projectsong.ProjectID(m.project.ID)).
			WithSong().
			All(context.Background())
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}

		m.projectSongs = projectSongs
		items := make([]list.Item, len(projectSongs))
		for i, ps := range projectSongs {
			items[i] = ProjectSongItem{
				projectSong: ps,
				song:        ps.Edges.Song,
			}
		}

		return projectSongsLoadedMsg{items: items}
	}
}

// Message sent when project songs are loaded
type projectSongsLoadedMsg struct {
	items []list.Item
}

func (m *ProjectView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// Handle keys
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Back):
			return m, SwitchView(ProjectListView)
		case key.Matches(msg, DefaultKeyMap.Add):
			// TODO: Implement adding a song to the project
			m.Error = "Adding a song to project not yet implemented"
		case key.Matches(msg, DefaultKeyMap.Delete):
			// Only allow deletion if there's a selected item
			if i, ok := m.list.SelectedItem().(ProjectSongItem); ok && len(m.projectSongs) > 0 {
				// TODO: Implement removing a song from the project
				m.Error = fmt.Sprintf("Removing song '%s' from project not yet implemented", i.song.Title)
			}
		case key.Matches(msg, DefaultKeyMap.Edit):
			// Only allow editing if there's a selected item
			if i, ok := m.list.SelectedItem().(ProjectSongItem); ok && len(m.projectSongs) > 0 {
				// TODO: Implement editing a project song
				m.Error = fmt.Sprintf("Editing song '%s' in project not yet implemented", i.song.Title)
			}
		}

	case DBConnectedMsg:
		m.Client = msg.Client
		if m.project != nil {
			cmds = append(cmds, m.loadProjectSongs())
		}

	case ProjectSelectedMsg:
		m.project = msg.Project
		m.list.Title = fmt.Sprintf("Project: %s", m.project.Title)
		cmds = append(cmds, m.loadProjectSongs())

	case ErrorMsg:
		m.Error = msg.Error

	case projectSongsLoadedMsg:
		m.list.SetItems(msg.items)
	}

	// Update the list
	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *ProjectView) View() string {
	if m.Error != "" {
		return ErrorStyle.Render(fmt.Sprintf("Error: %s", m.Error))
	}

	if m.project == nil {
		return "No project selected."
	}

	// Project info
	projectInfo := fmt.Sprintf("Project ID: %d", m.project.ID)
	projectInfoView := SubtitleStyle.Render(projectInfo)

	// Help text
	helpView := "\nj/k: navigate • a: add song • e: edit • d: delete • esc: back • q: quit"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		projectInfoView,
		lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()),
		lipgloss.NewStyle().Margin(0, 2).Foreground(subtleColor).Render(helpView),
	)
}
