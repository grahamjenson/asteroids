package vector2d

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test visualizations here
// https://www.w3schools.com/code/tryit.asp?filename=GHUYC3ODBQBP
var CollisionTests = []struct {
	description string
	colliding   bool
	s1          []float64
	s2          []float64
	overlapV    Vector
}{
	{
		"totally disjoint shapes",
		false,
		[]float64{0, 10, 10, 10, 20, 00},
		[]float64{30, 40, 40, 40, 50, 30},
		Vector{},
	},
	{
		"touching a concave shape",
		true,
		[]float64{0, 0, 50, 0, 50, 50, 25, 10, 0, 50},
		[]float64{20, 20, 30, 20, 50, 40},
		Vector{0, -30, 1},
	},
	{
		"two squares overlapping",
		true,
		[]float64{0, 0, 10, 0, 10, 10, 0, 10},
		[]float64{5, 5, 15, 5, 15, 15, 5, 15},
		Vector{0, -5, 1},
	},
	{
		"convex shapes rotated",
		false,
		[]float64{0, 10, 40, 40, 50, 30, 10, 0},
		[]float64{30, 10, 70, 40, 80, 30, 40, 0},
		Vector{},
	},
	{
		"B in A",
		true,
		[]float64{0, 0, 0, 50, 50, 50, 50, 0},
		[]float64{10, 10, 10, 20, 20, 20, 20, 10},
		Vector{20, 0, 1},
	},
}

func Test_CollidingWith(t *testing.T) {
	for i, tt := range CollisionTests {
		t.Run(fmt.Sprintf("%v_%s", i, tt.description), func(t *testing.T) {
			s1 := NewPolygon(tt.s1...)
			s2 := NewPolygon(tt.s2...)
			hc := s1.HitCheck(s2)

			assert.Equal(t, tt.colliding, hc.Collision)
			if hc.Collision {
				assert.Equal(t, tt.overlapV, hc.OverlapV)
			}
		})
	}
}

//https://bell0bytes.eu/centroid-convex/
var CentroidTests = []struct {
	description string
	s           []float64
	centroidx   float64
	centroidy   float64
}{
	{
		"triangle",
		[]float64{0, 1, 0, 2, 3, 0},
		1,
		1,
	},
	{
		"convex",
		[]float64{1, 0, 2, 1, 0, 3, -1, 2, -2, -1},
		-0.08,
		0.92,
	},
}

func Test_Centroid(t *testing.T) {
	for i, tt := range CentroidTests {
		t.Run(fmt.Sprintf("%v_%s", i, tt.description), func(t *testing.T) {
			s := NewPolygon(tt.s...)
			x, y := s.Centroid()

			assert.InDelta(t, x, tt.centroidx, 0.01)
			assert.InDelta(t, y, tt.centroidy, 0.01)
		})
	}
}

var IsConvexTests = []struct {
	description string
	s           []float64
	is          bool
}{
	{
		"convex",
		[]float64{0, 10, 10, 10, 20, 00},
		true,
	},
	{
		"concave",
		[]float64{0, 0, 50, 0, 50, 50, 25, 10, 0, 50},
		false,
	},
	{
		"flat",
		[]float64{0, 0, 10, 0, 20, 0},
		true,
	},
}

func Test_IsConvex(t *testing.T) {
	for i, tt := range IsConvexTests {
		t.Run(fmt.Sprintf("%v_%s", i, tt.description), func(t *testing.T) {
			s := NewPolygon(tt.s...)
			assert.Equal(t, s.IsConvex(), tt.is)
		})
	}
}

// PROJECTION
// -671.5792057939544 5e-324 [-0.6199034636337281 -0.7846780841688566,[687.6316996852223,312.63003058710115,654.1102372351858,339.1123181622106,662.1094593091648,331.3012717399208,656.0111626926846,312.25367928372236,644.9627923104654,310.540929477913
// -875.7947618476903 5e-324 [-0.6199034636337281 -0.7846780841688566,[976.3706931856734,344.7772950585037,1040.227761300695,282.1492274549643,939.5312529304391,230.57197275835424,901.4628121637226,242.85178617676567,875.6741848154176,293.20004036189357,926.0224390005454,318.9886677101987

var FlattenPointsTests = []struct {
	description string
	s           []float64
	n           Vector
	min         float64
	max         float64
}{
	{
		"box 1",
		[]float64{0, 0, 0, 10, 20, 10, 20, 00},
		Vector{1, 0, 1},
		0.0,
		20.0,
	},
	{
		"box 2",
		[]float64{0, 0, 0, 10, 20, 10, 20, 00},
		Vector{0, 1, 1},
		0.0,
		10.0,
	},
	{
		"box 3",
		[]float64{0, 0, 0, 10, 20, 10, 20, 00},
		Vector{1, 1, 1},
		0.0,
		21.213,
	},
	{
		"ship test",
		[]float64{687.6316996852223, 312.63003058710115, 654.1102372351858, 339.1123181622106, 662.1094593091648, 331.3012717399208, 656.0111626926846, 312.25367928372236, 644.9627923104654, 310.540929477913},
		Vector{-0.6199034636337281 - 0.7846780841688566, 1},
		-378.846,
		-336.179,
	},
}

func Test_flattenPointsOn(t *testing.T) {
	for i, tt := range FlattenPointsTests {
		t.Run(fmt.Sprintf("%v_%s", i, tt.description), func(t *testing.T) {
			s := NewPolygon(tt.s...)
			min, max := s.flattenPointsOn(tt.n.Normalize())
			assert.InDelta(t, tt.min, min, 0.001)
			assert.InDelta(t, tt.max, max, 0.001)
		})
	}
}
