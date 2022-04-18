package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	size := 25
	game := New(GameOptions{
		Size:         size,
		PlayerSymbol: X,
		GameType:     HOT_SEAT,
	})
	game.AddPlayer(HUMAN, "Player 1")
	game.AddPlayer(HUMAN, "Player 2")
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
				t.Errorf("For move %v at index %d wasn't X's turn", move, i)
			}
		case 1, 6, 8, 10:
			if game.State.Status != O_TURN {
				t.Errorf("For move %v at index %d wasn't O's turn", move, i)
			}
		}
		game.HandlePlayerTurn(move)
	}

	if game.State.Status != X_WON {
		t.Error("Game did not end to X victory", game)
	}
}
