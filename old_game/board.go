package game

import (
	"errors"
	"fmt"
	"strings"
)

type BoardCell struct {
	x         int
	y         int
	owner     PlayerSymbol
	adjacency map[Adjacency]int
}
type Board struct {
	size   int
	status GameStatus
	cells  []BoardCell
}

func newBoard(size int, player PlayerSymbol) *Board {
	cells := copy(size, b.cells)
	cells[y*b.size+x].owner = player
	newBoard := Board{
		size:   b.size,
		status: b.status,
		cells:  cells,
	}
	newBoard.updateCellAdjacencies(x, y, player)
	newBoard.updateGameStatus(x, y)
	return &newBoard
}

func (b *Board) New(x int, y int, player PlayerSymbol) *Board {
	cells := copy(b.size, b.cells)
	cells[y*b.size+x].owner = player
	newBoard := Board{
		size:   b.size,
		status: b.status,
		cells:  cells,
	}
	newBoard.updateCellAdjacencies(x, y, player)
	newBoard.updateGameStatus(x, y)
	return &newBoard
}

func (b *Board) asStateString() string {
	if len(b.cells) == 0 {
		return ""
	}
	arr := make([]string, len(b.cells))
	for i, v := range b.cells {
		arr[i] = v.owner.String()
	}
	return strings.Join(arr, "")
}

func createCells(size int, grid GameState) *[]BoardCell {
	cells := make([]BoardCell, size*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell := grid.getCellAt(x, y)
			cells[y*size+x] = BoardCell{
				x:     x,
				y:     y,
				owner: cell.Owner,
				adjacency: map[Adjacency]int{
					HORIZONTAL:             0,
					VERTICAL:               0,
					LEFT_TO_RIGHT_DIAGONAL: 0,
					RIGHT_TO_LEFT_DIAGONAL: 0,
				},
			}
		}
	}
	return &cells
}

func createBoard(state *GameState) *Board {
	return &Board{
		size:   state.Opts.GridSize,
		status: state.Status,
		cells:  *createCells(state.Opts.GridSize, *state),
	}
}

func (b *Board) isWithinBoard(x int, y int) bool {
	return x >= 0 && y >= 0 && x < b.size && y < b.size
}

func (b *Board) getCellAt(x int, y int) BoardCell {
	return b.cells[y*b.size+x]
}

func (b *Board) getScore(player PlayerSymbol) int {
	switch b.status {
	case X_WON:
		if player == X {
			return 10
		} else {
			return -10
		}
	case O_WON:
		if player == O {
			return 10
		} else {
			return -10
		}
	case TIE:
		return 0
	default:
		return 0
	}
}

func (b *Board) getNextSequentialEmptyCells(x int, y int, player PlayerSymbol) *[]ScoreCell {
	var cells []ScoreCell
	nextX, nextY, running := x, y, true
	for running {
		if !b.isWithinBoard(nextX, nextY) {
			running = false
			continue
		}
		cell := b.getCellAt(nextX, nextY)
		if cell.owner == EMPTY {
			cells = append(cells, ScoreCell{
				x:     nextX,
				y:     nextY,
				score: 1,
			})
		}
		if nextX+1 == b.size {
			nextX = 0
			nextY += 1
		} else {
			nextX += 1
		}
		fmt.Printf("nextX: %d, nextY: %d\n", nextX, nextY)
	}
	fmt.Printf("cells %v\n", cells)
	return &cells
}

func copy(size int, oldCells []BoardCell) []BoardCell {
	cells := make([]BoardCell, size*size)
	for i := 0; i < size*size; i++ {
		cells[i] = BoardCell{
			x:     oldCells[i].x,
			y:     oldCells[i].y,
			owner: oldCells[i].owner,
			adjacency: map[Adjacency]int{
				HORIZONTAL:             oldCells[i].adjacency[HORIZONTAL],
				VERTICAL:               oldCells[i].adjacency[VERTICAL],
				LEFT_TO_RIGHT_DIAGONAL: oldCells[i].adjacency[LEFT_TO_RIGHT_DIAGONAL],
				RIGHT_TO_LEFT_DIAGONAL: oldCells[i].adjacency[RIGHT_TO_LEFT_DIAGONAL],
			},
		}
	}
	return cells
}

