package game

import (
	"errors"
	"fmt"
	"strings"
)

type Adjacency uint8

const (
	HORIZONTAL Adjacency = iota
	VERTICAL
	LEFT_TO_RIGHT_DIAGONAL
	RIGHT_TO_LEFT_DIAGONAL
)

func (a Adjacency) String() string {
	return []string{"HORIZONTAL", "VERTICAL", "LEFT_TO_RIGHT_DIAGONAL", "RIGHT_TO_LEFT_DIAGONAL"}[a]
}

var Adjancies = [...]Adjacency{HORIZONTAL, VERTICAL, LEFT_TO_RIGHT_DIAGONAL, RIGHT_TO_LEFT_DIAGONAL}

type PlayerSymbol uint8

const (
	EMPTY PlayerSymbol = iota
	X
	O
)

func (p PlayerSymbol) String() string {
	return []string{"-", "X", "O"}[p]
}

type BoardCell struct {
	x         int
	y         int
	owner     PlayerSymbol
	adjacency map[Adjacency]int
}
type Board struct {
	size  int
	cells []BoardCell
}

func newBoard(size int) *Board {
	return &Board{
		size:  size,
		cells: *createCells(size),
	}
}

func (b *Board) clone() Board {
	cells := make([]BoardCell, b.size*b.size)
	copy(cells, b.cells)
	return Board{
		size:  b.size,
		cells: cells,
	}
}

func (b *Board) asStateString() string {
	arr := make([]string, b.size*b.size)
	for i := 0; i < b.size*b.size; i++ {
		arr[i] = b.cells[i].owner.String()
	}
	return strings.Join(arr, "")
}

func createCells(size int) *[]BoardCell {
	cells := make([]BoardCell, size*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cells[y*size+x] = BoardCell{
				x:         x,
				y:         y,
				owner:     EMPTY,
				adjacency: map[Adjacency]int{},
			}
		}
	}
	return &cells
}

func (b *Board) isWithinBoard(x int, y int) bool {
	return x >= 0 && y >= 0 && x < b.size && y < b.size
}

func (b *Board) getCellAt(x int, y int) BoardCell {
	return b.cells[y*b.size+x]
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
		} else {
			running = false
		}
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

func (b *Board) updateCell(x int, y int, player PlayerSymbol) {
	b.cells[y*b.size+x].owner = player
	for _, dir := range Adjancies {
		b.cells[y*b.size+x].adjacency[dir] = b.updateCellsInDirection(x, y, player, dir)
	}
}
