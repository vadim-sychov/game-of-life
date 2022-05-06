package main

import (
	"bytes"
	"fmt"
	"time"
)

const (
	fieldWidth           = 25
	fieldHeight          = 25
	aliveCellPlaceholder = "#"
	deadCellPlaceholder  = " "
	nextGenSleepTime     = 100
)

type Game struct {
	currentField, nextGenField GameField
}

func NewGame() *Game {
	currentField := NewGameField()
	currentField.CreateGliderSeedPattern()

	return &Game{currentField: currentField, nextGenField: NewGameField()}
}

func (game *Game) NextGeneration() {
	for i := 0; i < fieldHeight; i++ {
		for j := 0; j < fieldWidth; j++ {
			game.nextGenField[i][j] = game.currentField.IsCellAliveInNextGen(j, i)
		}
	}

	// declare current generation game field as next gen field and create a new one for next gen field
	game.currentField, game.nextGenField = game.nextGenField, NewGameField()
}

// String returns the current game field as a string
func (game *Game) String() string {
	var outputBuffer bytes.Buffer

	// show current status of each cell
	for _, fieldRow := range game.currentField {
		for _, isAliveCell := range fieldRow {
			if isAliveCell {
				outputBuffer.WriteString(aliveCellPlaceholder)
			} else {
				outputBuffer.WriteString(deadCellPlaceholder)
			}
		}
		outputBuffer.WriteString("\n")
	}

	// adds dividers below the game field
	for range game.currentField {
		outputBuffer.WriteString("-")
	}

	return outputBuffer.String()
}

type GameField [][]bool

func NewGameField() GameField {
	field := make(GameField, fieldHeight)
	for i := range field {
		field[i] = make([]bool, fieldWidth)
	}

	return field
}

// CreateGliderSeedPattern in this example just prints "glider" pattern at the middle of 25x25 game field,
// but it also can be any other seed pattern or generated using math/rand package
func (field GameField) CreateGliderSeedPattern() {
	field[12][13] = true
	field[13][14] = true
	field[14][12] = true
	field[14][13] = true
	field[14][14] = true
}

// CountAliveNeighborCells checks the adjacent cells and counts which are alive
func (field GameField) CountAliveNeighborCells(x, y int) int {
	// TODO slightly refactor algorithm
	var aliveNeighborCells int

	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			// Skip current cell
			if i == y && j == x {
				continue
			}

			// Check neighbors
			if field.IsAliveInCurrentGen(j, i) {
				aliveNeighborCells++
			}
		}
	}
	return aliveNeighborCells
}

// IsCellAliveInNextGen return next state of the cell according to the game rules
func (field GameField) IsCellAliveInNextGen(x, y int) bool {
	isAlive := field.IsAliveInCurrentGen(x, y)
	aliveNeighborCells := field.CountAliveNeighborCells(x, y)

	// Return next state of the cell according to the game rules
	if isAlive && aliveNeighborCells > 1 && aliveNeighborCells < 4 {
		// The cell is lives on to the next generation
		return true
	} else if !isAlive && aliveNeighborCells == 3 {
		// Reproduction case
		return true
	} else {
		// Underpopulation or overcrowding case
		return false
	}
}

// IsAliveInCurrentGen return cell status and handles out of bound case for coordinates,
// for example, an x value of -1 is treated as width-1
func (field GameField) IsAliveInCurrentGen(x, y int) bool {
	x += fieldWidth
	x %= fieldWidth
	y += fieldHeight
	y %= fieldHeight

	return field[y][x]
}

func main() {
	game := NewGame()

	for {
		fmt.Println(game)
		game.NextGeneration()
		time.Sleep(nextGenSleepTime * time.Millisecond)
	}
}
