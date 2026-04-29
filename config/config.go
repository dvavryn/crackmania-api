package config

type Car struct {
	Mass        int
	Size        [3]float64
	WhileRadius float64
}

type Engine struct {
force int
differential int
revForceRatio int
brakeForce int
coastBrakeForce int
frontBrakeBias int
reverseThreshold int
}

type Steering struct {
	MaxSteer float64
	MinSteer float64
	SteerSpeedMax int
}

type Suspension