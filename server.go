package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"unicode"
)

type Color int

const (
	NoColor Color = iota
	White
	Black
)

type Figure int

const (
	NoFigure Figure = iota
	Pawn
	Bishop
	Knight
	Rook
	Queen
	King
)

type GameStatus int

const (
	InProcess GameStatus = iota
	Check
	Checkmate
	Stalemate
	Draw
)

type Field struct {
	figure     Figure
	color     Color
	x, y          int
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Player struct {
	color Color
	kingPosition Coord
}

type Board [8][8]Field

func (b Board) At(c Coord) Field {
	return b[c.X][c.Y]
}

func (c Coord) InBounds() bool {
	return c.X >= 0 && c.X <= 7 && c.Y >= 0 && c.Y <= 7
}

func (c Color) opponent() Color {
	if c == White {
		return Black
	}
	return White
}

//forward is where the player's pawns go along y
func (c Color) forward() int {
	if c == White {
		return -1
	}
	return 1
}

//move to utility.go
func abs(x int) int {
	if x < 0 { return -x}
	return x
}

//func withinReach
func isCheck(board Board, turnColor Color, kingPosition [3]Coord) bool {
	//is there a check for a turnColor side in a given position?
	king := board[kingPosition[turnColor].X][kingPosition[turnColor].Y]
	//checking pawns
	var pawn1 Field
	var pawn2 Field
	if king.y+turnColor.forward() >= 0 && king.y+turnColor.forward() <= 7 {
		if king.x-1 >= 0 {
			pawn1 = board[king.x-1][king.y+turnColor.forward()]
		}
		if king.x+1 <= 7 {
			pawn2 = board[king.x+1][king.y+turnColor.forward()]
		}
	}
	if (pawn1.figure != 0 && pawn1.figure == Pawn && pawn1.color != turnColor) || (pawn2.figure != 0 && pawn2.figure == Pawn && pawn2.color != turnColor) {
		return true
	}
	for i := -1; i <= 1; i += 2 {
		for j := -1; j <= 1; j += 2 {
			//checking knights
			var checkingKnightField1 Field
			var checkingKnightField2 Field
			if king.x-(1*i) >= 0 && king.x-(1*i) <= 7 && king.y-(2*j) >= 0 && king.y-(2*j) <= 7 {
				checkingKnightField1 = board[king.x-(1*i)][king.y-(2*j)]
			}
			if king.x-(2*i) >= 0 && king.x-(2*i) <= 7 && king.y-(1*j) >= 0 && king.y-(1*j) <= 7 {
				checkingKnightField2 = board[king.x-(2*i)][king.y-(1*j)]
			}
			if (checkingKnightField1.figure != 0 && checkingKnightField1.figure == Knight && checkingKnightField1.color != turnColor) || (checkingKnightField2.figure != 0 && checkingKnightField2.figure == Knight && checkingKnightField2.color != turnColor) {
				return true
			}
			//checking bishops and queens(diag)
			for x, y := king.x+i, king.y+j; x <= 7 && x >= 0 && y <= 7 && y >= 0; x, y = x+i, y+j {
				if board[x][y].figure != NoFigure {
					if board[x][y].color != turnColor && (board[x][y].figure == Queen || board[x][y].figure == Bishop) {
						return true
					} else {
						break
					}
				}
			}
		}
		//checking rooks and queens(orth)
		//horisontally
		for x := king.x + i; x <= 7 && x >= 0; x += i {
			if board[x][king.y].figure != NoFigure {
				if board[x][king.y].color != turnColor && (board[x][king.y].figure == Queen || board[x][king.y].figure == Rook) {
					return true
				} else {
					break
				}
			}
		}
		//vertically
		for y := king.y + i; y <= 7 && y >= 0; y += i {
			if board[king.x][y].figure != NoFigure {
				if board[king.x][y].color != turnColor && (board[king.x][y].figure == Queen || board[king.x][y].figure == Rook) {
					return true
				} else {
					break
				}
			}
		}
	}
	//checking the enemy king
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			square := Coord{X: king.x + dx, Y: king.y + dy}
			if square.InBounds() {
				neighbour := board.At(square)
				if neighbour.figure == King && neighbour.color != turnColor {
					return true
				}
			}
		}
	}
	return false
}

//revised
func move(board *Board, fromCoord, toCoord Coord) error {
	from, to := &board[fromCoord.X][fromCoord.Y], &board[toCoord.X][toCoord.Y]
	to.figure, to.color = from.figure, from.color
	from.figure, from.color = NoFigure, NoColor
	//a pawn that reaches the last row becomes a queen
	if to.figure == Pawn && (toCoord.Y == 0 || toCoord.Y == 7) {
		to.figure = Queen
	}
	return nil
}

