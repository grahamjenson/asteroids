package game

import (
	"math"
	"math/rand"

	"github.com/grahamjenson/asteroids/vector2d"
)

///
// Asteroid
///

type Asteroid struct {
	template   *vector2d.Polygon
	Projection *vector2d.Polygon

	x, y                                      float64
	velocityX, velocityY, rotationD, rotation float64

	width, height int
	scale         float64
}

func (s *Asteroid) XY() (float64, float64) {
	return s.x, s.y
}
func (s *Asteroid) Init(width, height int, scale float64) {
	s.width = width
	s.height = height
	s.scale = scale

	xOffset := rand.Intn(500)
	yOffset := rand.Intn(500)

	if rand.Intn(2) == 0 {
		s.x = float64(width - xOffset)
	} else {
		s.x = float64(xOffset)
	}

	if rand.Intn(2) == 0 {
		s.y = float64(height - yOffset)
	} else {
		s.y = float64(yOffset)
	}

	initSpeed := 150
	s.velocityX = float64(rand.Intn(initSpeed*2) - (initSpeed))
	s.velocityY = float64(rand.Intn(initSpeed*2) - (initSpeed))
	s.rotationD = rand.Float64() * (math.Pi / 2)

	s.template = vector2d.NewPolygon(
		40, -40, 120, -20, 120, 40, 80, 80, 20, 80, -40, 100, -80, 60, -60, 00, -20, -40,
	)

	s.template.Scale(scale, scale)
	s.Projection = s.template.Clone()
}

func (s *Asteroid) Update(dt float64, pressedButtons map[int]bool) {

	s.x += s.velocityX * dt
	s.y += s.velocityY * dt
	s.rotation += s.rotationD * dt

	// BOUNDS
	// We limit rotation to tau
	s.rotation = math.Mod(s.rotation, (2 * math.Pi))

	s.x, s.y = WrapXY(s.x, s.y, s.width, s.height)

	// update Projection
	s.Projection = s.template.Clone().Translate(s.x, s.y).Rotate(s.rotation)
}
