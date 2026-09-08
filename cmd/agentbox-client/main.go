package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/AnandRaj2224/agentbox/internal/api"
	"github.com/AnandRaj2224/agentbox/internal/client"
	tea "github.com/charmbracelet/bubbletea"
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
}

func (m model) Init() tea.Cmd {
	return startExecution(m.grpcClient)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
	}

	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("AgentBox Terminal\n\n%s\n\nPress 'q' to quit\n", m.output)
}

func startExecution(c api.ExecutionServiceClient) tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		req := &api.ExecuteRequest{
			Runtime:    "python",
			SourceCode: "print(\"hello from TUI\")",
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

	m := model{
		grpcClient: c,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
