package ai

type AIResponse struct {
	Forward  bool `json:"forward"`
	Backward bool `json:"backward"`
	Left     bool `json:"left"`
	Right    bool `json:"right"`
}

type AIFrame struct {
	Position          [3]float64            `json:"position"`          // y, x, z
	Quaternion        [4]float64            `json:"quaternion"`        //
	Velocity          [3]float64            `json:"velocity"`          // y, x, z
	TotalCheckpoints  int                   `json:"totalCheckpoints"`  //
	CurrentCheckpoint int                   `json:"currentCheckpoint"` //
	Checkpoints       map[string][3]float64 `json:"checkpoints"`       // y, x, z
}

type D struct {
	vx, vy, vz     float64
	px, py, pz     float64
	qx, qy, qz, qw float64
	fx, fy, fz     float64

	dx, dz float64
}
