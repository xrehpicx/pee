package filepicker

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	headerMessage   string
	selectedMessage string
	filepicker      filepicker.Model
	selectedFile    string
	quitting        bool
	err             error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m model) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			m.quitting = true
			return m, tea.Quit
		case "q", "ctrl+c", "esc":
			m.quitting = true
			m.selectedFile = ""
			m.filepicker.FileSelected = ""
			return m, tea.Quit
		}
	case clearErrorMsg:
		m.err = nil
	}

	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.selectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.selectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.selectedFile == "" {
		s.WriteString(m.headerMessage)
	} else {
		s.WriteString(m.selectedMessage + m.filepicker.Styles.Selected.Render(m.selectedFile) + " <space> to select")
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n" + "<esc> or q to quit")
	return s.String()
}

func FilePicker(headerMessage string, selectedMessage string) (string, error) {
	fp := filepicker.New()
	// fp.AllowedTypes = []string{".mod", ".sum", ".go", ".txt", ".md"}
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.CurrentDirectory, _ = os.Getwd()
	fp.FileSelected, _ = os.Getwd()

	// Set default values for header and footer messages
	if headerMessage == "" {
		headerMessage = "Select file..."
	}
	if selectedMessage == "" {
		selectedMessage = "Selected file: "
	}
	m := model{
		headerMessage:   headerMessage,
		selectedMessage: selectedMessage,
		filepicker:      fp,
		selectedFile:    fp.FileSelected,
	}
	tm, err := tea.NewProgram(&m, tea.WithOutput(os.Stderr)).Run()
	mm := tm.(model)

	return mm.selectedFile, err
}
