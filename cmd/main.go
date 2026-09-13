package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
	"tic-tac-go/game"
	"tic-tac-go/network"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	game        game.Game
	cursor      [2]int
	appState    network.AppState
	conn        net.Conn
	isHost      bool
	joinAddress string
	hostPIN     string

	isChatting bool
	chatInput  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connectionMsg:
		if msg.err != nil {
			m.appState = network.Menu
			return m, nil
		}
		m.conn = msg.conn
		m.appState = network.Playing
		if !m.isHost {
			return m, listenToNetwork(m.conn)
		}
		// Host also listens continuously now!
		return m, listenToNetwork(m.conn)

	case netPayloadMsg:
		if msg.err != nil {
			m.appState = network.Menu
			return m, nil
		}
		if msg.payload.Type == "MOVE" {
			m.game.MakeMove(msg.payload.Row, msg.payload.Col)
		} else if msg.payload.Type == "CHAT" {
			sender := "Player 2"
			if m.isHost {
				sender = "Player 2"
			} else {
				sender = "Player 1"
			}
			m.game.Messages = append(m.game.Messages, game.Message{
				Sender: game.Player{Name: sender},
				Msg:    msg.payload.Text,
			})
		}
		return m, listenToNetwork(m.conn)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.appState {
		case network.Menu:
			switch msg.String() {
			case "h", "H":
				m.appState = network.Hosting
				m.isHost = true
				rand.Seed(time.Now().UnixNano())
				m.hostPIN = fmt.Sprintf("%04d", rand.Intn(10000))
				return m, startHosting(m.hostPIN)
			case "j", "J":
				m.appState = network.Joining
				m.isHost = false
			case "q", "Q":
				return m, tea.Quit
			}

		case network.Joining:
			switch msg.String() {
			case "enter":
				parts := strings.Split(m.joinAddress, ":")
				if len(parts) == 2 {
					return m, startJoining(parts[0], parts[1])
				}
			case "backspace":
				if len(m.joinAddress) > 0 {
					m.joinAddress = m.joinAddress[:len(m.joinAddress)-1]
				}
			case "esc":
				m.appState = network.Menu
			default:
				if len(msg.String()) == 1 {
					m.joinAddress += msg.String()
				}
			}

		case network.Hosting:
			if msg.String() == "esc" {
				m.appState = network.Menu
			}

		case network.Playing:
			if m.isChatting {
				switch msg.String() {
				case "esc":
					m.isChatting = false
					m.chatInput = ""
				case "enter":
					if m.chatInput != "" {
						sender := "O"
						if m.isHost {
							sender = "X"
						}
						m.game.Messages = append(m.game.Messages, game.Message{
							Sender: game.Player{Name: sender},
							Msg:    m.chatInput,
						})
						network.SendPayload(m.conn, network.Payload{
							Type: "CHAT",
							Text: m.chatInput,
						})
						m.chatInput = ""
						m.isChatting = false
					}
				case "backspace":
					if len(m.chatInput) > 0 {
						m.chatInput = m.chatInput[:len(m.chatInput)-1]
					}
				default:
					if len(msg.String()) == 1 || msg.String() == " " {
						m.chatInput += msg.String()
					}
				}
				return m, nil
			}

			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "t", "T":
				m.isChatting = true
			case "up", "k":
				if m.cursor[0] > 0 {
					m.cursor[0]--
				}
			case "left", "h":
				if m.cursor[1] > 0 {
					m.cursor[1]--
				}
			case "down", "j":
				if m.cursor[0] < 2 {
					m.cursor[0]++
				}
			case "right", "l":
				if m.cursor[1] < 2 {
					m.cursor[1]++
				}
			case "enter", " ":
				if m.game.Status != game.Finished {
					myTurn := (m.isHost && m.game.CurrentTurn == game.Cross) || (!m.isHost && m.game.CurrentTurn == game.Circle)
					if myTurn {
						err := m.game.MakeMove(m.cursor[0], m.cursor[1])
						if err == nil {
							network.SendPayload(m.conn, network.Payload{
								Type: "MOVE",
								Row:  m.cursor[0],
								Col:  m.cursor[1],
							})
						}
					}
				}
			}
		}
	}
	return m, nil
}

// ---------------------------------------------------------
// LIPGLOSS UI RENDERING
// ---------------------------------------------------------

func renderCell(c game.Cell, isCursor bool) string {
	str := ""
	switch c {
	case game.Circle:
		str = oStyle.Render("O")
	case game.Cross:
		str = xStyle.Render("X")
	}

	style := cellStyle
	if isCursor {
		style = cursorStyle
	}
	return style.Render(str)
}

