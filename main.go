package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var buildMode = "release"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func decideCommands(f AIFrame) AIResponse {
	var out AIResponse
	// p := fmt.Println

	/*====================
	Senior-Dev:
		1. Compute forward vector from quaternion
		2. Compute angle to next checkpoint
		3. If angle > threshold_right                 -> right     true
		   If angle < threshold_left                  -> left      true
		4. Compute speed = length(vx, vy, vz)
		5. Compute distance to next chckpoint
		6. If sharp turn ahead AND speed > safe_speed -> backwards true
		   Else                                       -> forwards  true
	====================*/

	var d D
	d.vx, d.vy, d.vz = f.Velocity[0], f.Velocity[1], f.Velocity[2]
	d.px, d.py, d.pz = f.Position[0], f.Position[1], f.Position[2]
	d.qx, d.qy, d.qz, d.qw = f.Quaternion[0], f.Quaternion[1], f.Quaternion[2], f.Quaternion[3]
	d.fx, d.fy, d.fz = 2*(d.qx*d.qz+d.qw*d.qy), 2*(d.qy*d.qz+d.qw*d.qx), 1-2*(d.qx*d.qx+d.qy*d.qy)

	d.dx, d.dz = f.Checkpoints["1"][0]-d.px, f.Checkpoints["1"][2]-d.pz

	cross := d.fx*d.dz - d.fz*d.dx
	dot := d.fx*d.dx + d.fz*d.dz
	angle := math.Atan2(cross, dot)
	if angle > 0.1 {
		out.Right = true
	} else if angle < -0.1 {
		out.Left = true
	}

	speed := math.Sqrt(d.vx*d.vx + d.vy*d.vy + d.vz*d.vz)
	if speed > 5 {
	} else {
		out.Forward = true
	}
	// p(speed)
	// ============
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
		// fmt.Println(string(msg))
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
