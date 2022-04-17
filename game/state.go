package game

import (
	"errors"
	"fmt"
)

type GameType int

const (
	HOT_SEAT GameType = iota
	LOCAL_AI
	MULTIPLAYER
)

func (t GameType) String() string {
	return []string{"HOT_SEAT", "LOCAL_AI", "MULTIPLAYER"}[t]
}

type GameStatus int

const (
	NOT_STARTED GameStatus = iota
	X_TURN
	O_TURN
	X_WON
	O_WON
	TIE
)

func (s GameStatus) String() string {
	return []string{"NOT_STARTED", "X_TURN", "O_TURN", "X_WON", "O_WON", "TIE"}[s]
}

type Adjacency int

const (
	HORIZONTAL Adjacency = iota
	VERTICAL
	LEFT_TO_RIGHT_DIAGONAL
	RIGHT_TO_LEFT_DIAGONAL
)

type PlayerSymbol int

const (
	EMPTY PlayerSymbol = iota
	X
	O
)

func (p PlayerSymbol) String() string {
	return []string{"Empty", "X", "O"}[p]
}

type GameOptions struct {
	GridSize     int
	PlayerSymbol PlayerSymbol
	GameType     GameType
}

type GridCell struct {
	X         int
	Y         int
	Owner     PlayerSymbol
	Adjacency map[Adjacency]int
}

type Move struct {
	X      int
	Y      int
	Player PlayerSymbol
}

type GameState struct {
	Opts        GameOptions
	Status      GameStatus
	NextMove    int
	GridMap     map[string]GridCell
	MoveHistory map[int]Move
	CreatedAt   int
}

func (g *GameState) New(opts GameOptions) *GameState {
	return &GameState{
		Opts: opts,
	}
}

func (g *GameState) GenerateGrid() *GameState {
	gridSize := g.Opts.GridSize
	g.GridMap = make(map[string]GridCell)
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			g.GridMap[fmt.Sprintf("%d:%d", x, y)] = GridCell{
				X:     x,
				Y:     y,
				Owner: EMPTY,
				Adjacency: map[Adjacency]int{
					HORIZONTAL:             0,
					VERTICAL:               0,
					LEFT_TO_RIGHT_DIAGONAL: 0,
					RIGHT_TO_LEFT_DIAGONAL: 0,
				},
			}
		}
	}
	return g
}

func (g *GameState) isWithinGrid(x int, y int) bool {
	return x >= 0 && y >= 0 && x < g.Opts.GridSize && y < g.Opts.GridSize
}

func (g *GameState) getCellAt(x int, y int) GridCell {
	return g.GridMap[fmt.Sprintf("%d:%d", x, y)]
}

func (g *GameState) getAdjacentInDirection(x int, y int, dir Adjacency, topSide bool) (GridCell, error) {
	if !g.isWithinGrid(x, y) {
		return GridCell{}, errors.New("x,y values were not inside the grid")
	}
	switch dir {
	case HORIZONTAL:
		if topSide {
			return g.getCellAt(x+1, y), nil
		} else {
			return g.getCellAt(x-1, y), nil
		}
	case VERTICAL:
		if topSide {
			return g.getCellAt(x, y+1), nil
		} else {
			return g.getCellAt(x, y-1), nil
		}
	case LEFT_TO_RIGHT_DIAGONAL:
		if topSide {
			return g.getCellAt(x-1, y+1), nil
		} else {
			return g.getCellAt(x+1, y-1), nil
		}
	case RIGHT_TO_LEFT_DIAGONAL:
		if topSide {
			return g.getCellAt(x+1, y+1), nil
		} else {
			return g.getCellAt(x-1, y-1), nil
		}
	default:
		return GridCell{}, errors.New("inside switch-case encountered unknown Adjacency value")
	}
}

// Gets adjacent cells in a direction until finds a non-player cell
func (g *GameState) getAdjacentCells(x int, y int, player PlayerSymbol, dir Adjacency) []GridCell {
	var adjacent []GridCell
	running, topSide, nowX, nowY, iters := true, true, x, y, 0
	for running {
		cell, err := g.getAdjacentInDirection(nowX, nowY, dir, topSide)
		if iters > 20 {
			fmt.Println("cell is ", cell)
			fmt.Println("player", player)
			fmt.Println("topSide", topSide)
			panic("infinite loop")
		}
		if err == nil && cell.Owner == player {
			adjacent = append(adjacent, cell)
			nowX = cell.X
			nowY = cell.Y
		} else if topSide {
			cell = g.getCellAt(x, y)
			topSide = false
			nowX = x
			nowY = y
		}
		running = cell.Owner == player || topSide
		iters += 1
	}
	return adjacent
}

func (g *GameState) updateCell(x int, y int, cell GridCell) *GameState {
	g.GridMap[fmt.Sprintf("%d:%d", x, y)] = cell
	return g
}

func (g *GameState) updateCellsInDirection(x int, y int, player PlayerSymbol, dir Adjacency) int {
	cells := g.getAdjacentCells(x, y, player, dir)
	adjacentCount := len(cells) + 1
	for _, cell := range cells {
		cell.Adjacency[dir] = adjacentCount
	}
	return adjacentCount
}

func (g *GameState) updateCellAdjacencies(x int, y int, player PlayerSymbol) *GameState {
	horizontal := g.updateCellsInDirection(x, y, player, HORIZONTAL)
	vertical := g.updateCellsInDirection(x, y, player, VERTICAL)
	leftToRightDiagonal := g.updateCellsInDirection(
		x,
		y,
		player,
		LEFT_TO_RIGHT_DIAGONAL,
	)
	rightToLeftDiagonal := g.updateCellsInDirection(
		x,
		y,
		player,
		RIGHT_TO_LEFT_DIAGONAL,
	)
	cell := g.getCellAt(x, y)
	cell.Adjacency[HORIZONTAL] = horizontal
	cell.Adjacency[VERTICAL] = vertical
	cell.Adjacency[LEFT_TO_RIGHT_DIAGONAL] = leftToRightDiagonal
	cell.Adjacency[RIGHT_TO_LEFT_DIAGONAL] = rightToLeftDiagonal
	return g
}

func (g *GameState) UpdateGameStatus(move Move) (*GameState, error) {
	status, playerWon := g.Status, g.CheckWin(move)
	if playerWon && status == X_TURN {
		status = X_WON
	} else if playerWon && status == O_TURN {
		status = O_WON
	} else if status == X_TURN {
		status = O_TURN
	} else if status == O_TURN {
		status = X_TURN
	} else {
		return nil, errors.New("incorrect game state for changing player!")
	}
	g.Status = status
	return g, nil
}

func (g *GameState) CheckWin(lastMove Move) bool {
	cell := g.getCellAt(lastMove.X, lastMove.Y)
	for _, count := range cell.Adjacency {
		if count == 5 {
			return true
		}
	}
	return false
}
