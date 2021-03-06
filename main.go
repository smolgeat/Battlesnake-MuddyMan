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

func AvoidSuicide(request GameRequest) (move string) {
	// take snake position
	// take walls position
	//  if move  causes head to hit it then new move

	myHead := request.You.Head
  myBody := request.You.Body
  opponents := request.Board.Snakes
  occupied := map[Coord]string{

  }
  //place opponent snake coordinates in a map
  for _, snake:= range opponents{
    for _,snakeBodySegment:= range snake.Body{
    
    occupied[snakeBodySegment] = "Opponent"
    }
  }
  for _, myBodySegment:= range myBody[1:]{
      occupied[myBodySegment] = "Me"
      }
	boundary := []int{-1, request.Board.Height}
	var newHead Coord
	newHeadP := &newHead
	// Choose a random direction to move in
	//map of moves and how they change head position
	moveResults := map[string][]int{
		"up":    {0, 1},
		"down":  {0, -1},
		"left":  {-1, 0},
		"right": {1, 0},
	}
	possibleMoves := []string{"up", "down", "left", "right"}

	moveP := &move
	safe := 0
	for safe != 1 {

		*moveP = possibleMoves[rand.Intn(len(possibleMoves))]
		// compute new head position
		switch *moveP {
		case "up":
			newHeadX := myHead.X + moveResults[move][0]
			newHeadY := myHead.Y + moveResults[move][1]
			*newHeadP = Coord{newHeadX, newHeadY}

		case "down":
			newHeadX := myHead.X + moveResults[move][0]
			newHeadY := myHead.Y + moveResults[move][1]
			*newHeadP = Coord{newHeadX, newHeadY}
		case "left":
			newHeadX := myHead.X + moveResults[move][0]
			newHeadY := myHead.Y + moveResults[move][1]
			*newHeadP = Coord{newHeadX, newHeadY}
		case "right":
			newHeadX := myHead.X + moveResults[move][0]
			newHeadY := myHead.Y + moveResults[move][1]
			*newHeadP = Coord{newHeadX, newHeadY}

		}
		// if newHead in myBody or boundray in newHead
		if (newHead.X == boundary[0]) || (newHead.X == boundary[1]) || (newHead.Y == boundary[0]) || (newHead.Y == boundary[1]) {
        fmt.Println(move,"IS WRONG, NEW MOVE")
		} else {
			safe = 1
		}// assigns result of occupied[newHead] to coordinate if no result then coordinate will be a zero value, exists is a boolean true if newHead position is occupied
    if coordinate, exists:=occupied[newHead];exists {
     fmt.Println(move,coordinate,"IS WRONG, NEW MOVE") 
    safe = 0
}
	}

	return move

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
		Head:       "viper",
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

  //timeout := request.Game.Timeout

	move := AvoidSuicide(request)
  fmt.Println(request.Turn)
	response := MoveResponse{
		Move: move,
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
