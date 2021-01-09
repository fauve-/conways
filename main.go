package main

import (
	// TODO: isn't there a module system these days?

	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowHeight = 768
	windowWidth  = 1024
	// They're square!
	pixelSize = 5
	// this yields 160972 pixels
	// thatsalotta pixels
	boardHeight = windowHeight / pixelSize
	boardWidth  = windowWidth / pixelSize
	boardSpaces = boardHeight * boardWidth
)

// MARK: - Game implementaiton

// oh no I forgot we don't have enums
type CellState = uint

const (
	dead  = CellState(2)
	alive = CellState(3)
)

func prepareBoard() []CellState {
	// initialized to zero, right?
	return make([]CellState, boardSpaces)
}

func populateBoard(board []CellState) {
	for x := 0; x < boardWidth; x++ {
		for y := 0; y < boardHeight; y++ {
			idx := xyToIndex(x, y)
			board[idx] = getRandomState()
		}
	}
}

func howManyLivingNeighbors(x, y int, board []CellState) int {
	neighborIndicies := []int{
		// Top row
		xyToIndex(x-1, y+1), xyToIndex(x, y+1), xyToIndex(x+1, y+1),
		xyToIndex(x-1, y) /*don't count self obvs*/, xyToIndex(x+1, y),
		xyToIndex(x-1, y-1), xyToIndex(x, y-1), xyToIndex(x+1, y-1),
	}

	count := 0
	for _, idx := range neighborIndicies {
		if board[idx] == alive {
			count++
		}
	}
	return count
}

func turn(board []CellState) []CellState {
	nextBoard := prepareBoard()
	for x := 0; x < boardWidth; x++ {
		for y := 0; y < boardHeight; y++ {
			// THE CONWAY COMMANDMENTS
			// Any live cell with two or three live neighbours survives.
			// Any dead cell with three live neighbours becomes a live cell.
			// All other live cells die in the next generation. Similarly, all other dead cells stay dead.
			index := xyToIndex(x, y)
			var newState CellState
			livingNeighbors := howManyLivingNeighbors(x, y, board)

			if board[index] == alive {
				if livingNeighbors == 2 || livingNeighbors == 3 {
					newState = alive
				} else {
					newState = dead
				}
			} else {
				if livingNeighbors == 3 {
					newState = alive
				} else {
					newState = dead
				}
			}
			nextBoard[index] = newState
		}
	}
	return nextBoard
}

func xyToIndex(x, y int) int {
	// TODO: decide on better rules
	// We'll wrap around for now
	if x < 0 {
		x = boardWidth - 1
	}
	if x > boardWidth-1 {
		x = 0
	}
	if y > boardHeight-1 {
		y = 0
	}
	if y < 0 {
		y = boardHeight - 1
	}
	return x + (boardWidth * y)
}

// TODO: make ratio configurable
// we'll do 2-1 dead to alive right now
func getRandomState() CellState {
	switch rand.Int31n(3) {
	case 0:
		fallthrough
	case 1:
		return dead
	case 2:
		return alive
	}
	panic("ya didn't specify the argument to rand.Int31n properly!")
}

// MARK: - Drawing Implementation
func square(imd *imdraw.IMDraw, x, y int, state CellState) {
	xf := float64(x)
	yf := float64(y)
	var colour color.Color
	if state == alive {
		colour = color.White
	} else {
		colour = color.Black
	}

	imd.SetColorMask(colour)
	imd.Push(pixel.V(xf, yf))
	imd.Push(pixel.V(xf+5, yf+5))
	imd.Rectangle(0)
}

// This is probably insanely inefficient
func drawBoard(imd *imdraw.IMDraw, board []CellState) {
	for x := 0; x < boardWidth; x++ {
		for y := 0; y < boardHeight; y++ {
			idx := xyToIndex(x, y)
			square(imd, x*pixelSize, y*pixelSize, board[idx])
		}
	}
}

func anyDead(board []CellState) bool {
	for _, b := range board {
		if b == dead {
			return true
		}
	}
	return false
}

func anyAlive(board []CellState) bool {
	for _, b := range board {
		if b == alive {
			return true
		}
	}
	return false
}

func run() {
	// allocate our game
	board := prepareBoard()
	populateBoard(board)
	screenRect := pixel.R(0, 0, windowWidth, windowHeight)
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: screenRect,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	canvas := pixelgl.NewCanvas(screenRect)
	imd := imdraw.New(nil)

	if err != nil {
		panic(err)
	}

	canvas.Clear(colornames.White)
	imd.Clear()
	win.Clear(colornames.Skyblue)
	second := time.Tick(time.Second / 30)
	turnNumber := 0
	// fmt.Printf("were any alive on turn %+v? %+v %+v \n", turnNumber, anyAlive(board), anyDead(board))
	for !win.Closed() {
		win.Clear(colornames.Black)
		canvas.Clear(colornames.Black)
		imd.Clear()

		drawBoard(imd, board)
		imd.Draw(canvas)

		win.SetMatrix(pixel.IM.Scaled(pixel.ZV,
			math.Min(
				win.Bounds().W()/canvas.Bounds().W(),
				win.Bounds().H()/canvas.Bounds().H(),
			),
		).Moved(win.Bounds().Center()))

		canvas.Draw(win, pixel.IM)
		//canvas.Draw(win, pixel.IM.Moved(canvas.Bounds().Center()))

		if win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Update()
		board = turn(board)
		turnNumber++
		// fmt.Printf("were any alive on turn %+v? %+v %+v \n", turnNumber, anyAlive(board), anyDead(board))
		<-second
	}
}

func main() {
	pixelgl.Run(run)
}
