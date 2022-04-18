package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	size := 25
	board := newBoard(size)
	count := 0
	var cell BoardCell
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			cell = board.getCellAt(x, y)
			if cell.owner != EMPTY {
				count += 1
			}
		}
	}
	if count != 0 {
		t.Errorf("Created new board that had %d non-empty cells", count)
	}
}

func CheckCellAdjancies(t *testing.T, b *Board, x int, y int, target map[Adjacency]int) {
	cell := b.getCellAt(x, y)
	for _, adj := range Adjancies {
		count := len(b.getAdjacentCells(x, y, cell.owner, adj))
		if count != target[adj] {
			t.Errorf("Cell at (%d, %d) should have %d adjancies in %s direction but instead had %d", x, y, target[adj], adj.String(), count)
		}
	}
}

func TestAdjancies(t *testing.T) {
	size := 5
	board := newBoard(size)
	cell := board.getCellAt(2, 2)
	if cell.owner != EMPTY {
		t.Error("Cell at (2, 2) wasn't empty!")
	}
	CheckCellAdjancies(t, board, 2, 2, map[Adjacency]int{
		HORIZONTAL:             4,
		VERTICAL:               4,
		LEFT_TO_RIGHT_DIAGONAL: 4,
		RIGHT_TO_LEFT_DIAGONAL: 4,
	})
	// Make following board:
	//
	// X|O|O| |X|
	// ----------
	// X|X|O|X|X|
	// ----------
	// X|X|X|O|X|
	// ----------
	//  | |X|X|O|
	// ----------
	// O|O|O|O|O|
	//
	cells := [][]PlayerSymbol{{
		X, O, O, EMPTY, X,
	}, {
		X, X, O, X, X,
	}, {
		X, X, X, O, X,
	}, {
		EMPTY, EMPTY, X, X, O,
	}, {
		O, O, O, O, O,
	}}

	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if cells[y][x] != EMPTY {
				board.updateCell(x, y, cells[y][x])
			}
		}
	}

	CheckCellAdjancies(t, board, 0, 0, map[Adjacency]int{
		HORIZONTAL:             0,
		VERTICAL:               2,
		LEFT_TO_RIGHT_DIAGONAL: 0,
		RIGHT_TO_LEFT_DIAGONAL: 3,
	})
	CheckCellAdjancies(t, board, 4, 0, map[Adjacency]int{
		HORIZONTAL:             0,
		VERTICAL:               2,
		LEFT_TO_RIGHT_DIAGONAL: 2,
		RIGHT_TO_LEFT_DIAGONAL: 0,
	})
	CheckCellAdjancies(t, board, 2, 1, map[Adjacency]int{
		HORIZONTAL:             0,
		VERTICAL:               1,
		LEFT_TO_RIGHT_DIAGONAL: 0,
		RIGHT_TO_LEFT_DIAGONAL: 3,
	})
	CheckCellAdjancies(t, board, 2, 2, map[Adjacency]int{
		HORIZONTAL:             2,
		VERTICAL:               1,
		LEFT_TO_RIGHT_DIAGONAL: 2,
		RIGHT_TO_LEFT_DIAGONAL: 3,
	})
	CheckCellAdjancies(t, board, 0, 4, map[Adjacency]int{
		HORIZONTAL:             4,
		VERTICAL:               0,
		LEFT_TO_RIGHT_DIAGONAL: 0,
		RIGHT_TO_LEFT_DIAGONAL: 0,
	})
	CheckCellAdjancies(t, board, 4, 4, map[Adjacency]int{
		HORIZONTAL:             4,
		VERTICAL:               1,
		LEFT_TO_RIGHT_DIAGONAL: 0,
		RIGHT_TO_LEFT_DIAGONAL: 0,
	})
}
