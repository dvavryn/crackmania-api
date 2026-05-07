package ai

import (
	"fmt"
	"math"
)

func catmullRom(p0, p1, p2, p3 [3]float64, t float64) [3]float64 {
	t2 := t * t
	t3 := t2 * t

	return [3]float64{
		0.5 * ((2 * p1[0]) + (-p0[0]+p2[0])*t + (2*p0[0]-5*p1[0]+4*p2[0]-p3[0])*t2 + (-p0[0]+3*p1[0]-3*p2[0]+p3[0])*t3),
		0.5 * ((2 * p1[1]) + (-p0[1]+p2[1])*t + (2*p0[1]-5*p1[1]+4*p2[1]-p3[1])*t2 + (-p0[1]+3*p1[1]-3*p2[1]+p3[1])*t3),
		0.5 * ((2 * p1[2]) + (-p0[2]+p2[2])*t + (2*p0[2]-5*p1[2]+4*p2[2]-p3[2])*t2 + (-p0[2]+3*p1[2]-3*p2[2]+p3[2])*t3),
	}
}

func closestT(p0, p1, p2, p3, car [3]float64) float64 {
	out := 0.0
	bestDist := math.MaxFloat64
	for i := 0; i <= 40; i++ {
		t := float64(i) / 40
		b := catmullRom(p0, p1, p2, p3, t)
		dx := b[0] - car[0]
		dz := b[2] - car[2]
		d := dx*dx + dz*dz
		if d < bestDist {
			bestDist = d
			out = t
		}
	}
	return out
}

// func Calculate(f AIFrame, cnf config.Config, buildMode string) AIResponse { // old -- config!!!
func Calculate(f AIFrame, buildMode string) AIResponse {
	// Unpack frame fields
	px, py, pz := f.Position[0], f.Position[1], f.Position[2]
	qx, qy, qz, qw := f.Quaternion[0], f.Quaternion[1], f.Quaternion[2], f.Quaternion[3]
	vx, vy, vz := f.Velocity[0], f.Velocity[1], f.Velocity[2]
	cc, ct := f.CurrentCheckpoint, f.TotalCheckpoints
	cp0, cp1, cp2 := f.Checkpoints["0"], f.Checkpoints["1"], f.Checkpoints["2"]

	// Derived motion values
	speed := math.Sqrt(vx*vx + vy*vy + vz*vz)
	fx, fz := 2*(qx*qz+qw*qy), 1-2*(qx*qx+qy*qy)

	// Build spline segment and find closes point
	p1, p2, p3 := f.Checkpoints["0"], f.Checkpoints["1"], f.Checkpoints["2"]
	p0 := [3]float64{2*p1[0] - p2[0], 2*p1[1] - p2[1], 2*p1[2] - p2[2]}
	car := [3]float64{px, py, pz}

	tCar := closestT(p0, p1, p2, p3, car)
	tTarget := math.Min(tCar+0.35, 0.99)
	target := catmullRom(p0, p1, p2, p3, tTarget)

	// Compute steering angle
	dx := target[0] - px
	dz := target[2] - pz
	dist := math.Sqrt(dx*dx + dz*dz)
	if dist > 0 {
		dx, dz = dx/dist, dz/dist
	}
	dot := fx*dx + fz*dz
	cross := fz*dx - fx*dz
	angle := math.Atan2(cross, dot)

	// Apply controls
	var out AIResponse
	if angle > 0.2 {
		out.Right = true
	} else if angle < -0.2 {
		out.Left = true
	}
	if speed < 80 {
		out.Forward = true
	}

	// Debug output
	if buildMode == "debug" {
		pf := fmt.Printf
		pf("\n\n\n\n\n\n\n\n\n\n")
		pf("pos\t%v\t%v\t%v\n", px, py, pz)
		pf("quat\t%v\t%v\t%v\t%v\n", qx, qy, qz, qw)
		pf("vel\t%v\t%v\t%v\n", vx, vy, vz)
		pf("cp\t%v / %v\n", cc, ct)
		pf("cp+0\t%v\t%v\t%v\n", cp0[0], cp0[1], cp0[2])
		pf("cp+1\t%v\t%v\t%v\n", cp1[0], cp1[1], cp1[2])
		pf("cp+2\t%v\t%v\t%v\n", cp2[0], cp2[1], cp2[2])
		pf("speed\t%.1f\n", speed)
		pf("ang: %v\n", angle)
	}
	return out
}
