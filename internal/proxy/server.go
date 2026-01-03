package proxy

import (
	"fmt"
	"net"
	"time"

	"github.com/Jonathanthedeveloper/havoc.git/internal/logger"
	"github.com/Jonathanthedeveloper/havoc.git/internal/state"
	"github.com/charmbracelet/lipgloss"
)

func Start(state *state.HavocState) error {
	connection := state.GetConnection()
	initialPort := connection.Port
	maxRetries := 10

	var listener net.Listener
	var err error
	var finalPort int

	for i := 0; i < maxRetries; i++ {
		port := initialPort + i
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			finalPort = port
			break
		}
		logger.ErrorF("failed to bind port %d: %s, trying next port...", port, err)
	}

	if err != nil {
		return fmt.Errorf("failed to find available port after %d attempts", maxRetries)
	}

	// Update the state with the actual port we are listening on
	state.SetPort(finalPort)

	logger.PrintF("proxy server started on port %d -> %s\n", finalPort, connection.Target)

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.ErrorF("failed to accept connection: %s", err)
			continue
		}

		go handleConnection(state, conn)
	}
}

func handleConnection(state *state.HavocState, conn net.Conn) error {
	defer conn.Close()

	targetConn, err := net.Dial("tcp", state.GetConnection().Target)

	if err != nil {
		logger.ErrorF("failed to connect to target %s: %s", state.GetConnection().Target, err)
		return err
	}

	logger.Debug("Connected to target: %s (RemoteAddr: %s)", state.GetConnection().Target, targetConn.RemoteAddr().String())

	defer targetConn.Close()

	clientSniffer := NewSniffer(conn, true)
	targetSniffer := NewSniffer(targetConn, false)

	startTime := time.Now()

	go Copy(state, clientSniffer, targetSniffer)
	err = Copy(state, targetSniffer, clientSniffer)

	duration := time.Since(startTime)

	method := clientSniffer.Method
	if method == "" {
		method = "???"
	}

	status := targetSniffer.Status
	if status == "" {
		status = "???"
	}

	methodColor := lipgloss.Color("205") // Pink
	switch method {
	case "GET":
		methodColor = lipgloss.Color("205") // Pink
	case "POST":
		methodColor = lipgloss.Color("39") // Blue
	case "PUT":
		methodColor = lipgloss.Color("39") // Blue
	case "PATCH":
		methodColor = lipgloss.Color("39") // Blue
	case "DELETE":
		methodColor = lipgloss.Color("196") // Red
	default:
		methodColor = lipgloss.Color("205") // Pink
	}

	statusColor := lipgloss.Color("46") // Green
	switch status {
	case "1":
		statusColor = lipgloss.Color("205") // Pink
	case "2":
		statusColor = lipgloss.Color("46") // Green
	case "3":
		statusColor = lipgloss.Color("226") // Yellow
	case "4":
		statusColor = lipgloss.Color("226") // Yellow
	case "5":
		statusColor = lipgloss.Color("196") // Red
	default:
		statusColor = lipgloss.Color("46") // Green
	}

	methodStyle := lipgloss.NewStyle().Foreground(methodColor).Bold(true)
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	logger.PrintF("%s %s %s %s",
		methodStyle.Render(method),
		clientSniffer.Path,
		statusStyle.Render(status),
		dimStyle.Render(duration.Round(time.Millisecond).String()),
	)

	return err
}
