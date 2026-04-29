package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var buildMode = "release"

type AIFrame struct {
	Position          [3]float64            `json:"position"`
	Quaternion        [4]float64            `json:"quaternion"`
	Velocity          [3]float64            `json:"velocity"`
	TotalCheckpoints  int                   `json:"totalCheckpoints"`
	CurrentCheckpoint int                   `json:"currentCheckpoint"`
	Checkpoints       map[string][3]float64 `json:"checkpoints"`
}

type AIResponse struct {
	Forward  bool `json:"forward"`
	Backward bool `json:"backward"`
	Left     bool `json:"left"`
	Right    bool `json:"right"`
	// Commands []string `json:"commands:`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func decideCommands(frame AIFrame) AIResponse {
	var out AIResponse
	out.Forward = true
	return out
}

func wsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		fmt.Print(".")
		var frame AIFrame
		if err := json.Unmarshal(msg, &frame); err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

		commands := decideCommands(frame)

		resp, _ := json.Marshal(commands)
		if err := conn.WriteMessage(websocket.TextMessage, resp); err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}

func main() {
	if buildMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.GET("/AI", wsHandler)

	r.Run(":8765")
}
