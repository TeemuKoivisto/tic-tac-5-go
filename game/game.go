package game

import (
	"errors"
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

type TicTacToe struct {
	ID         string
	State      GameState
	Players    map[PlayerSymbol]Player
	MaxPlayers int
}

func NewGame(opts GameOptions) *TicTacToe {
	return &TicTacToe{
		State: GameState{
			Opts: opts,
		},
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
	if t.State.Opts.GameType != LOCAL_AI {
		return nil, errors.New("game wasn't a local AI game")
	}
	if xPlayer, ok := t.Players[X]; ok && xPlayer.Type == AI {
		return &xPlayer, nil
	} else if oPlayer, ok := t.Players[O]; ok && oPlayer.Type == AI {
		return &oPlayer, nil
	}
	return nil, errors.New("there was no AI players")
}

func (t *TicTacToe) HandlePlayerTurn(move Move) (*GameState, error) {
	if t.State.Status != X_TURN && t.State.Status != O_TURN {
		return nil, errors.New("game has already ended")
	} else if !t.State.isWithinGrid(move.X, move.Y) {
		return nil, errors.New("x, y wasn't inside the grid")
	}
	current := t.State.getCellAt(move.X, move.Y)
	if current.Owner == EMPTY {
		current.Owner = move.Player
		t.State.updateCell(move.X, move.Y, current)
		t.State.updateCellAdjacencies(move.X, move.Y, move.Player)
		t.State.UpdateGameStatus(move)
		return &t.State, nil
	}
	return nil, errors.New("cell already selected")
}

func (t *TicTacToe) StartGame() *TicTacToe {
	t.State.GenerateGrid()
	t.State.Status = X_TURN
	return t
}

func (t *TicTacToe) EndGame() *TicTacToe {
	t.State.Status = TIE
	return t
}
