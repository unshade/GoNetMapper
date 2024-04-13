package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"main/cmd/scan_gateway"
	"main/cmd/scan_ports"
	"main/internal"
	"os"
	"strings"
	"time"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

type item struct {
	title       string
	description string
}

type tickMsg time.Time

type finishedMsg string

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type model struct {
	sub           chan string
	ready         bool
	list          list.Model
	viewport      viewport.Model
	commandOutput string
	progress      progress.Model

	needConfirm bool
	confirm     textinput.Model
}

func newModel() model {
	// Make initial list of items
	items := make([]list.Item, 4)
	items[0] = item{"Scan ports", "Scan ports of an ip address"}
	items[1] = item{"Scan gateways", "Scan your network gateways"}
	items[2] = item{"Clear", "Clear the command output"}
	items[3] = item{"Exit", "Exit the program"}

	// Setup list
	delegate := list.NewDefaultDelegate()
	commandList := list.New(items, delegate, 0, 0)
	commandList.Title = "Radar"
	commandList.Styles.Title = titleStyle

	progress := progress.New(progress.WithDefaultGradient())

	confirm := textinput.New()
	confirm.Placeholder = "Please enter an ip address..."
	confirm.CharLimit = 15

	return model{
		sub:         make(chan string),
		list:        commandList,
		progress:    progress,
		confirm:     confirm,
		needConfirm: false,
	}
}

func waitForFinishedMsg(sub chan string) tea.Cmd {
	return func() tea.Msg {
		return finishedMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		waitForFinishedMsg(m.sub),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

		if !m.ready {
			m.viewport = viewport.New(96, 19)
			m.viewport.YPosition = 0
			m.viewport.SetContent(m.commandOutput)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case tickMsg:
		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.SetPercent(internal.Progression)
		return m, tea.Batch(tickCmd(), cmd)

	case finishedMsg:
		// Set the command output to the result of the command
		m.commandOutput += fmt.Sprintf("%s", msg)
		m.viewport.SetContent(m.commandOutput)
		cmds = append(cmds, waitForFinishedMsg(m.sub))

	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		if msg.String() == "enter" {
			if m.confirm.Focused() {
				m.needConfirm = false
				m.confirm.Blur()
				ip := m.confirm.Value()
				go func() {
					read, write, err := os.Pipe()
					if err != nil {
						fmt.Println(err)
						return
					}

					scan_ports.ScanPortsCommand.SetOut(write)
					scan_ports.ScanPortsCommand.SetErr(write)

					stdout := os.Stdout
					stderr := os.Stderr

					os.Stdout = write
					os.Stderr = write

					go func() {
						// Always read from the pipe to prevent deadlock
						buf := make([]byte, 1024)
						for {
							n, err := read.Read(buf)
							if err != nil {
								break
							}
							m.sub <- string(buf[:n])
						}

					}()

					scan_ports.ScanPortsCommand.Run(scan_ports.ScanPortsCommand, []string{ip})

					write.Close()
					read.Close()

					os.Stdout = stdout
					os.Stderr = stderr
				}()
			} else {
				switch m.list.Cursor() {
				case 0:
					// Scan ports
					m.needConfirm = true
					m.confirm.Focus()

				case 1:
					// Scan gateways

					go func() {
						read, write, err := os.Pipe()
						if err != nil {
							fmt.Println(err)
							return
						}

						scan_gateway.ScanGatewayCommand.SetOut(write)
						scan_gateway.ScanGatewayCommand.SetErr(write)

						stdout := os.Stdout
						stderr := os.Stderr

						os.Stdout = write
						os.Stderr = write

						go func() {
							// Always read from the pipe to prevent deadlock
							buf := make([]byte, 1024)
							for {
								n, err := read.Read(buf)
								if err != nil {
									break
								}
								m.sub <- string(buf[:n])
							}

						}()

						scan_gateway.ScanGatewayCommand.Run(scan_gateway.ScanGatewayCommand, []string{})

						write.Close()
						read.Close()

						os.Stdout = stdout
						os.Stderr = stderr
					}()
				case 2:
					// Clear
					m.commandOutput = ""
					m.viewport.SetContent(m.commandOutput)
				case 3:
					// Exit
					return m, tea.Quit
				}
			}
		}
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.confirm, cmd = m.confirm.Update(msg)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	doc := strings.Builder{}

	commandLogStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}).
		Margin(1, 3, 0, 0).
		Padding(1, 2).
		Height(19).
		Width(96)

	if m.needConfirm {
		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			appStyle.Render(m.list.View()),
			lipgloss.JoinVertical(
				lipgloss.Top,
				m.confirm.View(),
				appStyle.Render(m.progress.View()),
				commandLogStyle.Render(m.viewport.View()),
			),
		))
	} else {
		doc.WriteString(lipgloss.JoinHorizontal(
			lipgloss.Top,
			appStyle.Render(m.list.View()),
			lipgloss.JoinVertical(
				lipgloss.Top,
				appStyle.Render(m.progress.View()),
				commandLogStyle.Render(m.viewport.View()),
			),
		))
	}

	return doc.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