func (b *Board) getAdjacentInDirection(x int, y int, dir Adjacency, topSide bool) (BoardCell, error) {
	xx, yy := x, y
	switch dir {
	case HORIZONTAL:
		if topSide {
			xx = x + 1
		} else {
			xx = x - 1
		}
	case VERTICAL:
		if topSide {
			yy = y + 1
		} else {
			yy = y - 1
		}
	case LEFT_TO_RIGHT_DIAGONAL:
		if topSide {
			xx, yy = x-1, y+1
		} else {
			xx, yy = x+1, y-1
		}
	case RIGHT_TO_LEFT_DIAGONAL:
		if topSide {
			xx, yy = x+1, y+1
		} else {
			xx, yy = x-1, y-1
		}
	default:
		panic("inside switch-case encountered unknown Adjacency value")
	}
	if !b.isWithinBoard(xx, yy) {
		return BoardCell{}, errors.New("x,y values were not inside the board")
	}
	return b.getCellAt(xx, yy), nil
}

// Gets adjacent cells in a direction until finds a non-player cell
func (b *Board) getAdjacentCells(x int, y int, player PlayerSymbol, dir Adjacency) []BoardCell {
	var adjacent []BoardCell
	running, topSide, nowX, nowY, iters := true, true, x, y, 0
	for running {
		cell, err := b.getAdjacentInDirection(nowX, nowY, dir, topSide)
		if iters > 20 {
			fmt.Println("cell is ", cell)
			fmt.Println("player", player)
			fmt.Println("topSide", topSide)
			panic("infinite loop")
		}
		if err == nil && cell.owner == player {
			adjacent = append(adjacent, cell)
			nowX = cell.x
			nowY = cell.y
		} else if topSide {
			cell = b.getCellAt(x, y)
			topSide = false
			nowX = x
			nowY = y
		}
		running = cell.owner == player || topSide
		iters += 1
	}
	return adjacent
}

func (b *Board) updateCellsInDirection(x int, y int, player PlayerSymbol, dir Adjacency) int {
	cells := b.getAdjacentCells(x, y, player, dir)
	adjacentCount := len(cells) + 1
	for _, cell := range cells {
		cell.adjacency[dir] = adjacentCount
	}
	return adjacentCount
}

func (b *Board) updateCellAdjacencies(x int, y int, player PlayerSymbol) *Board {
	horizontal := b.updateCellsInDirection(x, y, player, HORIZONTAL)
	vertical := b.updateCellsInDirection(x, y, player, VERTICAL)
	leftToRightDiagonal := b.updateCellsInDirection(
		x,
		y,
		player,
		LEFT_TO_RIGHT_DIAGONAL,
	)
	rightToLeftDiagonal := b.updateCellsInDirection(
		x,
		y,
		player,
		RIGHT_TO_LEFT_DIAGONAL,
	)
	cell := b.getCellAt(x, y)
	cell.adjacency[HORIZONTAL] = horizontal
	cell.adjacency[VERTICAL] = vertical
	cell.adjacency[LEFT_TO_RIGHT_DIAGONAL] = leftToRightDiagonal
	cell.adjacency[RIGHT_TO_LEFT_DIAGONAL] = rightToLeftDiagonal
	return b
}

func (b *Board) updateGameStatus(x int, y int) *Board {
	status, playerWon := b.status, b.checkWin(x, y)
	if playerWon && status == X_TURN {
		status = X_WON
	} else if playerWon && status == O_TURN {
		status = O_WON
	} else if status == X_TURN {
		status = O_TURN
	} else if status == O_TURN {
		status = X_TURN
	} else {
		panic("incorrect GameStatus for changing player!")
	}
	b.status = status
	return b
}

func (b *Board) checkWin(x int, y int) bool {
	cell := b.getCellAt(x, y)
	for _, count := range cell.adjacency {
		if count == 5 {
			return true
		}
	}
	return false
}
