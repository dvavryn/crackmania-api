package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gotta-go-fast-api/ai"
	"gotta-go-fast-api/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var buildMode = "release"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
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

		var frame ai.AIFrame
		if err := json.Unmarshal(msg, &frame); err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

		var cnf config.Config
		cnf, err = config.GetConfig()
		if err != nil {
			fmt.Println("Parse error:", err)
			continue
		}

		commands := ai.Calculate(frame, cnf)

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
