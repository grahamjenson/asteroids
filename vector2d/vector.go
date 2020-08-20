package vector2d

import "math"

// Should always be of len(3) [x,y,1]
type Vector []float64

// Perpendicular
func (v Vector) Clone() Vector {
	return Vector{v[0], v[1], v[2]}
}

func (v Vector) Perp() Vector {
	x := v[0]
	v[0] = v[1]
	v[1] = -x
	return v
}

// Reverse this vector.
func (v Vector) Reverse() Vector {
	v[0] = -v[0]
	v[1] = -v[1]
	return v
}

func (v Vector) Scale(x, y float64) Vector {
	v[0] *= x
	v[1] *= y
	return v
}

// Normalize this vector.  (make it have length of `1`)
func (v Vector) Normalize() Vector {
	var d = v.len()
	if d > 0 {
		v[0] = v[0] / d
		v[1] = v[1] / d
	}
	return v
}

// Add another vector to this one.
func (v Vector) Add(other Vector) Vector {
	v[0] += other[0]
	v[1] += other[1]
	return v
}

// Add another vector to this one.
func (v Vector) MidPoint(other Vector) Vector {
	return Vector{(v[0] + other[0]) / 2, (v[1] + other[1]) / 2, 1}
}

// Subtract another vector from this one.
func (v Vector) Sub(other Vector) Vector {
	v[0] -= other[0]
	v[1] -= other[1]
	return v
}

// Get the dot product of this vector and another.
func (v Vector) Dot(other Vector) float64 {
	return v[0]*other[0] + v[1]*other[1]
}

// Get the squared length of this vector.
func (v Vector) len2() float64 {
	return v.Dot(v)
}

// Get the length of this vector.
func (v Vector) len() float64 {
	return math.Sqrt(v.len2())
}