func (m model) boardView() string {
	myIdentity := "Player 2 (O)"
	if m.isHost {
		myIdentity = "Player 1 (X)"
	}
	identityUI := identityStyle.Render(fmt.Sprintf("Playing as %s", myIdentity))

	var rows []string
	for i := 0; i < 3; i++ {
		var cols []string
		for j := 0; j < 3; j++ {
			cols = append(cols, renderCell(m.game.Board[i][j], m.cursor[0] == i && m.cursor[1] == j))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cols...))
	}
	boardUI := lipgloss.JoinVertical(lipgloss.Left, rows...)

	statusUI := ""
	if m.game.Status == game.Finished {
		winner := m.game.CheckWinner()
		if winner == game.None {
			statusUI = activeTurnStyle.Render("Draw.") + statusStyle.Render("\nPress 'q' to quit.")
		} else {
			winChar := "X"
			if winner == game.Circle {
				winChar = "O"
			}
			statusUI = activeTurnStyle.Render(fmt.Sprintf("%s wins.", winChar)) + statusStyle.Render("\nPress 'q' to quit.")
		}
	} else {
		turnChar := "O"
		if m.game.CurrentTurn == game.Cross {
			turnChar = "X"
		}
		statusUI = statusStyle.Render("Turn: ") + activeTurnStyle.Render(turnChar) + statusStyle.Render("  •  Press 't' to chat")
	}

	leftSide := boardContainer.Render(lipgloss.JoinVertical(lipgloss.Left, identityUI, boardUI, statusUI))

	chatHeader := chatHeaderStyle.Render("Chat")
	
	var chatLines []string
	for _, msg := range m.game.Messages {
		senderStyle := chatSenderXStyle
		if msg.Sender.Name == "Player 2" || msg.Sender.Name == "O" {
			senderStyle = chatSenderOStyle
		}
		formattedMsg := fmt.Sprintf("%s %s", senderStyle.Render(msg.Sender.Name+":"), chatTextStyle.Render(msg.Msg))
		chatLines = append(chatLines, formattedMsg)
	}
	
	fullChat := lipgloss.JoinVertical(lipgloss.Left, chatLines...)
	lines := strings.Split(fullChat, "\n")
	for len(lines) < 13 {
		lines = append(lines, "")
	}
	if len(lines) > 13 {
		lines = lines[len(lines)-13:]
	}
	chatHistory := lipgloss.JoinVertical(lipgloss.Left, lines...)

	inputUI := ""
	if m.isChatting {
		inputUI = chatInputBox.Render(fmt.Sprintf("▶ %s_", m.chatInput))
	} else {
		inputUI = statusStyle.Render("▶ (press 't' to type)")
	}

	chatUI := lipgloss.JoinVertical(lipgloss.Left, chatHeader, chatHistory, inputUI)
	rightSide := chatContainer.Render(chatUI)

	content := lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
	title := titleStyle.Render("Tic-Tac-Go")
	return appLayout.Render(lipgloss.JoinVertical(lipgloss.Left, title, content))
}

func (m model) View() string {
	title := titleStyle.Render("Tic-Tac-Go")

	switch m.appState {
	case network.Menu:
		menu := statusStyle.Render("• [H] Host a game\n• [J] Join a game\n• [Q] Quit")
		return appLayout.Render(lipgloss.JoinVertical(lipgloss.Left, title, menu))
	case network.Hosting:
		msg := statusStyle.Render(fmt.Sprintf("Room created securely!\n\nTell your friend to connect to:\nYOUR_IP:%s", m.hostPIN))
		return appLayout.Render(lipgloss.JoinVertical(lipgloss.Left, title, msg))
	case network.Joining:
		msg := statusStyle.Render(fmt.Sprintf("Enter Connection String (IP:PIN)\nExample: 192.168.1.5:1234\n\n> %s_\n\n(Press Enter to connect)", m.joinAddress))
		return appLayout.Render(lipgloss.JoinVertical(lipgloss.Left, title, msg))
	case network.Playing:
		return m.boardView()
	}
	return ""
}

// ---------------------------------------------------------
// BACKGROUND TASKS
// ---------------------------------------------------------

type connectionMsg struct {
	conn net.Conn
	err  error
}

func startHosting(pin string) tea.Cmd {
	return func() tea.Msg {
		conn, err := network.HostGame("8080", pin)
		return connectionMsg{conn: conn, err: err}
	}
}

func startJoining(ip, pin string) tea.Cmd {
	return func() tea.Msg {
		conn, err := network.JoinGame(ip+":8080", pin)
		return connectionMsg{conn: conn, err: err}
	}
}

type netPayloadMsg struct {
	payload network.Payload
	err     error
}

func listenToNetwork(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		p, err := network.ReceivePayload(conn)
		return netPayloadMsg{payload: p, err: err}
	}
}

func main() {
	initialModel := model{
		game: game.Game{
			CurrentTurn: game.Cross, // X Goes first
		},
	}
	p := tea.NewProgram(initialModel)
	p.Run()
}
