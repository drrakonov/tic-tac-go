package game

import (
	"errors"
	"time"
)

type Player struct {
	Name string
	Symbol Cell
	IP string
}

type Cell int 
const (
	None Cell = iota
	Cross
	Circle
)

type Game struct {
	Board [3][3]Cell
	PlayerA Player
	PlayerB Player
	Status StatusValue
	CurrentTurn Cell
	Messages []Message
}

type Message struct {
	Sender Player
	Msg string
	SentAt time.Time
}


type StatusValue int
const (
	Waiting StatusValue = iota
	Playing
	Finished
)


func (g *Game) MakeMove(row, col int) error {
	if g.Board[row][col] != None {
		return errors.New("cell is already taken!")
	}

	//Make move & Change the current turn
	g.Board[row][col] = g.CurrentTurn
	
	if g.CurrentTurn == Cross {
		g.CurrentTurn = Circle
	} else {
		g.CurrentTurn = Cross
	}

	return nil

}

func (g *Game) CheckWinner() Cell {
	if g.Board[0][0] == g.Board[1][1] && g.Board[1][1] == g.Board[2][2] && g.Board[0][0] != None {
		return g.Board[0][0]
	}

	if g.Board[2][0] == g.Board[1][1] && g.Board[1][1] == g.Board[0][2] && g.Board[2][0] != None {
		return g.Board[2][0]
	}

	if g.Board[0][0] == g.Board[0][1] && g.Board[0][1] == g.Board[0][2] && g.Board[0][0] != None {
		return g.Board[0][0]
	}

	if g.Board[1][0] == g.Board[1][1] && g.Board[1][1] == g.Board[1][2] && g.Board[1][0] != None {
		return g.Board[1][0]
	}

	if g.Board[2][0] == g.Board[2][1] && g.Board[2][1] == g.Board[2][2] && g.Board[2][0] != None {
		return g.Board[2][0]
	}

	if g.Board[0][0] == g.Board[1][0] && g.Board[1][0] == g.Board[2][0] && g.Board[0][0] != None {
		return g.Board[0][0]
	}

	if g.Board[0][1] == g.Board[1][1] && g.Board[1][1] == g.Board[2][1] && g.Board[0][1] != None {
		return g.Board[0][1]
	}

	if g.Board[0][2] == g.Board[1][2] && g.Board[1][2] == g.Board[2][2] && g.Board[0][2] != None {
		return g.Board[0][2]
	}

	return None

}

func (g *Game) IsDraw() bool {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if g.Board[i][j] == None {
				return false
			}
		}
	}

	winner := g.CheckWinner()
	if winner == None {
		return true
	}

	return false
}
