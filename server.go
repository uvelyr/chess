package main

import (
	"fmt"
//	"os"
	"github.com/gin-gonic/gin"
	"time"
)
// p pawn, b bishop, q queen, K king, n knight, r rook
type Field struct {
	figure rune
	colorWhite *bool
	x int
	y int
}
type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}
func abs(x int) int {
	if x<0 {
		x=-x
	}
	return x
}
func handle(err error) {
	if(err!=nil) {
		fmt.Println("Error:", err)
	}
}
func Bool2int(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = 0
	}
	return i
}
func Int2bool(b int) bool {
	if b == 1 {
		return true
	} else {
		return false
	}
}
func Bool2int2(b bool) int {
	var i int
	if b {
		i = 1
	} else {
		i = -1
	}
	return i
}
func isCheck(board [8][8]Field, queryWhite bool, kingPosition[2]Coord) bool {
	//is there a check for a queryWhite side in a given position?
	king := board[kingPosition[Bool2int(queryWhite)].X][kingPosition[Bool2int(queryWhite)].Y]
	//checking pawns
	var pawn1 Field
	var pawn2 Field
	if(king.y+Bool2int2(!queryWhite)>= 0 && king.y+Bool2int2(!queryWhite)<=7){
		if(king.x-1>=0) {
			pawn1 = board[king.x-1][king.y+Bool2int2(!queryWhite)]
		}
		if(king.x+1<=7) {
			pawn2 = board[king.x+1][king.y+Bool2int2(!queryWhite)]
		}
	}
	if (pawn1.figure != 0 && pawn1.figure=='p' && *pawn1.colorWhite!=queryWhite) || (pawn2.figure != 0 && pawn2.figure=='p' && *pawn2.colorWhite!=queryWhite) {
		return true
	}
	for i :=-1; i<=1; i+=2 {
		for j :=-1; j<=1; j+=2 {
			//checking knights
			var checkingKnightField1 Field
			var checkingKnightField2 Field
			if(king.x-(1*i)>=0 && king.x-(1*i) <=7 && king.y-(2*j)>=0 && king.y-(2*j)<=7){
				checkingKnightField1 = board[king.x-(1*i)][king.y-(2*j)]
			}
			if(king.x-(2*i)>=0 && king.x-(2*i) <=7 && king.y-(1*j)>=0 && king.y-(1*j)<=7){
				checkingKnightField2 = board[king.x-(2*i)][king.y-(1*j)]
			}
			if (checkingKnightField1.figure != 0 && checkingKnightField1.figure == 'n' && *checkingKnightField1.colorWhite!=queryWhite) || (checkingKnightField2.figure != 0 && checkingKnightField2.figure == 'n' && *checkingKnightField2.colorWhite!=queryWhite) {
				return true
			}
			//checking bishops and queens(diag)
			for x, y := king.x+i, king.y+j; x<=7 && x>=0 && y<=7 && y>=0 ; x, y = x+i, y+j {
				if board[x][y].figure != 32 {
					if *board[x][y].colorWhite!=queryWhite && (board[x][y].figure=='q' || board[x][y].figure=='b') {
						return true
					} else {
						break
					}
				}
			}
		}
		//checking rooks and queens(orth)
		//horisontally
		for x := king.x+i; x<=7 && x>=0; x += i {
			if board[x][king.y].figure != 32 {
				if *board[x][king.y].colorWhite!=queryWhite && (board[x][king.y].figure=='q' || board[x][king.y].figure=='r'){
					return true
				} else {
					break
				}
			}
		}
		//vertically
		for y := king.y+i; y<=7 && y>=0; y += i {
			if board[king.x][y].figure != 32 {
				if *board[king.x][y].colorWhite!=queryWhite && (board[king.x][y].figure=='q' || board[king.x][y].figure=='r'){
					return true
				} else {
					break
				}
			}
		}
	}
	return false
}
func move(board *[8][8]Field, from, to *Field) {
	to.figure=from.figure
	to.colorWhite=from.colorWhite
	from.figure=' '
	from.colorWhite=nil
}
func isMoveCorrect(board [8][8]Field, from, to *Field, kingPosition [2]Coord, legalMoves []Coord, queryWhite bool) bool {
	for _, v := range legalMoves {
		if v.X==to.x && v.Y==to.y {
			return true
		}
	}
	return false
}
func isMate(board [8][8]Field, queryWhite bool, kingPosition [2]Coord) (bool, error) {
	for x := 0; x<8; x++ {
		for y := 0; y<8; y++ {
			if board[x][y].figure!=32 && *board[x][y].colorWhite==queryWhite{
				if legalMoves := getLegalMoves(board, &board[x][y], queryWhite, kingPosition); len(legalMoves) >=1 {
					return false, nil
				}
			}
		}
	}
	return true, nil
}
func renderBoard(board [8][8]Field) {
	for y:=0; y<8; y++ {
		for x:= 0; x<8; x++ {
			if board[x][y].colorWhite != nil {
				if *board[x][y].colorWhite == true {
					fmt.Printf("%c", board[x][y].figure)
				} else {
					fmt.Printf("%c", board[x][y].figure-32)
				}
			} else {
				fmt.Printf("%c", board[x][y].figure)
			}
		}
		fmt.Print("\n")
	}
}
func renderBoardWithLegalMoves(board [8][8]Field, legalMoves []Coord) {
	boardWithLegalMoves := board
	for _, v := range legalMoves {
		field := &boardWithLegalMoves[v.X][v.Y]
		if field.figure==32 {
			field.figure = '+'
		} else if *field.colorWhite{
			field.figure = rune(field.figure) + 119737
		} else {
			field.figure = rune(field.figure) + 119743
		}
	}
	renderBoard(boardWithLegalMoves)
}
func trivialValidation(board [8][8]Field, from, to *Field, queryWhite bool) (bool){
	if (to.colorWhite!=nil && *to.colorWhite==queryWhite) || from.figure==32 || *from.colorWhite!=queryWhite {
		return false
	}
	return true
}
func doesNotExceedBoundaries(x, y int) bool {
	return x>=0 && x<=7 && y>=0 && y<=7
}
func getLegalMoves(board [8][8]Field, from *Field, queryWhite bool, kingPosition [2]Coord) []Coord {
	var legalMoves []Coord
	var to Field

	switch from.figure {
	case 'p':
		yDirection := Bool2int2(!queryWhite)
		for dx := -1; dx <= 1; dx++ {
			for dy := yDirection; dy != yDirection*3; dy += yDirection {
				newX := from.x + dx
				newY := from.y + dy

				if !doesNotExceedBoundaries(newX, newY) || !trivialValidation(board, from, &board[newX][newY], queryWhite) {
					continue
				}

				to := board[newX][newY]

				if (dx != 0 && abs(dy) == 2) ||
					(dy == 2 && from.y != 1) ||
					(dy == -2 && from.y != 6) ||
					(dx != 0 && to.figure == 32) ||
					(dx == 0 && to.figure != 32) {
					continue
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: newX, Y: newY})
				}
			}
		}

	case 'n':
		for i := 1; i <= 2; i++ {
			for j := -1; j <= 1; j += 2 {
				for c := -1; c <= 1; c += 2 {
					dx := i * j
					dy := (3 - i) * c
					newX := from.x + dx
					newY := from.y + dy

					if !doesNotExceedBoundaries(newX, newY) || !trivialValidation(board, from, &board[newX][newY], queryWhite) {
						continue
					}

					to = board[newX][newY]

					if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
						legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
					}
				}
			}
		}

	case 'k':
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx == 0 && dy == 0 {
					continue
				}

				newX := from.x + dx
				newY := from.y + dy

				if !doesNotExceedBoundaries(newX, newY) || !trivialValidation(board, from, &board[newX][newY], queryWhite) {
					continue
				}

				to = board[newX][newY]

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
				}
			}
		}

	case 'b':
		for dx := -1; dx <= 1; dx += 2 {
			for dy := -1; dy <= 1; dy += 2 {
				for x, y := from.x+dx, from.y+dy; x <= 7 && x >= 0 && y <= 7 && y >= 0; x, y = x+dx, y+dy {
					to = board[x][y]

					if !trivialValidation(board, from, &to, queryWhite) {
						break
					}

					if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
						legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
					}

					if to.figure != 32 {
						break
					}
				}
			}
		}

	case 'r':
		for d := -1; d <= 1; d += 2 {
			// Horizontal
			for dx := d; from.x+dx <= 7 && from.x+dx >= 0; dx += d {
				to = board[from.x+dx][from.y]

				if !trivialValidation(board, from, &to, queryWhite) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
				}

				if to.figure != 32 {
					break
				}
			}

			// Vertical
			for dy := d; from.y+dy <= 7 && from.y+dy >= 0; dy += d {
				to = board[from.x][from.y+dy]

				if !trivialValidation(board, from, &to, queryWhite) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
				}

				if to.figure != 32 {
					break
				}
			}
		}

	case 'q':
		for d := -1; d <= 1; d += 2 {
			// Horizontal
			for dx := d; from.x+dx <= 7 && from.x+dx >= 0; dx += d {
				to = board[from.x+dx][from.y]

				if !trivialValidation(board, from, &to, queryWhite) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
				}

				if to.figure != 32 {
					break
				}
			}

			// Vertical
			for dy := d; from.y+dy <= 7 && from.y+dy >= 0; dy += d {
				to = board[from.x][from.y+dy]

				if !trivialValidation(board, from, &to, queryWhite) {
					break
				}

				if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
					legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
				}

				if to.figure != 32 {
					break
				}
			}
		}

		// Diagonals
		for dx := -1; dx <= 1; dx += 2 {
			for dy := -1; dy <= 1; dy += 2 {
				for x, y := from.x+dx, from.y+dy; x <= 7 && x >= 0 && y <= 7 && y >= 0; x, y = x+dx, y+dy {
					to = board[x][y]

					if !trivialValidation(board, from, &to, queryWhite) {
						break
					}

					if !isCheckAfterMove(board, from, &to, kingPosition, queryWhite) {
						legalMoves = append(legalMoves, Coord{X: to.x, Y: to.y})
					}

					if to.figure != 32 {
						break
					}
				}
			}
		}
	}

	return legalMoves
}

