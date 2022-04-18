package game

import (
	"errors"
	"fmt"
)

type PlayerType int

const (
	HUMAN PlayerType = iota
	AI
)

type Player struct {
	Type            PlayerType
	User            string
	Symbol          PlayerSymbol
	AcceptedRematch bool
}

type GameOptions struct {
	Size         int
	PlayerSymbol PlayerSymbol
	GameType     GameType
}

type TicTacToe struct {
	ID         string
	Opts       GameOptions
	State      GameState
	Players    map[PlayerSymbol]Player
	MaxPlayers int
}

func New(opts GameOptions) *TicTacToe {
	return &TicTacToe{
		ID:         "unique-id",
		Opts:       opts,
		State:      GameState{},
		Players:    map[PlayerSymbol]Player{},
		MaxPlayers: 2,
	}
}

func (t *TicTacToe) AddPlayer(playerType PlayerType, user string) (*Player, error) {
	if len(t.Players) == t.MaxPlayers {
		return nil, errors.New("game already full")
	}
	var player = Player{}
	if _, ok := t.Players[X]; !ok {
		player = Player{
			Symbol:          X,
			Type:            playerType,
			User:            user,
			AcceptedRematch: false,
		}
	} else if _, ok := t.Players[O]; !ok {
		player = Player{
			Symbol:          O,
			Type:            playerType,
			User:            user,
			AcceptedRematch: false,
		}
	} else {
		return nil, errors.New("game not full but has no available slots")
	}
	t.Players[player.Symbol] = player
	return &player, nil
}

func (t *TicTacToe) getAIPlayer() (*Player, error) {
	if t.Opts.GameType != LOCAL_AI {
		return nil, errors.New("game wasn't a local AI game")
	}
	if xPlayer, ok := t.Players[X]; ok && xPlayer.Type == AI {
		return &xPlayer, nil
	} else if oPlayer, ok := t.Players[O]; ok && oPlayer.Type == AI {
		return &oPlayer, nil
	}
	return nil, errors.New("there was no AI players")
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

func (t *TicTacToe) StartGame() *TicTacToe {
	t.State = *newState(t.Opts.Size)
	t.State.Status = X_TURN
	return t
}

func (t *TicTacToe) EndGame() *TicTacToe {
	t.State.Status = TIE
	return t
}
