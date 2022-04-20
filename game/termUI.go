package game

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Play() {
	fmt.Println("### TicTac5 ###")
	game := New(GameOptions{
		Size: 5,
	})
	game.AddPlayer(HUMAN, User{
		ID:   "me",
		name: "Player",
	})
	game.AddPlayer(HUMAN, User{
		ID:   "opponent",
		name: "Opponent",
	})
	game.StartGame()
	for game.isRunning() {
		PrintBoard(game)
		var x, y int
		var err error
		if game.State.Status == X_TURN {
			x, y, err = PromptMove()
		} else {
			x, y, err = PromptMove()
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
		err = game.HandlePlayerTurn(move)
		if err != nil {
			fmt.Println("error from handlePlayerTurn", err)
			continue
		}
	}
	PrintBoard(game)
	fmt.Println(game.State.Status.String())
	fmt.Println("### GAME ENDED ###")
}

func ClearScreen() {
	c := exec.Command("cmd", "/c", "cls")
	c.Stdout = os.Stdout
	c.Run()
}

func PrintBoard(t *TicTacToe) {
	ClearScreen()
	for y := 0; y < t.Opts.Size; y++ {
		for x := 0; x < t.Opts.Size; x++ {
			cell := t.State.Board.getCellAt(x, y)
			if cell.owner == EMPTY {
				fmt.Printf(" |")
			} else {
				fmt.Printf(cell.owner.String() + "|")
			}
		}
		fmt.Println()
		fmt.Println(strings.Repeat("--", t.Opts.Size))
	}
}

func PromptMove() (int, int, error) {
	var readX, readY int
	fmt.Println("Enter x,y coordinates separated by space (eg 0 1): ")
	_, err := fmt.Scanf("%d %d", &readX, &readY)
	return readX, readY, err
}
