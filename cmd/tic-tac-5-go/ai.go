package ticTacToe

import (
	"fmt"
	// "errors"
)

type ScoreCell struct {
	x     int
	y     int
	score int
}

type ScoreBoard struct {
	key      string
	turns    int
	score    int
	move     Move
	board    Board
	children []string
}

type AIOptions struct {
	scanDepth int
}

type LocalAI struct {
	opts   AIOptions
	symbol PlayerSymbol
	turn   int
	scores map[string]ScoreBoard
}

func (ai *LocalAI) New(opts AIOptions, symbol PlayerSymbol) *LocalAI {
	return &LocalAI{
		opts:   opts,
		symbol: symbol,
		turn:   0,
		scores: make(map[string]ScoreBoard),
	}
}

// add parent *ScoreBoard
func (ai *LocalAI) updateScores(board *Board, cell *ScoreCell, maximizingPlayer PlayerSymbol) *ScoreCell {
	boardAsString := fmt.Sprintf("%v", board.cells)
	score := board.getScore(maximizingPlayer)
	ai.scores[boardAsString] = ScoreBoard{
		key:   boardAsString,
		turns: 0,
		score: score,
		move: Move{
			X:      cell.x,
			Y:      cell.y,
			Player: maximizingPlayer,
		},
		board: *board,
	}
	// fmt.Printf("set scoreboard %v\n", boardAsString)
	return &ScoreCell{
		x:     cell.x,
		y:     cell.y,
		score: score,
	}
}

func getOppositePlayer(player PlayerSymbol) PlayerSymbol {
	if player == X {
		return O
	} else {
		return X
	}
}

func (ai *LocalAI) minimax(board *Board, cell *ScoreCell, depth int, maximizingPlayer PlayerSymbol) *ScoreCell {
	fmt.Printf("minimax depth: %d cell: %v player: %v maximizingPlayer\n", depth, cell, maximizingPlayer)
	if depth == 0 || board.status == X_WON || board.status == O_WON || board.status == TIE {
		return ai.updateScores(board, cell, maximizingPlayer)
	}
	var value int
	minimizingPlayer := getOppositePlayer(maximizingPlayer)
	if (board.status == X_TURN && maximizingPlayer == X) || (board.status == O_TURN && maximizingPlayer == O) {
		value = -10000000
		emptyCells := board.getNextSequentialEmptyCells(0, 0, maximizingPlayer)
		fmt.Printf("emptyCells: %v\n", emptyCells)
		iters := 0
		for _, cell := range *emptyCells {
			result := ai.minimax(board.New(cell.x, cell.y, maximizingPlayer), &cell, depth-1, minimizingPlayer)
			// fmt.Printf("maximized result: %v\n", result)
			fmt.Printf("iters: %v\n", iters)
			if result.score > value {
				value = result.score
			}
			iters += 1
		}
	} else {
		value = 10000000
		emptyCells := board.getNextSequentialEmptyCells(0, 0, maximizingPlayer)
		for _, cell := range *emptyCells {
			result := ai.minimax(board.New(cell.x, cell.y, maximizingPlayer), &cell, depth-1, minimizingPlayer)
			fmt.Printf("minimized result: %v\n", result)
			if result.score < value {
				value = result.score
			}
		}
	}
	return ai.updateScores(board, cell, maximizingPlayer)
}

func (ai *LocalAI) getAIMove(state *GameState) (int, int, error) {
	board := createBoard(state)
	boardAsString := board.asStateString()
	fmt.Printf("board: %v\n", board)
	fmt.Printf("boardAsString: %s\n", boardAsString)
	if oldBoard, exists := ai.scores[boardAsString]; exists {
		fmt.Printf("oldBoard: %v\n", oldBoard)
		// scan children to ai.opts.scanDepth
		return oldBoard.move.X, oldBoard.move.Y, nil
	}
	best := &ScoreCell{
		x:     -1,
		y:     -1,
		score: -10000000,
	}
	maximizingPlayer, minimizingPlayer := ai.symbol, getOppositePlayer(ai.symbol)
	emptyCells := board.getNextSequentialEmptyCells(0, 0, maximizingPlayer)
	fmt.Printf("emptyCells: %v\n", emptyCells)
	for _, cell := range *emptyCells {
		result := ai.minimax(board.New(cell.x, cell.y, maximizingPlayer), &cell, ai.opts.scanDepth-1, minimizingPlayer)
		fmt.Printf("root maximized result: %v\n", result)
		if result.score > best.score {
			best = result
		}
	}
	return best.x, best.y, nil
}

// function  minimax(node, depth, maximizingPlayer) is
// 	if depth = 0 or node is a terminal node then
// 			return the heuristic value of node
// 	if maximizingPlayer then
// 			value := −∞
// 			for each child of node do
// 					value := max(value, minimax(child, depth − 1, FALSE))
// 			return value
// 	else (* minimizing player *)
// 			value := +∞
// 			for each child of node do
// 					value := min(value, minimax(child, depth − 1, TRUE))
// 			return value

// if depth = 0 or node is a terminal node then
// 	return the heuristic value of node
// if maximizingPlayer then
// 	value := −∞
// 	for each child of node do
// 			value := max(value, alphabeta(child, depth − 1, α, β, FALSE))
// 			if value ≥ β then
// 					break (* β cutoff *)
// 			α := max(α, value)
// 	return value
// else
// 	value := +∞
// 	for each child of node do
// 			value := min(value, alphabeta(child, depth − 1, α, β, TRUE))
// 			if value ≤ α then
// 					break (* α cutoff *)
// 			β := min(β, value)
// 	return value
