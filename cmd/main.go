package main

import (
	"fmt"
	"net"
	"tic-tac-go/game"
	"tic-tac-go/network"
	tea "github.com/charmbracelet/bubbletea"
)


type model struct {
	game game.Game
	cursor [2]int
	appState network.AppState
	conn net.Conn
	isHost bool
	joinAddress string
}


func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connectionMsg:
		if msg.err != nil {
			// If connection failed, go back to menu
			m.appState = network.Menu
			return m, nil
		}
		// Connection successful! Save it and start playing!
		m.conn = msg.conn
		m.appState = network.Playing
		if !m.isHost {
			return m, waitForOpponent(m.conn) // Joiner waits for Host!
		}
		return m, nil
	case gameMsg:
		if msg.err != nil {
			m.appState = network.Menu
			return m, nil
		}
		m.game = msg.game
		return m, waitForOpponent(m.conn)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		switch m.appState {

		case network.Menu:

			switch msg.String() {

			case "h", "H":
				m.appState = network.Hosting
				m.isHost = true
				return m, startHosting()

			case "j", "J":
				m.appState = network.Joining
				m.isHost = false

			case "q", "Q":
				return m, tea.Quit

			}

		case network.Joining:

			switch msg.String() {

			case "enter":
				return m, startJoining(m.joinAddress)
			case "backspace":
				if len(m.joinAddress) > 0 {
					m.joinAddress = m.joinAddress[:len(m.joinAddress) - 1]
				}
			default:
				if len(msg.String()) == 1 {
					m.joinAddress += msg.String()
				}
			}

		case network.Hosting:
			// will do it later
		
		case network.Playing:
			
			switch msg.String() {
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
					// Are we allowed to move?
					myTurn := (m.isHost && m.game.CurrentTurn == game.Cross) || (!m.isHost && m.game.CurrentTurn == game.Circle)
					
					if myTurn {
						m.game.MakeMove(m.cursor[0], m.cursor[1]) // Make the move locally
						network.SendGameState(m.conn, m.game)     // Send it to the opponent!
						return m, waitForOpponent(m.conn)         // Now we wait for their move
					}
				}
			}
		}


		
	}
	return m, nil
}

func cellToString(c game.Cell, isCursor bool) string {
	str := " "
	switch c {
	case game.Circle:
		str = "O"
	case game.Cross:
		str = "X"
	}

	if isCursor {
		return "[" + str + "]"
	}
	return " " + str + " "
}



func (m model) boardView() string {
	row0 := fmt.Sprintf("%s|%s|%s\n", 
		cellToString(m.game.Board[0][0], m.cursor[0] == 0 && m.cursor[1] == 0), 
		cellToString(m.game.Board[0][1], m.cursor[0] == 0 && m.cursor[1] == 1), 
		cellToString(m.game.Board[0][2], m.cursor[0] == 0 && m.cursor[1] == 2),
	)
	row1 := fmt.Sprintf("%s|%s|%s\n", 
		cellToString(m.game.Board[1][0], m.cursor[0] == 1 && m.cursor[1] == 0), 
		cellToString(m.game.Board[1][1], m.cursor[0] == 1 && m.cursor[1] == 1), 
		cellToString(m.game.Board[1][2], m.cursor[0] == 1 && m.cursor[1] == 2),
	)
	row2 := fmt.Sprintf("%s|%s|%s\n", 
		cellToString(m.game.Board[2][0], m.cursor[0] == 2 && m.cursor[1] == 0), 
		cellToString(m.game.Board[2][1], m.cursor[0] == 2 && m.cursor[1] == 1), 
		cellToString(m.game.Board[2][2], m.cursor[0] == 2 && m.cursor[1] == 2),
	)
	
	boardStr := row0 + "-----------\n" + row1 + "-----------\n" + row2

	if m.game.Status == game.Finished {
		winner := m.game.CheckWinner()
		if winner == game.None {
			boardStr += "\n\nGame Over! It's a Draw! 🤝"
		} else {
			boardStr += fmt.Sprintf("\n\nGame Over! %s Wins! 🎉", cellToString(winner, false))
		}
		boardStr += "\nPress 'q' to quit."
	} else {
		// Tell them whose turn it is
        boardStr += fmt.Sprintf("\n\nCurrent Turn: %s", cellToString(m.game.CurrentTurn, false))
	}

	return boardStr

}




func (m model) View() string {
	switch m.appState {
	case network.Menu:
		return "Welcome to Tic-Tac-Go!\n\n"
	case network.Hosting:
		return "..."
	case network.Joining:
		return fmt.Sprintf("Enter Host IP: %s", m.joinAddress)
	case network.Playing:
		return m.boardView()
	}

	return  ""
}

func main() {
	initialModel := model {
		game: game.Game {
			CurrentTurn:  game.Cross, // X Goes first
		},
	}

	p := tea.NewProgram(initialModel)
	p.Run()
}
type connectionMsg struct {
	conn net.Conn
	err  error
}

func startHosting() tea.Cmd {
	return func() tea.Msg {
		conn, err := network.HostGame("8080")
		return connectionMsg{conn: conn, err: err}
	}
}

func startJoining(ip string) tea.Cmd {
	return func() tea.Msg {
		conn, err := network.JoinGame(ip + ":8080")
		return connectionMsg{conn: conn, err: err}
	}
}

type gameMsg struct {                                                                                                                                   
	game game.Game                                                                                                                                         
	err  error                                                                                                                                             
}                                                                                                                                                       
                                                                                                                                                            
func waitForOpponent(conn net.Conn) tea.Cmd {                                                                                                           
	return func() tea.Msg {                                                                                                                                
		g, err := network.ReceiveGameState(conn)                                                                                                               
		return gameMsg{game: g, err: err}                                                                                                                      
	}                                                                                                                                                      
}
