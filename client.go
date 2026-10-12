package main
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"bytes"
	"os"
)
func fetchBoard(endpoint string) [8][8]rune{
	resp, err := http.Get(url+endpoint)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res BoardResponse
	json.Unmarshal(body, &res)
	board := res.Board
	return board
}
func renderBoard(board [8][8]rune) {
	fmt.Print("\033[2J\033[H")
	fmt.Println(" ┌───────────────┐")
	for y :=0; y<8; y++ {
		fmt.Print(8-y)
		fmt.Print(string('│'))
		for x:=0; x<8; x++ {
			fmt.Printf("%c", board[x][y])
			if x!=7 {
				fmt.Print(string(' '))
			}
		}
		fmt.Print(string('│') + "\n")
	}
	fmt.Println(" └───────────────┘")
	fmt.Println("  a b c d e f g h ")
}
type inputError string

func (e inputError) Error() string {
	return fmt.Sprintf(string(e))
}

func sendInput(endpoint string, x, y int) error {
	data := []byte(fmt.Sprintf(`{"x":%d,"y":%d}`, x, y))
	resp, err := do("POST", endpoint, data)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var res InputResponse
	json.Unmarshal(body, &res)
	if res.Status == "invalid" {
		return inputError("Invalid move")
	} else {
		 return nil
	}
}
type InputResponse struct {
	Status string `json: "status"`
}
type BoardResponse struct{
	Board [8][8]rune `json: "board"`
}
type QueryResponse struct {
	QueryWhite bool `json:"queryWhite"`
}
type ColorResponse struct {
	ColorWhite bool `json:"colorWhite"`
	Token string `json:"token"`
}
var token string
//same as http.Get/Post but sends our token, so server knows who we are
func do(method, endpoint string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url+endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", token)
	return http.DefaultClient.Do(req)
}
type ConditionResponse struct {
	Condition int `json:"condition"`
}
func getColor() bool {
	resp, err := http.Post(
		url+"/joinGame",
		"application/json",
		nil,
	)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res ColorResponse
	json.Unmarshal(body, &res)
	token = res.Token
	if resp.StatusCode == 200 {
		if res.ColorWhite {
			fmt.Println("You are playing white")
		} else {
			fmt.Println("You are playing black")
		}
	} else {
		fmt.Println("Error. Game is already full")
		os.Exit(1)
	}
	return res.ColorWhite
}
var colorWhite bool
func isMyQuery() bool {
	resp, err := http.Get(url+"/getQuery")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res QueryResponse
	json.Unmarshal(body, &res)
	queryWhite := res.QueryWhite
	return colorWhite == queryWhite
}
var url string
func waitProcessing() {
	http.Get(url+"/doneProcessing")
}
func waitQuery() {
	do("GET", "/queryChanged", nil)
}
func getGameCondition() int {
	resp, err := do("GET", "/gameCondition", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var res ConditionResponse
	json.Unmarshal(body, &res)
	return res.Condition
}
func main() {
	var ip, port string
	fmt.Println("Server IP:")
	fmt.Scan(&ip)
	fmt.Println("Port:")
	fmt.Scan(&port)
	url = "http://" + ip + ":" + port

	colorWhite = getColor()
	//game loop
	if !colorWhite {
		waitQuery()
	}
	var figureBoard, legalBoard [8][8]rune
	for {
		gameCondition := getGameCondition()
		if gameCondition==3 {
			fmt.Println("Mate.")
			fmt.Println("You won.")
			return
		}

		figureBoard = fetchBoard("/getFigureBoard")
		renderBoard(figureBoard)

		if gameCondition==1 {
			fmt.Println("Check.")
		} else if gameCondition==2 {
			fmt.Println("Mate.")
			fmt.Println("You lost.")
			return
		} else if gameCondition==4 {
			fmt.Println("Stalemate. It is a draw.")
			return
		}
		var f string
		fmt.Println("From:")
		fmt.Scan(&f)
		fx := int(f[0])-97
		fy := 56-int(f[1])
		sendInput("/inputFrom", fx, fy)
		waitProcessing()
		legalBoard = fetchBoard("/getLegalBoard")
		fmt.Println("10")
		renderBoard(legalBoard)

		var t string
		fmt.Println("To:")
		fmt.Scan(&t)
		tx := int(t[0])-97
		ty := 56-int(t[1])
		err := sendInput("/inputTo", tx, ty)
		for err!=nil {
			fmt.Println(err)
			fmt.Println("From:")
			fmt.Scan(&f)
			fx = int(f[0])-97
			fy = 56-int(f[1])
			sendInput("/inputFrom", fx, fy)
			fmt.Println("To:")
			fmt.Scan(&t)
			tx := int(t[0])-97
			ty := 56-int(t[1])
			err = sendInput("/inputTo", tx, ty)
		}
		waitProcessing()
		figureBoard = fetchBoard("/getFigureBoard")
		renderBoard(figureBoard)
		waitQuery()
	}
}