//revised
func isMoveLegal(toCoord Coord, legalTargets []Coord) bool {
	for _, v := range legalTargets {
		if toCoord == v {
			return true
		}
	}
	return false
}

//revised, but depends on getlegalTargets
func haslegalTargets(board Board, turnColor Color, kingPosition [3]Coord) bool {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			//check if at least any of one's pieces has legal moves
			field := board[x][y]
			if field.figure != NoFigure && field.color == turnColor {
				legalTargets := getlegalTargets(board, &board[x][y], turnColor, kingPosition);
				if len(legalTargets) >= 1 {
					return true
				}
			}
		}
	}
	return false
}

func bold(r rune) rune {
	switch {
	case 'A' <= r && r <= 'Z':
		return r - 'A' + 0x1D400
	case 'a' <= r && r <= 'z':
		return r - 'a' + 0x1D41A
	}
	return r
}

const figureRunes = " pbnrqk"

//revised
func (f Field) Rune() rune {
	r := rune(figureRunes[f.figure])
	if f.color == Black {
		r = unicode.ToUpper(r)
	}
	return r
}

//revised
func toRunes(board Board, legalTargets []Coord) [8][8]rune {
	var out [8][8]rune
	for x :=0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			out[x][y] = board[x][y].Rune()
		}
	}
	for _, t := range legalTargets {
		if board.At(t).figure == NoFigure {
			out[t.X][t.Y] = '+'
		} else {
			out[t.X][t.Y] = bold(out[t.X][t.Y])
		}
	}
	return out
}

//revised
func renderBoard(board Board, legalTargets ...Coord){
	runeBoard := toRunes(board, legalTargets)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			fmt.Printf("%c", runeBoard[x][y])
		}
		fmt.Print("\n")
	}
}

//perhaps could rework
func trivialValidation(from, to Field, turnColor Color) bool {
	if to.color == turnColor || from.figure == NoFigure || from.color != turnColor {
		return false
	}
	return true
}

func getlegalTargets(board Board, from Field, turnColor Color, kingPosition [3]Coord) []Coord {
	var legalTargets []Coord
	var to Field

	switch from.figure {
	case Pawn:
		yDirection := turnColor.forward()
		for dx := -1; dx <= 1; dx++ {
			for dy := yDirection; dy <= abs(yDirection*2); dy += yDirection {
				newX := from.x + dx
				newY := from.y + dy
				to := board[newX][newY]

				if !trivialValidation(from, to, turnColor){
					continue
				}

				//if pawn eats, its dy must be just 1
				//pawn may move forward by 2 only if it moves from the starting position
				//if pawn moves sideways, it has to eat a figure
				//
				if (dx != 0 && abs(dy) == 2) ||
					(dy == 2 && from.y != 1) ||
					(dy == -2 && from.y != 6) ||
					(dx != 0 && to.figure == NoFigure) ||
					(dx == 0 && to.figure != NoFigure) ||
					(abs(dy) == 2 && board[from.x][from.y+yDirection].figure != NoFigure) {
					continue
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: newX, Y: newY})
				}
			}
		}

	case Knight:
		for i := 1; i <= 2; i++ {
			for j := -1; j <= 1; j += 2 {
				for c := -1; c <= 1; c += 2 {
					dx := i * j
					dy := (3 - i) * c
					newX := from.x + dx
					newY := from.y + dy

					if !(Coord{X: newX, Y: newY}).InBounds() || !trivialValidation(board, from, &board[newX][newY], turnColor) {
						continue
					}

					to = board[newX][newY]

					if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
						legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
					}
				}
			}
		}

	case King:
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}

				newX := from.x + dx
				newY := from.y + dy

				if !(Coord{X: newX, Y: newY}).InBounds() || !trivialValidation(board, from, &board[newX][newY], turnColor) {
					continue
				}

				to = board[newX][newY]

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
				}
			}
		}

	case Bishop:
		for dx := -1; dx <= 1; dx += 2 {
			for dy := -1; dy <= 1; dy += 2 {
				for x, y := from.x+dx, from.y+dy; x <= 7 && x >= 0 && y <= 7 && y >= 0; x, y = x+dx, y+dy {
					to = board[x][y]

					if !trivialValidation(board, from, &to, turnColor) {
						break
					}

					if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
						legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
					}

					if to.figure != NoFigure {
						break
					}
				}
			}
		}

	case Rook:
		for d := -1; d <= 1; d += 2 {
			// Horizontal
			for dx := d; from.x+dx <= 7 && from.x+dx >= 0; dx += d {
				to = board[from.x+dx][from.y]

				if !trivialValidation(board, from, &to, turnColor) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
				}

				if to.figure != NoFigure {
					break
				}
			}

			// Vertical
			for dy := d; from.y+dy <= 7 && from.y+dy >= 0; dy += d {
				to = board[from.x][from.y+dy]

				if !trivialValidation(board, from, &to, turnColor) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
				}

				if to.figure != NoFigure {
					break
				}
			}
		}

	case Queen:
		for d := -1; d <= 1; d += 2 {
			// Horizontal
			for dx := d; from.x+dx <= 7 && from.x+dx >= 0; dx += d {
				to = board[from.x+dx][from.y]

				if !trivialValidation(board, from, &to, turnColor) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
				}

				if to.figure != NoFigure {
					break
				}
			}

			// Vertical
			for dy := d; from.y+dy <= 7 && from.y+dy >= 0; dy += d {
				to = board[from.x][from.y+dy]

				if !trivialValidation(board, from, &to, turnColor) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
				}

				if to.figure != NoFigure {
					break
				}
			}
		}

		// Diagonals
		for dx := -1; dx <= 1; dx += 2 {
			for dy := -1; dy <= 1; dy += 2 {
				for x, y := from.x+dx, from.y+dy; x <= 7 && x >= 0 && y <= 7 && y >= 0; x, y = x+dx, y+dy {
					to = board[x][y]

					if !trivialValidation(board, from, &to, turnColor) {
						break
					}

					if !isCheckAfterMove(board, from, &to, kingPosition, turnColor) {
						legalTargets = append(legalTargets, Coord{X: to.x, Y: to.y})
					}

					if to.figure != NoFigure {
						break
					}
				}
			}
		}
	}

	return legalTargets
}

