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

// SongItem represents a song in the list
type SongItem struct {
	song *ent.Song
}

func (i SongItem) Title() string { return i.song.Title }
func (i SongItem) Description() string {
	genre := i.song.Genre
	if genre == "" {
		genre = "No genre"
	}
	return fmt.Sprintf("ID: %d | Filename: %s | Genre: %s", i.song.ID, i.song.Filename, genre)
}
func (i SongItem) FilterValue() string { return i.song.Title }

// SongList model
type SongList struct {
	BaseModel
	list  list.Model
	songs []*ent.Song
}

// Create a new song list
func NewSongList() *SongList {
	m := &SongList{
		BaseModel: BaseModel{
			ActiveView: SongListView,
			PrevView:   MainMenuView,
		},
	}

	// Setup list
	delegate := list.NewDefaultDelegate()
	m.list = list.New([]list.Item{}, delegate, 0, 0)
	m.list.Title = "Songs"
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

func (m *SongList) Init() tea.Cmd {
	return tea.Batch(
		m.InitDB(),
		m.loadSongs(),
	)
}

// Load songs from database
func (m *SongList) loadSongs() tea.Cmd {
	return func() tea.Msg {
		if m.Client == nil {
			return nil
		}

		songs, err := m.Client.Song.Query().All(context.Background())
		if err != nil {
			return ErrorMsg{Error: err.Error()}
		}

		m.songs = songs
		items := make([]list.Item, len(songs))
		for i, song := range songs {
			items[i] = SongItem{song: song}
		}

		return songsLoadedMsg{items: items}
	}
}

// Message sent when songs are loaded
type songsLoadedMsg struct {
	items []list.Item
}

func (m *SongList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			return m, SwitchView(MainMenuView)
		case key.Matches(msg, DefaultKeyMap.Add):
			// TODO: Implement song creation
			// For now, just display an error
			m.Error = "Song creation not yet implemented"
		case key.Matches(msg, DefaultKeyMap.Select) || msg.String() == "enter":
			if len(m.songs) > 0 {
				if i, ok := m.list.SelectedItem().(SongItem); ok {
					return m, tea.Batch(
						SelectSong(i.song),
						SwitchView(SongDetailView),
					)
				}
			}
		}

	case DBConnectedMsg:
		m.Client = msg.Client
		cmds = append(cmds, m.loadSongs())

	case ErrorMsg:
		m.Error = msg.Error

	case songsLoadedMsg:
		m.list.SetItems(msg.items)
	}

	// Update the list
	newList, cmd := m.list.Update(msg)
	m.list = newList
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SongList) View() string {
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
