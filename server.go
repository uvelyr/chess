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
	figure Figure
	color  Color
	x, y   int
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Player struct {
	color        Color
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

// TODO: move this helper to utility.go.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// func withinReach
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

func move(board *Board, sourceCoord, targetCoord Coord) error {
	source := &board[sourceCoord.X][sourceCoord.Y]
	target := &board[targetCoord.X][targetCoord.Y]
	target.figure, target.color = source.figure, source.color
	source.figure, source.color = NoFigure, NoColor

	//pawn that reaches the last rank is promoted to a queen
	if target.figure == Pawn && (targetCoord.Y == 0 || targetCoord.Y == 7) {
		target.figure = Queen
	}
	return nil
}

func isMoveLegal(targetCoord Coord, legalTargets []Coord) bool {
	for _, legalTarget := range legalTargets {
		if targetCoord == legalTarget {
			return true
		}
	}
	return false
}

//returns true if at least one piece has at least one legal target
func hasLegalTargets(board Board, turnColor Color, kingPosition [3]Coord) bool {
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			source := &board[x][y]
			if source.figure != NoFigure && source.color == turnColor {
				if len(getLegalTargets(board, source, turnColor, kingPosition)) > 0 {
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

func (f Field) Rune() rune {
	r := rune(figureRunes[f.figure])
	if f.color == Black {
		r = unicode.ToUpper(r)
	}
	return r
}

func toRunes(board Board, legalTargets []Coord) [8][8]rune {
	var out [8][8]rune
	for x := 0; x < 8; x++ {
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

func renderBoard(board Board, legalTargets ...Coord) {
	runeBoard := toRunes(board, legalTargets)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			fmt.Printf("%c", runeBoard[x][y])
		}
		fmt.Print("\n")
	}
}

func trivialMoveValidation(source, target *Field, turnColor Color) bool {
	return target.color != turnColor && source.figure != NoFigure && source.color == turnColor
}

func getLegalTargets(board Board, source *Field, turnColor Color, kingPosition [3]Coord) []Coord {
	var legalTargets []Coord

	appendIfLegal := func(targetX, targetY int) {
		targetCoord := Coord{X: targetX, Y: targetY}
		if !targetCoord.InBounds() {
			return
		}
		target := &board[targetX][targetY]
		if !trivialMoveValidation(source, target, turnColor) {
			return
		}
		if isCheckAfterMove(board, source, target, kingPosition, turnColor) {
			return
		}
		legalTargets = append(legalTargets, targetCoord)
	}

	switch source.figure {
	case Pawn:
		yDirection := turnColor.forward()
		startRank := 1
		if turnColor == White {
			startRank = 6
		}
		for step := 1; step <= 2; step++ {
			dy := step * yDirection
			for dx := -1; dx <= 1; dx++ {
				targetX, targetY := source.x+dx, source.y+dy
				coord := Coord{X: targetX, Y: targetY}
				if !coord.InBounds() {
					continue
				}
				target := &board[targetX][targetY]
				if !trivialMoveValidation(source, target, turnColor) {
					continue
				}
				//pawns move straight, capture diagonally, and may advance two
				//squares only from their starting rank with a clear path.
				if dx == 0 && target.figure != NoFigure {
					continue
				}
				if dx != 0 && target.figure == NoFigure {
					continue
				}
				if step == 2 && (dx != 0 || source.y != startRank || board[source.x][source.y+yDirection].figure != NoFigure) {
					continue
				}
				if !isCheckAfterMove(board, source, target, kingPosition, turnColor) {
					legalTargets = append(legalTargets, coord)
				}
			}
		}

	case Knight:
		for longStep := 1; longStep <= 2; longStep++ {
			for _, xSign := range []int{-1, 1} {
				for _, ySign := range []int{-1, 1} {
					dx := longStep * xSign
					dy := (3 - longStep) * ySign
					appendIfLegal(source.x+dx, source.y+dy)
				}
			}
		}

	case King:
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}
				appendIfLegal(source.x+dx, source.y+dy)
			}
		}

	case Bishop, Rook, Queen:
		directions := make([]Coord, 0, 8)
		if source.figure == Bishop || source.figure == Queen {
			for dx := -1; dx <= 1; dx += 2 {
				for dy := -1; dy <= 1; dy += 2 {
					directions = append(directions, Coord{X: dx, Y: dy})
				}
			}
		}
		if source.figure == Rook || source.figure == Queen {
			directions = append(directions,
				Coord{X: -1, Y: 0}, Coord{X: 1, Y: 0},
				Coord{X: 0, Y: -1}, Coord{X: 0, Y: 1},
			)
		}
		for _, direction := range directions {
			for x, y := source.x+direction.X, source.y+direction.Y; (Coord{X: x, Y: y}).InBounds(); x, y = x+direction.X, y+direction.Y {
				target := &board[x][y]
				if !trivialMoveValidation(source, target, turnColor) {
					break
				}
				if !isCheckAfterMove(board, source, target, kingPosition, turnColor) {
					legalTargets = append(legalTargets, Coord{X: x, Y: y})
				}
				if target.figure != NoFigure {
					break
				}
			}
		}
	}

	return legalTargets
}

func isCheckAfterMove(board Board, source, target *Field, kingPosition [3]Coord, turnColor Color) bool {
	boardAfterMove := board
	move(&boardAfterMove,
		Coord{X: source.x, Y: source.y},
		Coord{X: target.x, Y: target.y},
	)
	movedTarget := boardAfterMove[target.x][target.y]
	if movedTarget.figure == King {
		kingPosition[movedTarget.color] = Coord{X: target.x, Y: target.y}
	}
	return isCheck(boardAfterMove, turnColor, kingPosition)
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

// codes the client gets source /gameCondition
const (
	condNormal = iota
	condCheck
	condLost
	condWon
	condStalemate
)

// Game is everything the handlers share. Every handler locks mu first.
type Game struct {
	mu           sync.Mutex
	turnChange   *sync.Cond //wakes up waiting players when the turn changes or the game ends
	board        Board
	turn         Color
	kingPos      [3]Coord
	tokens       [3]string
	players      int
	source       Coord //the piece chosen with /inputSource
	hasSource    bool
	legalTargets []Coord //legal moves of that piece
	inCheck      bool    //is the side to move in check
	over         bool
	winner       Color //NoColor if a stalemate
}

func newGame() *Game {
	g := &Game{turn: White}
	g.turnChange = sync.NewCond(&g.mu)
	initialiseBoard(&g.board)
	g.kingPos[Black] = Coord{X: 4, Y: 0}
	g.kingPos[White] = Coord{X: 4, Y: 7}
	return g
}

//generate a random token
func newToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

//playerColor finds the player by the X-Token header
func (g *Game) playerColor(c *gin.Context) Color {
	t := c.GetHeader("X-Token")
	if t == "" {
        //unidentified
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

	if g.players >= 2 {
		c.JSON(400, gin.H{"error": "Game is full"})
		return
	}
	var color Color
	if g.players == 0 {
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			c.JSON(500, gin.H{"error": "Failed to assign player color"})
			return
		}
		if b[0]%2 == 1 {
			color = Black
		} else {    
            color = White
        }
	} else {
		//the second player is of the opposite color.
		if g.tokens[White] != "" {
			color = Black
		} else {
			color = White
		}
	}
	g.players++
	g.tokens[color] = newToken()

	c.JSON(200, gin.H{
		"colorWhite": color == White,
		"token":      g.tokens[color],
	})
	fmt.Println(c.ClientIP() + " joined the game")
}

//readCoord does the checks /inputSource and /inputTarget share. If the input has to be
//rejected it answers the request itself and returns false.
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
		reject(c, "That square is outside of the board")
	default:
		return pos, true
	}
	return pos, false
}