func isCheckAfterMove(board Board, from, to *Field, kingPosition [3]Coord, turnColor Color) bool {
	boardMoved := board
	move(&boardMoved, Coord{X: from.x, Y: from.y}, Coord{X: to.x, Y: to.y})
	if boardMoved[to.x][to.y].figure == King {
		kingPosition[boardMoved[to.x][to.y].color] = Coord{X: to.x, Y: to.y}
	}
	return isCheck(boardMoved, turnColor, kingPosition)
}
func initialiseBoard(board *Board) error {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			switch {
			case y == 1 || y == 6:
				board[x][y].figure = Pawn
			case (y == 0 || y == 7) && (x == 0 || x == 7):
				board[x][y].figure = Rook
			case (y == 0 || y == 7) && (x == 1 || x == 6):
				board[x][y].figure = Knight
			case (y == 0 || y == 7) && (x == 2 || x == 5):
				board[x][y].figure = Bishop
			case (y == 0 || y == 7) && x == 3:
				board[x][y].figure = Queen
			case (y == 0 || y == 7) && x == 4:
				board[x][y].figure = King
			default:
				board[x][y].figure = NoFigure
			}
			board[x][y].x, board[x][y].y = x, y
			if y < 2 {
				board[x][y].color = Black
			} else if y > 5 {
				board[x][y].color = White
			}
		}
	}
	return nil
}

func declareWin(turnColor Color) {
	if turnColor == White {
		fmt.Println("White won")
	} else if turnColor == Black {
		fmt.Println("Black won")
	}
}

// codes the client gets from /gameCondition
const (
	condNormal = iota
	condCheck
	condLost
	condWon
	condStalemate
)

// Game is everything the handlers share. Every handler locks mu first.
type Game struct {
	mu         sync.Mutex
	turnChange *sync.Cond // wakes up waiting players when the turn changes or the game ends
	board      Board
	turn       Color
	kingPos    [3]Coord
	tokens     [3]string
	players    int
	from       Coord   // the piece chosen with /inputFrom
	hasFrom    bool
	legalTargets []Coord // legal moves of that piece
	inCheck    bool    // is the side to move in check
	over       bool
	winner     Color // NoColor if the game ended in a stalemate
}

func newGame() *Game {
	g := &Game{turn: White}
	g.turnChange = sync.NewCond(&g.mu)
	initialiseBoard(&g.board)
	g.kingPos[Black] = Coord{X: 4, Y: 0}
	g.kingPos[White] = Coord{X: 4, Y: 7}
	return g
}

func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// playerColor finds the player by the X-Token header. NoColor means unknown player.
func (g *Game) playerColor(c *gin.Context) Color {
	t := c.GetHeader("X-Token")
	if t == "" {
		return NoColor
	}
	if t == g.tokens[White] {
		return White
	}
	if t == g.tokens[Black] {
		return Black
	}
	return NoColor
}

func reject(c *gin.Context, message string) {
	c.JSON(400, gin.H{"status": "invalid", "error": message})
}