func isCheckAfterMove(board [8][8]Field, from, to *Field, kingPosition[2]Coord, queryWhite bool) (bool){
	boardMoved := board
	f := &boardMoved[from.x][from.y]
	t := &boardMoved[to.x][to.y]
	move(&boardMoved, f, t)
	if t.figure=='k' {
		a := Bool2int(*t.colorWhite)
		kingPosition[a].X = to.x
		kingPosition[a].Y = to.y
	}
	isCheck := isCheck(boardMoved, queryWhite, kingPosition)
	return isCheck
}
func simpleRender(board [8][8]rune) {
	for y:=0; y<8; y++ {
		for x:=0; x<8; x++ {
			fmt.Printf("%c", board[x][y])
		}
		fmt.Print("\n")
	}
}
func processToFigureBoard(board [8][8]Field) [8][8]rune {
	var figureBoard [8][8]rune
	for x:=0; x<8; x++ {
		for y:=0; y<8; y++ {
			if board[x][y].colorWhite != nil {
				if *board[x][y].colorWhite == true {
					figureBoard[x][y] = board[x][y].figure
				} else {
					figureBoard[x][y] = board[x][y].figure-32
				}
			} else {
				figureBoard[x][y] = 32
			}
		}
	}
	return figureBoard
}
func processToLegalBoard(board [8][8]Field, legalMoves []Coord) [8][8]rune {
	var boardWithLegalMoves = processToFigureBoard(board)
	for _, v := range legalMoves {
		field := board[v.X][v.Y]
		if field.figure==32 {
			boardWithLegalMoves[v.X][v.Y] = '+'
		} else if *field.colorWhite{
			boardWithLegalMoves[v.X][v.Y] = rune(field.figure) + 119737
		} else {
			boardWithLegalMoves[v.X][v.Y] = rune(field.figure) + 119743
		}
	}
	return boardWithLegalMoves
}
func initialiseBoard(board *[8][8]Field) error {
	for y:=0; y<8; y++ {
		for x:=0; x<8; x++ {
			switch {
				case y==1 || y==6:
					board[x][y].figure = 'p'
				case (y==0 || y==7) && (x==0 || x==7):
					board[x][y].figure = 'r'
				case (y==0 || y==7) && (x==1 || x==6):
					board[x][y].figure = 'n'
				case (y==0 || y==7) && (x==2 || x==5):
					board[x][y].figure = 'b'
				case (y==0 || y==7) && x==3:
					board[x][y].figure = 'q'
				case (y==0 || y==7) && x==4:
					board[x][y].figure = 'k'
				default: board[x][y].figure = ' '
			}
			board[x][y].x=x
			board[x][y].y=y
			if y<2 {
				board[x][y].colorWhite = BoolPtr(false)
			} else if y>5 {
				board[x][y].colorWhite = BoolPtr(true)
			}
		}
	}
	return nil
}
func BoolPtr(b bool) *bool {
	return &b
}
func declareWin(queryWhite bool) {
	if queryWhite {
		fmt.Println("White won")
	} else {
		fmt.Println("Black won")
	}
}

