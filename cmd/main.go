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
		}
	}
	return m, nil
}

func cellToString(c game.Cell, isCursor bool) string {
	str := " "
	if c == game.Circle {
		str = "O"
	} else if c == game.Cross {
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
	
	return row0 + "-----------\n" + row1 + "-----------\n" + row2
}

func main() {
	p := tea.NewProgram(model{})
	p.Run()
}