func (g *Game) joinGame(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.players == 2 {
		c.JSON(400, gin.H{"error": "Game full"})
		return
	}
	color := White //first player is white, second is black
	if g.players == 1 {
		color = Black
	}
	g.players++
	g.tokens[color] = newToken()
	c.JSON(200, gin.H{"colorWhite": color == White, "token": g.tokens[color]})
	fmt.Println(c.ClientIP() + " joined the game")
}

// readCoord does the checks /inputFrom and /inputTo share. If the input has to be
// rejected it answers the request itself and returns false.
func (g *Game) readCoord(c *gin.Context) (Coord, bool) {
	var pos Coord
	if err := c.ShouldBindJSON(&pos); err != nil {
		reject(c, "Invalid JSON")
		return pos, false
	}
	switch {
	case g.over:
		reject(c, "The game is over")
	case g.playerColor(c) != g.turn:
		reject(c, "It is currently not your turn")
	case !pos.InBounds():
		reject(c, "That square is outside the board")
	default:
		return pos, true
	}
	return pos, false
}

func (g *Game) inputFrom(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	pos, ok := g.readCoord(c)
	if !ok {
		return
	}
	field := g.board.At(pos)
	if field.figure == NoFigure || field.color != g.turn {
		g.hasFrom = false
		g.legalTargets = nil
		reject(c, "That is not your piece")
		return
	}
	g.from = pos
	g.hasFrom = true
	g.legalTargets = getlegalTargets(g.board, &g.board[pos.X][pos.Y], g.turn, g.kingPos)
	c.JSON(200, gin.H{"status": "ok"})
}

func (g *Game) inputTo(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	to, ok := g.readCoord(c)
	if !ok {
		return
	}
	if !g.hasFrom || !isMoveLegal(to, g.legalTargets) {
		reject(c, "Invalid move")
		return
	}
	g.play(to)
	c.JSON(200, gin.H{"status": "ok"})
}

// play makes the move chosen with /inputFrom and works out what it means for the opponent.
func (g *Game) play(to Coord) {
	move(&g.board, g.from, to)
	if g.board.At(to).figure == King {
		g.kingPos[g.turn] = to
	}
	g.hasFrom = false
	g.legalTargets = nil
	renderBoard(g.board)

	opponent := g.turn.opponent()
	g.inCheck = isCheck(g.board, opponent, g.kingPos)
	if !haslegalTargets(g.board, opponent, g.kingPos) {
		g.over = true
		if g.inCheck {
			g.winner = g.turn //checkmate, otherwise it is a stalemate and winner stays NoColor
			declareWin(g.turn)
		}
	}
	g.turn = opponent
	g.turnChange.Broadcast()
}

// queryChanged blocks until it is the caller's turn (or the game is over).
func (g *Game) queryChanged(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	me := g.playerColor(c)
	if me == NoColor {
		c.JSON(401, gin.H{"error": "unknown player"})
		return
	}
	for g.turn != me && !g.over {
		g.turnChange.Wait()
	}
	c.JSON(200, gin.H{"queryWhite": g.turn == White})
}

func (g *Game) gameCondition(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	me := g.playerColor(c)
	if me == NoColor {
		c.JSON(401, gin.H{"error": "unknown player"})
		return
	}
	condition := condNormal
	switch {
	case g.over && g.winner == me:
		condition = condWon
	case g.over && g.winner == NoColor:
		condition = condStalemate
	case g.over:
		condition = condLost
	case g.inCheck && g.turn == me:
		condition = condCheck
	}
	c.JSON(200, gin.H{"condition": condition})
}

func (g *Game) figureBoard(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	c.JSON(200, gin.H{"board": toRunes(g.board, nil)})
}

func (g *Game) legalBoard(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	c.JSON(200, gin.H{"board": toRunes(g.board, g.legalTargets)})
}

func main() {
	g := newGame()
	router := gin.Default()
	router.POST("/joinGame", g.joinGame)
	router.POST("/inputFrom", g.inputFrom)
	router.POST("/inputTo", g.inputTo)
	router.GET("/getFigureBoard", g.figureBoard)
	router.GET("/getLegalBoard", g.legalBoard)
	router.GET("/queryChanged", g.queryChanged)
	router.GET("/gameCondition", g.gameCondition)
	//moves are processed before /inputTo answers, so there is nothing to wait for any more.
	//The route stays so the client keeps working.
	router.GET("/doneProcessing", func(c *gin.Context) {
		c.JSON(200, gin.H{"isDone": true})
	})
	if err := router.Run(); err != nil {
		fmt.Println("server stopped:", err)
	}
}
