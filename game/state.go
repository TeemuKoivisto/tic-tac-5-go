package game

import (
	"errors"
)

type GameType int

const (
	HOT_SEAT GameType = iota
	LOCAL_AI
	MULTIPLAYER
)

func (t GameType) String() string {
	return []string{"HOT_SEAT", "LOCAL_AI", "MULTIPLAYER"}[t]
}

type GameStatus int

const (
	NOT_STARTED GameStatus = iota
	X_TURN
	O_TURN
	X_WON
	O_WON
	TIE
)

func (s GameStatus) String() string {
	return []string{"NOT_STARTED", "X_TURN", "O_TURN", "X_WON", "O_WON", "TIE"}[s]
}

type Move struct {
	X      int
	Y      int
	Player PlayerSymbol
}

type GameState struct {
	Board  Board
	Status GameStatus
}

func newState(size int) *GameState {
	return &GameState{
		Board:  *newBoard(size),
		Status: NOT_STARTED,
	}
}

func (g *GameState) updateGameStatus(move Move) error {
	status, playerWon := g.Status, g.CheckWin(move)
	if playerWon && status == X_TURN {
		status = X_WON
	} else if playerWon && status == O_TURN {
		status = O_WON
	} else if status == X_TURN {
		status = O_TURN
	} else if status == O_TURN {
		status = X_TURN
	} else {
		return errors.New("incorrect game state for changing player!")
	}
	g.Status = status
	return nil
}

func (g *GameState) CheckWin(lastMove Move) bool {
	cell := g.Board.getCellAt(lastMove.X, lastMove.Y)
	for _, count := range cell.adjacency {
		if count == 5 {
			return true
		}
	}
	return false
}
