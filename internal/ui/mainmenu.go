package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Main menu item
type menuItem struct {
	title       string
	description string
	view        View
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

// MainMenu model
type MainMenu struct {
	BaseModel
	list list.Model
}

// Initialize a new main menu
func NewMainMenu() *MainMenu {
	m := &MainMenu{
		BaseModel: BaseModel{
			ActiveView: MainMenuView,
		},
	}

	// Define menu items
	items := []list.Item{
		menuItem{
			title:       "Projects",
			description: "List, view, and manage projects",
			view:        ProjectListView,
		},
		menuItem{
			title:       "Songs",
			description: "List, view, and manage songs",
			view:        SongListView,
		},
	}

	// Setup the list
	m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	m.list.Title = "Zupfmanager"
	m.list.SetShowStatusBar(false)
	m.list.SetFilteringEnabled(false)
	m.list.Styles.Title = TitleStyle

	// Add custom keybindings
	m.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			DefaultKeyMap.Quit,
		}
	}

	return m
}

func (m *MainMenu) Init() tea.Cmd {
	return m.InitDB()
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		h, v := lipgloss.NewStyle().Margin(1, 2).GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		// First, handle our custom keys
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit
		}

		// Handle list selection
		if msg.String() == "enter" {
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				cmds = append(cmds, SwitchView(i.view))
			}
		}

	case DBConnectedMsg:
		m.Client = msg.Client

	case ErrorMsg:
		m.Error = msg.Error
	}

	// Update the list
	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *MainMenu) View() string {
	if m.Error != "" {
		return ErrorStyle.Render(fmt.Sprintf("Error: %s", m.Error))
	}

	return lipgloss.NewStyle().Margin(1, 2).Render(m.list.View())
}
