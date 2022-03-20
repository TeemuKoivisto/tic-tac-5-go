package ticTacToe

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Play() {
	game := NewGame(GameOptions{
		GridSize: 5,
	})
	game.AddPlayer(HUMAN, "human 1")
	// game.AddPlayer(HUMAN, "human 2")
	game.AddPlayer(AI, "computer")
	game.StartGame()
	ai := &LocalAI{
		opts: AIOptions{
			scanDepth: 2,
		},
		symbol: O,
		turn:   0,
		scores: make(map[string]ScoreBoard),
	}
	running := true
	for running {
		PrintGrid(game)
		var x, y int
		var err error
		if game.State.Status == X_TURN {
			x, y, err = PromptMove()
		} else {
			x, y, err = ai.getAIMove(&game.State)
			fmt.Printf("AI move (x: %d y: %d)\n", x, y)
		}
		if err != nil {
			fmt.Println("error", err)
			continue
		}
		var player PlayerSymbol
		if game.State.Status == X_TURN {
			player = X
		} else {
			player = O
		}
		move := Move{
			X:      x,
			Y:      y,
			Player: player,
		}
		state, err := game.HandlePlayerTurn(move)
		if err != nil {
			fmt.Println("error from handlePlayerTurn", err)
			continue
		}
		if state.Status != X_TURN && state.Status != O_TURN {
			running = false
		}
	}
	PrintGrid(game)
	fmt.Println("GAME ENDED: ", game.State.Status.String())
}

func ClearScreen() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func PrintGrid(t *TicTacToe) {
	ClearScreen()
	for y := 0; y < t.State.Opts.GridSize; y++ {
		for x := 0; x < t.State.Opts.GridSize; x++ {
			cell := t.State.getCellAt(x, y)
			if cell.Owner == EMPTY {
				fmt.Printf(" |")
			} else {
				fmt.Printf(cell.Owner.String() + "|")
			}
		}
		fmt.Println()
		fmt.Println(strings.Repeat("--", t.State.Opts.GridSize))
	}
}

func PromptMove() (int, int, error) {
	// return 1, 1, nil
	var readX, readY int
	fmt.Println("Enter x,y coordinates separated by space (eg 0 1): ")
	_, err := fmt.Scanf("%d %d", &readX, &readY)
	return readX, readY, err
}
