package ui

import (
	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Available views
type View string

const (
	MainMenuView      View = "mainMenu"
	ProjectListView   View = "projectList"
	ProjectDetailView View = "projectDetail"
	SongListView      View = "songList"
	SongDetailView    View = "songDetail"
	AddSongView       View = "addSong"
	EditSongView      View = "editSong"
)

// Common keymap
type KeyMap struct {
	Quit    key.Binding
	Back    key.Binding
	Select  key.Binding
	Add     key.Binding
	Edit    key.Binding
	Delete  key.Binding
	Confirm key.Binding
	Cancel  key.Binding
}

// Default keymap
var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Confirm: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "confirm"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "cancel"),
	),
}

// Common model
type BaseModel struct {
	ActiveView View
	PrevView   View
	Client     *database.Client
	Error      string
	Width      int
	Height     int
}

// Initialize database connection for models
func (m *BaseModel) InitDB() tea.Cmd {
	return func() tea.Msg {
		client, err := database.New()
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}
		return DBConnectedMsg{Client: client}
	}
}

// Common messages
type ErrorMsg struct {
	Error string
}

type DBConnectedMsg struct {
	Client *database.Client
}

type ViewChangeMsg struct {
	View View
}

type ProjectSelectedMsg struct {
	Project *ent.Project
}

type SongSelectedMsg struct {
	Song *ent.Song
}

type ProjectSongSelectedMsg struct {
	ProjectSong *ent.ProjectSong
	Project     *ent.Project
	Song        *ent.Song
}

// Helper functions for models
func SwitchView(view View) tea.Cmd {
	return func() tea.Msg {
		return ViewChangeMsg{View: view}
	}
}

func SelectProject(project *ent.Project) tea.Cmd {
	return func() tea.Msg {
		return ProjectSelectedMsg{Project: project}
	}
}

func SelectSong(song *ent.Song) tea.Cmd {
	return func() tea.Msg {
		return SongSelectedMsg{Song: song}
	}
}

func SelectProjectSong(ps *ent.ProjectSong, p *ent.Project, s *ent.Song) tea.Cmd {
	return func() tea.Msg {
		return ProjectSongSelectedMsg{
			ProjectSong: ps,
			Project:     p,
			Song:        s,
		}
	}
}
