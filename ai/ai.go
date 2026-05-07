package ai

import (
	"fmt"
	"gotta-go-fast-api/config"
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
	best := 0.0
	bestDist := math.MaxFloat64
	for i := 0; i <= 40; i++ {
		t := float64(i) / 40
		b := catmullRom(p0, p1, p2, p3, t)
		dx := b[0] - car[0]
		dz := b[2] - car[2]
		d := dx*dx + dz*dz
		if d < bestDist {
			bestDist = d
			best = t
		}
	}
	return best
}

func Calculate(f AIFrame, cnf config.Config) AIResponse {
	var out AIResponse
	for i := range 15 {
		fmt.Println()
		_ = i
	}

	// position x, y, z
	px, py, pz := f.Position[0], f.Position[1], f.Position[2]
	fmt.Printf("pos\t%v\t%v\t%v\n", px, py, pz)

	// quaternion x, y, z, w
	qx, qy, qz, qw := f.Quaternion[0], f.Quaternion[1], f.Quaternion[2], f.Quaternion[3]
	fmt.Printf("quat\t%v\t%v\t%v\t%v\n", qx, qy, qz, qw)

	// velocity x, y, z
	vx, vy, vz := f.Velocity[0], f.Velocity[1], f.Velocity[2]
	fmt.Printf("vel\t%v\t%v\t%v\n", vx, vy, vz)

	// current checkpoint, total checkpoints
	cc, ct := f.CurrentCheckpoint, f.TotalCheckpoints
	fmt.Printf("cp\t%v / %v\n", cc, ct)

	var cp0 []float64
	cp0 = append(cp0, f.Checkpoints["0"][0], f.Checkpoints["0"][1], f.Checkpoints["0"][2])
	fmt.Printf("cp+0\t%v\t%v\t%v\n", cp0[0], cp0[1], cp0[2])

	var cp1 []float64
	cp1 = append(cp1, f.Checkpoints["1"][0], f.Checkpoints["1"][1], f.Checkpoints["1"][2])
	fmt.Printf("cp+1\t%v\t%v\t%v\n", cp1[0], cp1[1], cp1[2])

	var cp2 []float64
	cp2 = append(cp2, f.Checkpoints["2"][0], f.Checkpoints["2"][1], f.Checkpoints["2"][2])
	fmt.Printf("cp+2\t%v\t%v\t%v\n", cp2[0], cp2[1], cp2[2])

	speed := math.Sqrt(vx*vx + vy*vy + vz*vz)
	fmt.Printf("speed\t%.1f\n", speed)

	fx := 2 * (qx*qz + qw*qy)
	fz := 1 - 2*(qx*qx+qy*qy)

	p1 := [3]float64{f.Checkpoints["0"][0], f.Checkpoints["0"][1], f.Checkpoints["0"][2]}
	p2 := [3]float64{f.Checkpoints["1"][0], f.Checkpoints["1"][1], f.Checkpoints["1"][2]}
	p3 := [3]float64{f.Checkpoints["2"][0], f.Checkpoints["2"][1], f.Checkpoints["2"][2]}
	p0 := [3]float64{2*p1[0] - p2[0], 2*p1[1] - p2[1], 2*p1[2] - p2[2]}

	car := [3]float64{px, py, pz}
	tCar := closestT(p0, p1, p2, p3, car)
	tTarget := math.Min(tCar+0.35, 0.99)

	target := catmullRom(p0, p1, p2, p3, tTarget)
	dx := target[0] - px
	dz := target[2] - pz

	dist := math.Sqrt(dx*dx + dz*dz)
	if dist > 0 {
		dx, dz = dx/dist, dz/dist
	}

	dot := fx*dx + fz*dz
	cross := fz*dx - fx*dz
	angle := math.Atan2(cross, dot)

	fmt.Printf("\nang: %v\n", angle)
	if angle > 0.2 {
		out.Right = true
	} else if angle < -0.2 {
		out.Left = true
	}

	// out.Forward = dot > 0
	// out.Backward = dot <= 0 && speed < float64(cnf.DriveTrain.ReverseThreshold)
	if speed < 80 {
		out.Forward = true
	}

	return out
}
