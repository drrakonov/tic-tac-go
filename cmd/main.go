package main

import (
	"fmt"
	"tic-tac-go/game"

	tea "github.com/charmbracelet/bubbletea"
)


type model struct {
	game game.Game
	cursor [2]int
}


func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
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
				m.game.MakeMove(m.cursor[0], m.cursor[1])
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

func (m model) View() string {
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

func main() {
	initialModel := model {
		game: game.Game {
			CurrentTurn:  game.Cross, // X Goes first
		},
	}

	p := tea.NewProgram(initialModel)
	p.Run()
}