var router *gin.Engine
var inputChan = make(chan Coord)
func inputCoord(coord *Coord) gin.HandlerFunc{
        return func(c *gin.Context) {
                var input Coord
                if err := c.BindJSON(&input); err!=nil {
                        c.JSON(
                                400,
                                gin.H{
                                        "error": "Invalid JSON",
                                },
                        )
                        fmt.Println("mistake 1")
			return
                }
		if c.ClientIP()!=playerIP[Bool2int(queryWhite)] {
			c.JSON(
				400,
				gin.H{
					"error": "Not your query, dumbass",
				},
			)
                        fmt.Println("mistake 2")
			return
		}
		inputChan <- input
                c.JSON(
                        200,
                        gin.H{
                                "receivedX": input.X,
                                "receivedY": input.Y,
                        },
                )
                fmt.Println("mistake 3")
                fmt.Println(input.X, input.Y)
		*coord = input
        }
}
func inputField(field **Field, input *Coord, board *[8][8]Field) {
	*input = <- inputChan
	*field = &(*board)[input.X][input.Y]
}
var count = 0
var playerIP [2]string
func joinGame(c *gin.Context) {
	var colorWhite bool
	if count == 0 {
		colorWhite = true//if 1 - white, else 0 - black
	} else if count == 1 {
		colorWhite = false
	} else {
		c.JSON(
			400,
			gin.H{
				"error": "Game full",
			},
		)
		return
	}
	playerIP[Bool2int(colorWhite)] = c.ClientIP()
	count++
	c.JSON(
		200,
		gin.H{
			"colorWhite":colorWhite,
		},
	)
	if colorWhite{
		fmt.Println(c.ClientIP() + " is playing white")
	} else {
		fmt.Println(c.ClientIP() + " is playing black")
	}
}
var conditionChan = make([]chan int, 2)
var queryWhite = true
var fromInput, toInput Coord
var legalMoves []Coord
func main() {
	var board [8][8]Field
	initialiseBoard(&board)
	var kingPosition [2]Coord
	kingPosition[0] = Coord{X:4, Y:0}
	kingPosition[1] = Coord{X:4, Y:7}
	conditionChan[0]=make(chan int, 1)
	conditionChan[1]=make(chan int, 1)
	router := gin.Default()
	router.POST("/inputFrom", inputCoord(&fromInput))
	router.POST("/inputTo", inputCoord(&toInput))
	router.POST("/joinGame", joinGame)
	router.GET("/getFigureBoard", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{
				"board": processToFigureBoard(board),
			},
		)
	})
	router.GET("/getLegalBoard", func(c *gin.Context) {
		c.JSON(
			200,
			gin.H{
				"board": processToLegalBoard(board, legalMoves),
			},
		)
	})
	doneChan := make(chan bool)
	router.GET("/doneProcessing", func(c *gin.Context) {
		isDone := <- doneChan
		c.JSON(
			200,
			gin.H{
				"isDone": isDone,
			},
		)
	})
	queryChan := make(chan bool)
	router.GET("/queryChanged", func(c *gin.Context) {
		fmt.Println("Someonse is checking queryChanged")
		queryWhite := <- queryChan
		fmt.Println("Someonse waited till querywhite got a value")
		c.JSON(
			200,
			gin.H{
				"queryWhite": queryWhite,
			},
		)
	})
	router.GET("/gameCondition", func(c *gin.Context) {
		var gameCondition int
		if playerIP[0] == c.ClientIP(){
			gameCondition = <-conditionChan[0]
		} else if playerIP[1] == c.ClientIP(){
			gameCondition = <-conditionChan[1]
		}
		c.JSON(
			200,
			gin.H{
				"condition": gameCondition,
			},
		)
	})
	conditionChan[1] <- 0
	go router.Run()
	time.Sleep(time.Second)
	defer fmt.Println("thanks for playing")
	//game loop
	for {
		//render the board
		simpleRender(processToFigureBoard(board))

		//input from user
		var from *Field
		inputField(&from, &fromInput, &board)

		//calculate legal moves
		legalMoves = getLegalMoves(board, from, queryWhite, kingPosition)

		doneChan <- true
		//render the board with legal moves
//		renderBoardWithLegalMoves(board, legalMoves)
		simpleRender(processToLegalBoard(board, legalMoves))

		var to *Field
		inputField(&to, &toInput, &board)

		fmt.Println(to)
		//move valid check, repeated input if needed
		correct := isMoveCorrect(board, from, to, kingPosition, legalMoves, queryWhite)
		fmt.Println(correct)
		for !correct {
			fmt.Println("Invalid move")
//			renderBoard(board)
//			from = inputFrom(&board, queryWhite)
//			legalMoves = getLegalMoves(board, from, queryWhite, kingPosition)
//			renderBoardWithLegalMoves(board, legalMoves)
//			to = inputTo(&board, queryWhite)
//			correct = isMoveCorrect(board, from, to, kingPosition, legalMoves, queryWhite)
		}
		//move the pieces
		if correct {
			move(&board, from, to)
			fmt.Println("moved pieces")
			if to.figure=='k' {
				a := Bool2int(*to.colorWhite)
				kingPosition[a].X = to.x
				kingPosition[a].Y = to.y
			}
			fmt.Println("1")
			//is there check to opponent?
			isCheck := isCheck(board, !queryWhite, kingPosition)
			if isCheck {
				//is there mate to opponent?
				fmt.Println("2")
				isMate, err := isMate(board, !queryWhite, kingPosition)
				handle(err)
				if isMate {
					renderBoard(board)
					declareWin(queryWhite)
					conditionChan[Bool2int(!queryWhite)] <- 2
					queryChan <- !queryWhite
					conditionChan[Bool2int(queryWhite)] <- 3
					fmt.Println("sent a winning condition code:")
				} else {
					fmt.Println("3")
					conditionChan[Bool2int(!queryWhite)] <- 1
				}
			} else {
				fmt.Println("4")
				fmt.Println("6")
				conditionChan[Bool2int(!queryWhite)] <- 0
				fmt.Println("5")
			}
		}
		fmt.Println("5")
		doneChan <- true
		//change the query
		fmt.Println("sent doneChan")
		queryWhite = !queryWhite
		queryChan <- queryWhite
		fmt.Println("sent queryChan")
	}
}
