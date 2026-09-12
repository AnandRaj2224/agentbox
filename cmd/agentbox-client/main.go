package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/AnandRaj2224/agentbox/internal/api"
	"github.com/AnandRaj2224/agentbox/internal/client"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"google.golang.org/grpc"
)

type streamMsg string
type streamDoneMsg struct{}
type streamStartMsg struct {
	stream grpc.ServerStreamingClient[api.ExecuteResponse]
}
type errMsg struct{ err error }

type model struct {
	output     string
	grpcClient api.ExecutionServiceClient
	stream     grpc.ServerStreamingClient[api.ExecuteResponse]
	codebox    textarea.Model
	runtime    string
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+r":
			m.output = ""
			return m, startExecution(m.grpcClient, m.codebox.Value(), m.runtime)
		case "ctrl+l":
			if m.runtime == "python" {
				m.runtime = "go"
			} else {
				m.runtime = "python"
			}
			return m, nil
		default:
			var cmd tea.Cmd
			m.codebox, cmd = m.codebox.Update(msg)
			return m, cmd
		}

	case streamStartMsg:
		m.stream = msg.stream
		return m, readNextChunk(m.stream)
	case streamMsg:
		m.output += string(msg)
		return m, readNextChunk(m.stream)
	case streamDoneMsg:
		m.output += "\n[Execution Complete]"
		return m, nil
	case errMsg:
		m.output += fmt.Sprintf("\n[Error: %v]", msg.err)
		return m, nil
	default:
		var cmd tea.Cmd
		m.codebox, cmd = m.codebox.Update(msg)
		return m, cmd

	}
}

var (
	codePanelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1).
			Width(50)

	outputPanelStyle = lipgloss.NewStyle().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("39")).
				Padding(0, 1).
				Width(50)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			Bold(true)
)

func (m model) View() string {
	styledCode := codePanelStyle.Render("Code:\n\n" + m.codebox.View())

	// 2. Render the right panel (Execution Output)
	styledOutput := outputPanelStyle.Render("Stream:\n\n" + m.output)

	// 3. Join them side-by-side
	topPanels := lipgloss.JoinHorizontal(lipgloss.Top, styledCode, styledOutput)

	// 4. Create a dynamic status bar at the bottom
	statusText := fmt.Sprintf(" Runtime: %s | [Ctrl+R] Run | [Ctrl+L] Swap Runtime | [Ctrl+C] Quit ", m.runtime)
	styledStatus := statusStyle.Render(statusText)

	// 5. Stack the panels on top of the status bar
	return lipgloss.JoinVertical(
		lipgloss.Left,
		topPanels,
		styledStatus,
	)
}

func startExecution(c api.ExecutionServiceClient, code string, runtime string) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		req := &api.ExecuteRequest{
			Runtime:    runtime,
			SourceCode: code,
		}

		stream, err := c.Execute(ctx, req)
		if err != nil {
			return errMsg{err: err}
		}

		return streamStartMsg{stream: stream}
	}
}

func readNextChunk(stream grpc.ServerStreamingClient[api.ExecuteResponse]) tea.Cmd {
	return func() tea.Msg {
		resp, err := stream.Recv()
		if err == io.EOF {
			return streamDoneMsg{}
		}
		if err != nil {
			return errMsg{err: err}
		}

		return streamMsg(resp.Output)
	}
}

func main() {
	c, err := client.Connect("localhost:8000")
	if err != nil {
		fmt.Printf("Failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	ta := textarea.New()
	ta.Placeholder = "Write code here..."
	ta.Focus()

	m := model{
		grpcClient: c,
		codebox:    ta,
		runtime:    "python",
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
