package ai

import (
	"gotta-go-fast-api/config"
	"math"
)

// func Calculate(f AIFrame, cnf config.Config) AIResponse {
// 	var out AIResponse
// 	// p := fmt.Println

// 	/*====================
// 	Senior-Dev:
// 		1. Compute forward vector from quaternion
// 		2. Compute angle to next checkpoint
// 		3. If angle > threshold_right                 -> right     true
// 		   If angle < threshold_left                  -> left      true
// 		4. Compute speed = length(vx, vy, vz)
// 		5. Compute distance to next chckpoint
// 		6. If sharp turn ahead AND speed > safe_speed -> backwards true
// 		   Else                                       -> forwards  true
// 	====================*/

// 	var d D
// 	d.vx, d.vy, d.vz = f.Velocity[0], f.Velocity[1], f.Velocity[2]
// 	// d.px, d.py, d.pz = f.Position[0], f.Position[1], f.Position[2]
// 	// d.qx, d.qy, d.qz, d.qw = f.Quaternion[0], f.Quaternion[1], f.Quaternion[2], f.Quaternion[3]
// 	// d.fx, d.fy, d.fz = 2*(d.qx*d.qz+d.qw*d.qy), 2*(d.qy*d.qz+d.qw*d.qx), 1-2*(d.qx*d.qx+d.qy*d.qy)

// 	// d.dx, d.dz = f.Checkpoints["1"][0]-d.px, f.Checkpoints["1"][2]-d.pz

// 	// cross := d.fx*d.dz - d.fz*d.dx
// 	// dot := d.fx*d.dx + d.fz*d.dz
// 	// angle := math.Atan2(cross, dot)
// 	// if angle > 0.1 {
// 	// 	out.Right = true
// 	// } else if angle < -0.1 {
// 	// 	out.Left = true
// 	// }

// 	speed := math.Sqrt(d.vx*d.vx + d.vy*d.vy + d.vz*d.vz)
// 	if speed > 5 {
// 	} else {
// 		out.Forward = true
// 	}
// 	out.Right = true
// 	// p(speed)
// 	// ============
// 	return out
// }

func catmullRom(p0, p1, p2, p3 [3]float64, t float64) [3]float64 {
	t2 := t * t
	t3 := t2 * t
	return [3]float64{
		0.5 * ((2 * p1[0]) + (-p0[0]+p2[0])*t + (2*p0[0]-5*p1[0]+4*p2[0]-p3[0])*t2 + (-p0[0]+3*p1[0]-3*p2[0]+p3[0])*t3),
		0.5 * ((2 * p1[1]) + (-p0[1]+p2[1])*t + (2*p0[1]-5*p1[1]+4*p2[1]-p3[1])*t2 + (-p0[1]+3*p1[1]-3*p2[1]+p3[1])*t3),
		0.5 * ((2 * p1[2]) + (-p0[2]+p2[2])*t + (2*p0[2]-5*p1[2]+4*p2[2]-p3[2])*t2 + (-p0[2]+3*p1[2]-3*p2[2]+p3[2])*t3)
	}
}

func closestT(p0, p1, p2, p3, car [3]float64) float64 {
	best := 0.0
	bestDist := math.MaxFloat64
	for i := 0; i <= 40; i++ {
		t := float64(i) / 40
		b := catmullRom(p0, p1, p2, t)
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
	vx := f.Velocity[0]
	vy := f.Velocity[1]
	vz := f.Velocity[2]

	px := f.Position[0]
	pz := f.Position[2]

	qx := f.Quaternion[0]
	qy := f.Quaternion[1]
	qz := f.Quaternion[2]
	qw := f.Quaternion[3]

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
	cross := fx*dz - fz*dx
	angle := math.Atan2(cross, dot)

	if angle > 0.1 {
		out.Right = true
	} else if angle < -0.1 {
		out.Left = true
	}

	speed := math.Sqrt(vx*vx + vy*vy + vz*vz)
	out.Forward = dot > 0
	out.Backward = dot <= 0 && speed < float64(cnf.DriveTrain.ReverseThreshold)

	return out
}

/*

Algorithm:
	1. Read frame-input:
		position, quaternion, velocity
	2. Extract forward vector
		rotate world fwd by quaterion
	3. Vector to target checkpoint
		nectCheckpoint - position
	4. Dot product
		fwd * cp_normalized -> ahead or behind?
	5. Signed steer error
		atan2(cross, dot) -> left or right?
	6. Forward speed
		speed * dot -> moving forward or back?
	7. Map to outputs
		forward
			dot > 0 & fwdSpeed ok
		backward
			dot < 0 or stuck
		left
			steerErr < -deadzone
		right steerErr > +deadzone
	8. Emit frame-output
	9. step 1

*/
