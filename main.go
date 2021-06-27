package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

type Game struct {
	ID      string `json:"id"`
	Timeout int32  `json:"timeout"`
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Battlesnake struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Health int32   `json:"health"`
	Body   []Coord `json:"body"`
	Head   Coord   `json:"head"`
	Length int32   `json:"length"`
	Shout  string  `json:"shout"`
}

type Board struct {
	Height int           `json:"height"`
	Width  int           `json:"width"`
	Food   []Coord       `json:"food"`
	Snakes []Battlesnake `json:"snakes"`
}

type BattlesnakeInfoResponse struct {
	APIVersion string `json:"apiversion"`
	Author     string `json:"author"`
	Color      string `json:"color"`
	Head       string `json:"head"`
	Tail       string `json:"tail"`
}

type GameRequest struct {
	Game  Game        `json:"game"`
	Turn  int         `json:"turn"`
	Board Board       `json:"board"`
	You   Battlesnake `json:"you"`
}

type MoveResponse struct {
	Move  string `json:"move"`
	Shout string `json:"shout,omitempty"`
}

func AddCoord (a Coord,b Coord)Coord {
     newX := a.X + b.X
     newY := a.Y + b.Y
     var newCoord Coord
     newCoord.X = newX
     newCoord.Y = newY
    return newCoord
}

func GenWalls(request GameRequest) [2]int{

  var boardMax int = request.Board.Height 
  var boardMin int = 0
  var walls [2]int = [2]int{boardMin,boardMax}
  return walls
}

func GetMoveResult(move string, head Coord) Coord{
  var moveResult Coord
  var possibleResult [4]Coord = [4]Coord{{0,1},{0,-1},{-1,0},{1,0}}

  switch move {
    case "up":
        moveResult = AddCoord(head,possibleResult[0])
    case "down":
        moveResult = AddCoord(head,possibleResult[1])
    case "left":
        moveResult = AddCoord(head,possibleResult[2])
    case "right":
        moveResult = AddCoord(head,possibleResult[3])    
    }

  return moveResult
}

func CheckMove(moveResult Coord, walls [2]int)bool{
  //if walls in result false else true
  var result bool
  result = true
  for i:=0;i<1;i++{
    if (moveResult.X == walls[i] || moveResult.Y ==walls[i]){
      result = false
    }
  }
  return result
}

//To do get head location and avoid walls and self


// HandleIndex is called when your Battlesnake is created and refreshed
// by play.battlesnake.com. BattlesnakeInfoResponse contains information about
// your Battlesnake, including what it should look like on the game board.
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "smolgeat",
		Color:      "#4d0019",
		Head:       "#316000",
		Tail:       "coffee",
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleStart is called at the start of each game your Battlesnake is playing.
// The GameRequest object contains information about the game that's about to start.
// TODO: Use this function to decide how your Battlesnake is going to look on the board.
func HandleStart(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("START\n")
}

// HandleMove is called for each turn of each game.
// Valid responses are "up", "down", "left", or "right".
// TODO: Use the information in the GameRequest object to determine your next move.
func HandleMove(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}
  var finalMove string
	// Choose a random direction to move in
	possibleMoves := []string{"up", "down", "left", "right"}
  for {
	move := possibleMoves[rand.Intn(len(possibleMoves))]
  walls := GenWalls(request)
  moveResult := GetMoveResult(move, request.You.Head)
  result :=CheckMove(moveResult, walls)
  if (result == true){
    finalMove = move
    break
  }
  }
	response := MoveResponse{
		Move: finalMove,
	}

	fmt.Printf("MOVE: %s\n", response.Move)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Fatal(err)
	}
}

// HandleEnd is called when a game your Battlesnake was playing has ended.
// It's purely for informational purposes, no response required.
func HandleEnd(w http.ResponseWriter, r *http.Request) {
	request := GameRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Fatal(err)
	}

	// Nothing to respond with here
	fmt.Print("END\n")
}

func main() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/start", HandleStart)
	http.HandleFunc("/move", HandleMove)
	http.HandleFunc("/end", HandleEnd)

	fmt.Printf("Starting Battlesnake Server at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
