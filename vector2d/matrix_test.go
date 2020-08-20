package vector2d

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New3x3Matrix(t *testing.T) {
	assert.NoError(t, nil)
}

func Test_New3x3Matrix_Extract3x3Matrix(t *testing.T) {
	var a, b, c, d, e, f, g, h, i float64 = 1, 2, 3, 4, 5, 6, 7, 8, 9
	m := New3x3Matrix(a, b, c, d, e, f, g, h, i)
	a1, b1, c1, d1, e1, f1, g1, h1, i1 := m.Extract3x3Matrix()
	assert.Equal(t, a, a1)
	assert.Equal(t, b, b1)
	assert.Equal(t, c, c1)
	assert.Equal(t, d, d1)
	assert.Equal(t, e, e1)
	assert.Equal(t, f, f1)
	assert.Equal(t, g, g1)
	assert.Equal(t, h, h1)
	assert.Equal(t, i, i1)
}

func Test_Determinant(t *testing.T) {
	// from https://www.mathsisfun.com/algebra/matrix-determinant.html
	m := New3x3Matrix(
		6, 1, 1,
		4, -2, 5,
		2, 8, 7,
	)
	assert.Equal(t, m.Determinant(), -306.0)
}

func Test_MinorMatrix(t *testing.T) {
	// FROM https://www.mathsisfun.com/algebra/matrix-inverse-minors-cofactors-adjugate.html
	m := New3x3Matrix(
		3, 0, 2,
		2, 0, -2,
		0, 1, 1,
	)
	m1 := m.MinorMatrix()

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	var a2, b2, c2, d2, e2, f2, g2, h2, i2 float64 = 2, 2, 2, -2, 3, 3, 0, -10, 0

	assert.Equal(t, a2, a1)
	assert.Equal(t, b2, b1)
	assert.Equal(t, c2, c1)
	assert.Equal(t, d2, d1)
	assert.Equal(t, e2, e1)
	assert.Equal(t, f2, f1)
	assert.Equal(t, g2, g1)
	assert.Equal(t, h2, h1)
	assert.Equal(t, i2, i1)
}

func Test_CofactorMatrix(t *testing.T) {
	// FROM https://www.mathsisfun.com/algebra/matrix-inverse-minors-cofactors-adjugate.html
	m := New3x3Matrix(
		2, 2, 2,
		-2, 3, 3,
		0, -10, 0,
	)
	m1 := m.CofactorMatrix()

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	var a2, b2, c2, d2, e2, f2, g2, h2, i2 float64 = 2, -2, 2, 2, 3, -3, 0, 10, 0

	assert.Equal(t, a2, a1)
	assert.Equal(t, b2, b1)
	assert.Equal(t, c2, c1)
	assert.Equal(t, d2, d1)
	assert.Equal(t, e2, e1)
	assert.Equal(t, f2, f1)
	assert.Equal(t, g2, g1)
	assert.Equal(t, h2, h1)
	assert.Equal(t, i2, i1)
}

func Test_Transpose(t *testing.T) {
	// FROM https://www.mathsisfun.com/algebra/matrix-inverse-minors-cofactors-adjugate.html
	m := New3x3Matrix(
		2, -2, 2,
		2, 3, -3,
		0, 10, 0,
	)
	m1 := m.Transpose()

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	var a2, b2, c2, d2, e2, f2, g2, h2, i2 float64 = 2, 2, 0, -2, 3, 10, 2, -3, 0

	assert.Equal(t, a2, a1)
	assert.Equal(t, b2, b1)
	assert.Equal(t, c2, c1)
	assert.Equal(t, d2, d1)
	assert.Equal(t, e2, e1)
	assert.Equal(t, f2, f1)
	assert.Equal(t, g2, g1)
	assert.Equal(t, h2, h1)
	assert.Equal(t, i2, i1)
}

func Test_MultiplyByScalar(t *testing.T) {
	m := New3x3Matrix(
		2, 3, 4,
		5, 6, 7,
		8, 9, 10,
	)
	m1 := m.MultiplyByScalar(2)

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	var a2, b2, c2, d2, e2, f2, g2, h2, i2 float64 = 4, 6, 8, 10, 12, 14, 16, 18, 20

	assert.Equal(t, a2, a1)
	assert.Equal(t, b2, b1)
	assert.Equal(t, c2, c1)
	assert.Equal(t, d2, d1)
	assert.Equal(t, e2, e1)
	assert.Equal(t, f2, f1)
	assert.Equal(t, g2, g1)
	assert.Equal(t, h2, h1)
	assert.Equal(t, i2, i1)
}

func Test_Inverse(t *testing.T) {
	// FROM https://www.mathsisfun.com/algebra/matrix-inverse-minors-cofactors-adjugate.html

	m := New3x3Matrix(
		3, 0, 2,
		2, 0, -2,
		0, 1, 1,
	)

	m1 := m.Inverse()

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	var a2, b2, c2, d2, e2, f2, g2, h2, i2 float64 = 0.2, 0.2, 0, -0.2, 0.3, 1, 0.2, -0.3, 0

	assert.InDelta(t, a2, a1, 0.000001)
	assert.InDelta(t, b2, b1, 0.000001)
	assert.InDelta(t, c2, c1, 0.000001)
	assert.InDelta(t, d2, d1, 0.000001)
	assert.InDelta(t, e2, e1, 0.000001)
	assert.InDelta(t, f2, f1, 0.000001)
	assert.InDelta(t, g2, g1, 0.000001)
	assert.InDelta(t, h2, h1, 0.000001)
	assert.InDelta(t, i2, i1, 0.000001)
}

func Test_Product_Matrix(t *testing.T) {
	m1 := New3x3Matrix(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)

	m2 := New3x3Matrix(
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	)

	p := m1.Product(m2)

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := p.Extract3x3Matrix()

	assert.Equal(t, 30.0, a1)
	assert.Equal(t, 36.0, b1)
	assert.Equal(t, 42.0, c1)
	assert.Equal(t, 66.0, d1)
	assert.Equal(t, 81.0, e1)
	assert.Equal(t, 96.0, f1)
	assert.Equal(t, 102.0, g1)
	assert.Equal(t, 126.0, h1)
	assert.Equal(t, 150.0, i1)
}

func Test_Product_IdentityMatrix(t *testing.T) {
	i := IdentityMatrix()
	m1 := New3x3Matrix(
		3, 0, 2,
		2, 0, -2,
		0, 1, 1,
	)

	a1, b1, c1,
		d1, e1, f1,
		g1, h1, i1 := m1.Extract3x3Matrix()

	m2 := m1.Product(i)
	a2, b2, c2,
		d2, e2, f2,
		g2, h2, i2 := m2.Extract3x3Matrix()

	assert.Equal(t, a2, a1)
	assert.Equal(t, b2, b1)
	assert.Equal(t, c2, c1)
	assert.Equal(t, d2, d1)
	assert.Equal(t, e2, e1)
	assert.Equal(t, f2, f1)
	assert.Equal(t, g2, g1)
	assert.Equal(t, h2, h1)
	assert.Equal(t, i2, i1)
}
