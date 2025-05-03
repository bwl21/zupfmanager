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

// ProjectItem represents a project in the list that contains the song
type SongProjectItem struct {
	projectSong *ent.ProjectSong
	project     *ent.Project
}

func (i SongProjectItem) Title() string { return i.project.Title }
func (i SongProjectItem) Description() string {
	return fmt.Sprintf("ID: %d | Priority: %d | Difficulty: %s",
		i.project.ID, i.projectSong.Priority, i.projectSong.Difficulty)
}
func (i SongProjectItem) FilterValue() string { return i.project.Title }

// SongView model for displaying song details and associated projects
type SongView struct {
	BaseModel
	song         *ent.Song
	list         list.Model
	projectSongs []*ent.ProjectSong
	projects     []*ent.Project
}

// Create a new song view
func NewSongView() *SongView {
	m := &SongView{
		BaseModel: BaseModel{
			ActiveView: SongDetailView,
			PrevView:   SongListView,
		},
	}

	delegate := list.NewDefaultDelegate()
	m.list = list.New([]list.Item{}, delegate, 0, 0)
	m.list.Styles.Title = TitleStyle
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeyMap.Back,
			DefaultKeyMap.Quit,
		}
	}

	return m
}

func (m *SongView) Init() tea.Cmd {
	return m.InitDB()
}

// Load projects that contain this song
func (m *SongView) loadSongProjects() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil || m.song == nil {
			return nil
		}

		// Load projects that contain this song
		projectSongs, err := m.Client.ProjectSong.Query().
			Where(projectsong.SongID(m.song.ID)).
			WithProject().
			All(context.Background())
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}

		m.projectSongs = projectSongs
		items := make([]list.Item, len(projectSongs))
		for i, ps := range projectSongs {
			items[i] = SongProjectItem{
				projectSong: ps,
				project:     ps.Edges.Project,
			}
		}

		return songProjectsLoadedMsg{items: items}
	}
}

// Message sent when song projects are loaded
type songProjectsLoadedMsg struct {
	items []list.Item
}

func (m *SongView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, SwitchView(SongListView)
		case key.Matches(msg, DefaultKeyMap.Select) || msg.String() == "enter":
			if len(m.projectSongs) > 0 {
				if i, ok := m.list.SelectedItem().(SongProjectItem); ok {
					return m, tea.Batch(
						SelectProject(i.project),
						SwitchView(ProjectDetailView),
					)
				}
			}
		}

	case DBConnectedMsg:
		m.Client = msg.Client
		if m.song != nil {
			cmds = append(cmds, m.loadSongProjects())
		}

	case SongSelectedMsg:
		m.song = msg.Song
		m.list.Title = fmt.Sprintf("Song: %s", m.song.Title)
		cmds = append(cmds, m.loadSongProjects())

	case ErrorMsg:
		m.Error = msg.Error

	case songProjectsLoadedMsg:
		m.list.SetItems(msg.items)
	}

	// Update the list
	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SongView) View() string {
	if m.Error != "" {
		return ErrorStyle.Render(fmt.Sprintf("Error: %s", m.Error))
	}

	if m.song == nil {
		return "No song selected."
	}

	// Song info
	genre := m.song.Genre
	if genre == "" {
		genre = "No genre"
	}
	songInfo := fmt.Sprintf("Song ID: %d\nFilename: %s\nGenre: %s", m.song.ID, m.song.Filename, genre)
	songInfoView := SubtitleStyle.Render(songInfo)

	// List title
	projectsHeader := TitleStyle.Render("Projects containing this song")

	// Help text
	helpView := "\nj/k: navigate • enter: view project • esc: back • q: quit"

	return lipgloss.JoinVertical(
		lipgloss.Left,
		songInfoView,
		lipgloss.NewStyle().Margin(1, 0).Render(projectsHeader),
		lipgloss.NewStyle().Margin(1, 2).Render(m.list.View()),
		lipgloss.NewStyle().Margin(0, 2).Foreground(subtleColor).Render(helpView),
	)
}
