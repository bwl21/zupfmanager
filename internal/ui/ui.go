package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// UI is the main model that manages all sub-models and views
type UI struct {
	mainMenu    *MainMenu
	projectList *ProjectList
	projectView *ProjectView
	songList    *SongList
	songView    *SongView
	activeModel tea.Model
	activeView  View
}

// Initialize a new UI
func NewUI() *UI {
	mainMenu := NewMainMenu()
	projectList := NewProjectList()
	projectView := NewProjectView()
	songList := NewSongList()
	songView := NewSongView()

	return &UI{
		mainMenu:    mainMenu,
		projectList: projectList,
		projectView: projectView,
		songList:    songList,
		songView:    songView,
		activeModel: mainMenu,
		activeView:  MainMenuView,
	}
}

func (m *UI) Init() tea.Cmd {
	return m.activeModel.Init()
}

func (m *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case ViewChangeMsg:
		// Switch the active view based on the message
		m.activeView = msg.View
		switch msg.View {
		case MainMenuView:
			m.activeModel = m.mainMenu
		case ProjectListView:
			m.activeModel = m.projectList
		case ProjectDetailView:
			m.activeModel = m.projectView
		case SongListView:
			m.activeModel = m.songList
		case SongDetailView:
			m.activeModel = m.songView
		default:
			return m, nil
		}
		cmd = m.activeModel.Init()
		cmds = append(cmds, cmd)

	case DBConnectedMsg:
		// Share the DB client across all models
		m.mainMenu.Client = msg.Client
		m.projectList.Client = msg.Client
		m.projectView.Client = msg.Client
		m.songList.Client = msg.Client
		m.songView.Client = msg.Client
	}

	// Update the active model
	newModel, cmd := m.activeModel.Update(msg)
	if newModel != nil {
		// Handle different model types
		switch v := newModel.(type) {
		case *MainMenu:
			m.mainMenu = v
			m.activeModel = v
		case *ProjectList:
			m.projectList = v
			m.activeModel = v
		case *ProjectView:
			m.projectView = v
			m.activeModel = v
		case *SongList:
			m.songList = v
			m.activeModel = v
		case *SongView:
			m.songView = v
			m.activeModel = v
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *UI) View() string {
	if m.activeModel == nil {
		return "Loading..."
	}

	return fmt.Sprintf("%s", m.activeModel.View())
}

// Run the UI application
func RunUI() error {
	p := tea.NewProgram(NewUI(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
