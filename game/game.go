package game

import (
	"errors"
	"fmt"
)

type PlayerType int

const MaxPlayers = 2

const (
	HUMAN PlayerType = iota
	AI
)

type User struct {
	ID   string
	name string
}

type Player struct {
	User            User
	Type            PlayerType
	Symbol          PlayerSymbol
	AcceptedRematch bool
}

type GameOptions struct {
	Size         int
	PlayerSymbol PlayerSymbol
	GameType     GameType
}

type TicTacToe struct {
	ID      string
	Opts    GameOptions
	State   GameState
	XPlayer *Player
	OPlayer *Player
}

func New(opts GameOptions) *TicTacToe {
	return &TicTacToe{
		ID:    "unique-id",
		Opts:  opts,
		State: GameState{},
	}
}

func (t *TicTacToe) isFull() bool {
	return t.XPlayer != nil && t.OPlayer != nil
}

func (t *TicTacToe) isRunning() bool {
	return t.State.Status == X_TURN || t.State.Status == O_TURN
}

func (t *TicTacToe) AddPlayer(playerType PlayerType, user User) (*Player, error) {
	if t.isFull() {
		return nil, errors.New("game already full")
	}
	var player *Player
	if t.XPlayer == nil {
		player = &Player{
			User:            user,
			Type:            playerType,
			Symbol:          X,
			AcceptedRematch: false,
		}
		t.XPlayer = player
	} else {
		player = &Player{
			User:            user,
			Type:            playerType,
			Symbol:          O,
			AcceptedRematch: false,
		}
		t.OPlayer = player
	}
	return player, nil
}

func (t *TicTacToe) HandlePlayerTurn(move Move) error {
	if t.State.Status != X_TURN && t.State.Status != O_TURN {
		return errors.New("game has already ended")
	} else if t.State.Status == X_TURN && move.Player != X {
		return fmt.Errorf("%s tried to move on X's turn", move.Player.String())
	} else if t.State.Status == O_TURN && move.Player != O {
		return fmt.Errorf("%s tried to move on O's turn", move.Player.String())
	} else if !t.State.Board.isWithinBoard(move.X, move.Y) {
		return errors.New("x, y wasn't inside the board")
	}
	current := t.State.Board.getCellAt(move.X, move.Y)
	if current.owner != EMPTY {
		return errors.New("cell already selected")
	}
	t.State.Board.updateCell(move.X, move.Y, move.Player)
	return t.State.updateGameStatus(move)
}

func (t *TicTacToe) StartGame() error {
	if !t.isFull() {
		return errors.New("game is not full")
	}
	t.State = *newState(t.Opts.Size)
	t.State.Status = X_TURN
	return nil
}

func (t *TicTacToe) EndGame() {
	t.State.Status = TIE
}
