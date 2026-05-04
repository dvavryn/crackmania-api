package config

type Chassis struct {
	Mass        float64 `json:"mass"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	Length      float64 `json:"length"`
	WheelRadius float64 `json:"wheelRadius"`
}

type DriveTrain struct {
	Force            float64 `json:"force"`
	Differential     float64 `json:"differential"`   // fraction fof force sent to read (0 = FWD, 1 = RWD, 0.5 = AWD)
	RevForceRation   float64 `json:"revForceRation"` // reverse force as a fraction of total force
	BrakeForce       float64 `json:"brakeForce"`
	CoastBrakeForce  float64 `json:"coastBrakeForce"`  // passive develeration applied when no throttle/brake input
	FrontBrakeBias   float64 `json:"frontBrakeBias"`   // front brake fraction when turning (reduces udnersteer)
	ReverseThreshold float64 `json:"reverseThreshold"` // forward speed (m/s) below which braking switches to reverse
}

type Steering struct {
	MaxSteer      float64 `json:"maxSteer"`
	MinSteer      float64 `json:"minSteer"`
	SteerSpeedMax float64 `json:"steerSpeedMax"` // speed (m/s) at which steering is clamped to minSteer
}

type Suspension struct {
	SuspensionStiffness  float64 `json:"suspensionStiffness"`
	SuspensionRestLength float64 `json:"suspensionRestLength"`
	DampingRelaxation    float64 `json:"dampingRelaxation"`
	DampingCompression   float64 `json:"dampingCompression"`
	MaxSuspensionForce   float64 `json:"maxSuspensionForce"`
	RollInfluence        float64 `json:"rollInfluence"`
	MaxSuspensionTravel  float64 `json:"maxSuspensionTravel"`
}

type Wheels struct {
	RearFrictionSlip                float64 `json:"rearFrictionSlip"`
	FrontFrictionSlip               float64 `json:"frontFrictionSlip"`
	CustumSlidingRotationalSpeed    float64 `json:"custumSlidingRotationalSpeed"`
	UseCustomSlidingRotationalSpeed bool    `json:"useCustomSlidingRotationalSpeed"`
	WheelCylinderThickness          float64 `json:"wheelCylinderThickness"`
	WheelCylinderSegments           float64 `json:"wheelCylinderSegments"`
	WheelMass                       float64 `json:"wheelMass"`
}

type WheelGeometryFactors struct { // multiplied against chassis dimensions
	WheelHeightFactor         float64 `json:"wheelHeightFactor"`
	WheelWidthFactor          float64 `json:"wheelWidthFactor"`
	WheelFrontOffsetFactor    float64 `json:"wheelFrontOffsetFactor"`
	WheelRearOffsetFactor     float64 `json:"wheelRearOffsetFactor"`
	FrontWheelWidthMultiplier float64 `json:"frontWheelWidthMultiplier"`
}

type Config struct {
	Chassis              Chassis              `json:"chassis"`
	DriveTrain           DriveTrain           `json:"drivetrain"`
	Steering             Steering             `json:"steering"`
	Suspension           Suspension           `json:"suspension"`
	Wheels               Wheels               `json:"wheels"`
	WheelGeometryFactors WheelGeometryFactors `json:"wheelGeometryFactors"`
}
