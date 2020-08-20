package vector2d

import "math"

// Useful https://github.com/bugra/matrix/blob/master/matrix.go
// AND https://github.com/grahamjenson/pathby/blob/master/lib/transformations.rb
type Matrix []Vector

func New3x3Matrix(a, b, c, d, e, f, g, h, i float64) Matrix {
	return []Vector{
		Vector{a, d, g},
		Vector{b, e, h},
		Vector{c, f, i},
	}
}

func (m Matrix) Extract3x3Matrix() (a, b, c, d, e, f, g, h, i float64) {
	if len(m) != 3 {
		panic("Matrix not 3x3")
	}

	a = m[0][0]
	b = m[1][0]
	c = m[2][0]
	d = m[0][1]
	e = m[1][1]
	f = m[2][1]
	g = m[0][2]
	h = m[1][2]
	i = m[2][2]
	return
}

func (m Matrix) Determinant() float64 {
	// https://www.mathsisfun.com/algebra/matrix-determinant.html
	a, b, c, d, e, f, g, h, i := m.Extract3x3Matrix()
	return (a * (e*i - f*h)) - (b * (d*i - f*g)) + (c * (d*h - e*g))
}

func (m Matrix) MinorMatrix() Matrix {
	a, b, c, d, e, f, g, h, i := m.Extract3x3Matrix()

	return New3x3Matrix(
		(e*i)-(f*h), // a
		(d*i)-(f*g), // b
		(d*h)-(e*g), // c
		(b*i)-(c*h), // d
		(a*i)-(c*g), // e
		(a*h)-(b*g), // f
		(b*f)-(c*e), // g
		(a*f)-(c*d), // h
		(a*e)-(b*d), // i
	)
}

func (m Matrix) CofactorMatrix() Matrix {
	a, b, c, d, e, f, g, h, i := m.Extract3x3Matrix()

	return New3x3Matrix(
		a, -b, c,
		-d, e, -f,
		g, -h, i,
	)
}

func (m Matrix) Transpose() Matrix {
	a, b, c, d, e, f, g, h, i := m.Extract3x3Matrix()

	return New3x3Matrix(
		a, d, g,
		b, e, h,
		c, f, i,
	)
}

func (m Matrix) MultiplyByScalar(s float64) Matrix {
	a, b, c, d, e, f, g, h, i := m.Extract3x3Matrix()

	return New3x3Matrix(
		a*s, b*s, c*s,
		d*s, e*s, f*s,
		g*s, h*s, i*s,
	)
}

func (m Matrix) Inverse() Matrix {
	// Only work on
	if len(m) != 3 {
		panic("Matrix not 3x3")
	}

	// https://www.mathsisfun.com/algebra/matrix-inverse-minors-cofactors-adjugate.html
	determinant := m.Determinant()
	minorMatrix := m.MinorMatrix()
	coFactor := minorMatrix.CofactorMatrix()
	adj := coFactor.Transpose()

	return adj.MultiplyByScalar(1.0 / determinant)
}

func (m1 Matrix) Product(m2 Matrix) Matrix {
	// A.B A must have exactly 3 Columns
	a, b, c, d, e, f, g, h, i := m1.Extract3x3Matrix()

	nm := make(Matrix, len(m2))
	for j, v := range m2 {
		nm[j] = Vector{
			a*v[0] + b*v[1] + c*v[2],
			d*v[0] + e*v[1] + f*v[2],
			g*v[0] + h*v[1] + i*v[2],
		}
	}
	return nm
}

////
// Specific Matricies
////

func RotateMatrix(radian float64) Matrix {
	return New3x3Matrix(
		math.Cos(radian), math.Sin(radian), 0,
		-math.Sin(radian), math.Cos(radian), 0,
		0, 0, 1,
	)
}

func TranslateMatrix(dx, dy float64) Matrix {
	return New3x3Matrix(
		1, 0, dx,
		0, 1, dy,
		0, 0, 1,
	)
}

func ScaleMatrix(sx, sy float64) Matrix {
	return New3x3Matrix(
		sx, 0, 0,
		0, sy, 0,
		0, 0, 1,
	)
}

func IdentityMatrix() Matrix {
	return New3x3Matrix(
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	)
}

// This returns a matrix that will transform about a point (e.g. center of shape)
func CenterTransformationAt(transofmrationMatrix Matrix, x, y float64) Matrix {
	translateTo := TranslateMatrix(-x, -y)
	// A^-1 * (B * A)
	return translateTo.Inverse().Product(transofmrationMatrix.Product(translateTo))
}