func (g *Game) inputSource(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	pos, ok := g.readCoord(c)
	if !ok {
		return
	}
	field := g.board.At(pos)
	if field.figure == NoFigure || field.color != g.turn {
		g.hasSource = false
		g.legalTargets = nil
		reject(c, "That is not your piece")
		return
	}
	g.source = pos
	g.hasSource = true
	g.legalTargets = getLegalTargets(g.board, &g.board[pos.X][pos.Y], g.turn, g.kingPos)
	c.JSON(200, gin.H{"status": "ok"})
}

func (g *Game) inputTarget(c *gin.Context) {
	g.mu.Lock()
	defer g.mu.Unlock()
	target, ok := g.readCoord(c)
	if !ok {
		return
	}
	if !g.hasSource || !isMoveLegal(target, g.legalTargets) {
		reject(c, "Invalid move")
		return
	}
	g.play(target)
	c.JSON(200, gin.H{"status": "ok"})
}

// play makes the move chosen with /inputSource and works out what it means for the opponent
func (g *Game) play(target Coord) {
	move(&g.board, g.source, target)
    //if a player moves the king, update its coordinates
	if g.board.At(target).figure == King {
		g.kingPos[g.turn] = target
	}
	g.hasSource = false
	g.legalTargets = nil
	renderBoard(g.board)

	opponent := g.turn.opponent()
	g.inCheck = isCheck(g.board, opponent, g.kingPos)
    //if there are no legal moves, game is over
	if !hasLegalTargets(g.board, opponent, g.kingPos) {
		g.over = true
		if g.inCheck {
			g.winner = g.turn //checkmate, otherwise it is a stalemate and winner stays NoColor
			declareWin(g.turn)
		}
	}
	g.turn = opponent
	g.turnChange.Broadcast()
}

//queryChanged blocks until it is the caller's turn (or the game is over)
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
	router.POST("/inputSource", g.inputSource)
	router.POST("/inputTarget", g.inputTarget)
	router.GET("/getFigureBoard", g.figureBoard)
	router.GET("/getLegalBoard", g.legalBoard)
	router.GET("/queryChanged", g.queryChanged)
	router.GET("/gameCondition", g.gameCondition)
	router.GET("/doneProcessing", func(c *gin.Context) {
		c.JSON(200, gin.H{"isDone": true})
	})
	if err := router.Run(); err != nil {
		fmt.Println("server stopped:", err)
	}
}
