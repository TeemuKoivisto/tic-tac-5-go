package game

import (
	"testing"
)

func TestAddingPlayers(t *testing.T) {
	size := 25
	game := New(GameOptions{
		Size:         size,
		PlayerSymbol: X,
		GameType:     HOT_SEAT,
	})
	var startErr, addErr error
	startErr = game.StartGame()
	if startErr == nil {
		t.Error("Was able to start unready game without players")
	}
	_, addErr = game.AddPlayer(HUMAN, User{ID: "1", name: "Player 1"})
	if addErr != nil {
		t.Error("Failed to add X player")
	}
	startErr = game.StartGame()
	if startErr == nil {
		t.Error("Was able to start unready with only X player")
	}
	_, addErr = game.AddPlayer(HUMAN, User{ID: "2", name: "Player 2"})
	if addErr != nil {
		t.Error("Failed to add O player")
	}
	startErr = game.StartGame()
	if startErr != nil {
		t.Error("Was not able to start game with full players")
	} else if !game.isRunning() {
		t.Error("Started game but it was not set running")
	}
	game.EndGame()
	if game.isRunning() {
		t.Error("Ended game but it was still shown running")
	}
}

func TestNewGame(t *testing.T) {
	size := 25
	game := New(GameOptions{
		Size:         size,
		PlayerSymbol: X,
		GameType:     HOT_SEAT,
	})
	game.AddPlayer(HUMAN, User{ID: "1", name: "Player 1"})
	game.AddPlayer(HUMAN, User{ID: "2", name: "Player 2"})
	game.StartGame()

	moves := []Move{
		{X: 0, Y: 0, Player: X},
		{X: 10, Y: 0, Player: O},
		{X: 11, Y: 0, Player: O},
		{X: 0, Y: 0, Player: X},
		{X: 11, Y: 0, Player: O},
		{X: 2, Y: 2, Player: X},
		{X: 12, Y: 0, Player: O},
		{X: 1, Y: 1, Player: X},
		{X: 13, Y: 0, Player: O},
		{X: 4, Y: 4, Player: X},
		{X: 0, Y: 0, Player: O},
		{X: 14, Y: 0, Player: O},
		{X: 3, Y: 3, Player: X},
		{X: 11, Y: 0, Player: O},
	}

	for i, move := range moves {
		switch i {
		case 0, 2, 3, 4, 5, 7, 9, 12:
			if game.State.Status != X_TURN {
				t.Errorf("Move %v at index %d wasn't on X's turn", move, i)
			}
		case 1, 6, 8, 10:
			if game.State.Status != O_TURN {
				t.Errorf("Move %v at index %d wasn't on O's turn", move, i)
			}
		}
		game.HandlePlayerTurn(move)
	}

	if game.State.Status != X_WON {
		t.Error("Game did not end to X victory", game)
	}
}
