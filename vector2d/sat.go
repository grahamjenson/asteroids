package vector2d

import (
	"math"
)

// From https://github.com/jriecken/sat-js/blob/master/SAT.js

type HitCheckResponse struct {
	overlap  float64 // amount of overlap
	overlapN Vector  // the normal of the overlap

	aInB bool // A is totally inside B
	bInA bool // B is totally inside A

	Collision bool   // Are the polygons colliding
	OverlapV  Vector // The vector of the overlap (the amount and direction A needs to move to be free)
}

func (p1 *Polygon) HitCheck(p2 *Polygon) *HitCheckResponse {
	response := &HitCheckResponse{
		Collision: false,
		aInB:      true,
		bInA:      true,
		overlap:   math.MaxFloat64,
	}

	// If any of the edge normals of A is a separating axis, no intersection.
	for _, e := range p1.Edges() {
		nResp := checkNormalForSeperatingAxis(p1, p2, e.Normal)
		if !nResp.Collision {
			return nResp
		}

		// min overlap
		if nResp.overlap < response.overlap {
			response.overlap = nResp.overlap
			response.overlapN = nResp.overlapN
		}
		response.aInB = response.aInB && nResp.aInB
		response.bInA = response.bInA && nResp.bInA
	}

	// If any of the edge normals of B is a separating axis, no intersection.
	for _, e := range p2.Edges() {
		nResp := checkNormalForSeperatingAxis(p1, p2, e.Normal)
		if !nResp.Collision {
			return nResp
		}

		// min overlap
		if nResp.overlap < response.overlap {
			response.overlap = nResp.overlap
			response.overlapN = nResp.overlapN
		}
		response.aInB = response.aInB && nResp.aInB
		response.bInA = response.bInA && nResp.bInA
	}

	response.Collision = true

	response.OverlapV = response.overlapN.Clone().Scale(response.overlap, response.overlap)

	return response
}

// SAT algorithm
// We project the vector onto the normal axis to check for separation
func checkNormalForSeperatingAxis(a, b *Polygon, normal Vector) *HitCheckResponse {
	response := &HitCheckResponse{
		Collision: false,

		aInB: false, // assume equality
		bInA: false, // assume equality

		overlap:  math.MaxFloat64, // We are looking for smallest overlap
		overlapN: normal.Clone(),
	}

	// Project the polygons onto the axis.
	minA, maxA := a.flattenPointsOn(normal)
	minB, maxB := b.flattenPointsOn(normal)

	// Check overlap of points
	if minA > maxB || minB > maxA {
		// No overlap on normal
		return response
	}

	response.Collision = true

	// We know they intersect so either:
	// 1. a is in b
	// 2. b is in a
	// 3. a is greater than b
	// 4. b is greater than a
	// 5. a and b are equal (default)

	// That is, either
	// Either  |---|--|---|
	//         a---b--a---b // B is greater than A
	//         a---b--b---a // B is inside A
	//         b---a--b---a // A is greater than B
	//         b---a--a---b // A is inside B

	// Direction Matters, e.g.
	// in a--------b-a-b the a will move left (<) which is the reverse direction
	// in b-a-b--------a then a would move right and the reverse is not needed
	response.aInB = true
	response.bInA = true

	if minA < minB {
		response.aInB = false
		// Either  |---|--|---|
		//         a---b--a---b // B is greater than A
		//         a---b--b---a // B is inside A
		if maxA < maxB {
			//         a---b--a---b // B is greater than A
			// overlap     ----
			response.overlap = maxA - minB
			response.overlapN.Reverse()
			response.bInA = false
			// B is fully inside A.  Pick the shortest way out.
		} else {
			//         a---b--b---a // B is inside A
			// overlap     --------
			//    or   --------
			// select smallest overlap
			option1 := maxA - minB
			option2 := maxB - minA
			if option1 < option2 {
				response.overlap = option1
				response.overlapN.Reverse()
			} else {
				response.overlap = option2
			}
		}
	} else {
		// Either
		//         b---a--b---a // A is greater than B
		//         b---a--a---b // A is inside B
		response.bInA = false
		if maxA > maxB {
			//         b---a--b---a // A is greater than B
			// overlap     ----
			response.overlap = maxB - minA
			response.aInB = false
		} else {
			//         b---a--a---b // A is inside B
			// overlap     --------
			//    or   --------
			// select smallest overlap
			option1 := maxB - minA
			option2 := maxA - minB
			if option1 < option2 {

				response.overlap = option1
			} else {
				response.overlap = option2
				response.overlapN.Reverse()
			}
		}
	}

	return response
}

func (s *Polygon) flattenPointsOn(normal Vector) (float64, float64) {
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, p := range s.Matrix {
		// The magnitude of the projection of the point onto the normal
		dot := p.Dot(normal)
		if dot < min {
			min = dot
		}
		if dot > max {
			max = dot
		}
	}

	return min, max
}
