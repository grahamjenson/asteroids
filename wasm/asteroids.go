package main

import (
	"math"
	"math/rand"

	"github.com/grahamjenson/asteroids/canvas"
	"github.com/grahamjenson/asteroids/vector2d"
)

///
// Asteroid
///

type Asteroid struct {
	template   *vector2d.Polygon
	projection *vector2d.Polygon

	x, y                                      float64
	velocityX, velocityY, rotationD, rotation float64

	width, height int
	scale         float64
}

var random *rand.Rand = rand.New(rand.NewSource(100))

func CoinToss() bool { return rand.Intn(2) == 0 }

func (s *Asteroid) Init(width, height int, scale float64) {
	s.width = width
	s.height = height
	s.scale = scale

	xOffset := random.Intn(400)
	yOffset := rand.Intn(400)

	if random.Intn(2) == 0 {
		s.x = float64(width - xOffset)
	} else {
		s.x = float64(xOffset)
	}

	if random.Intn(2) == 0 {
		s.y = float64(height - yOffset)
	} else {
		s.y = float64(yOffset)
	}

	s.velocityX = float64(random.Intn(400) - 200) // -200 to 200
	s.velocityY = float64(random.Intn(400) - 200)
	s.rotationD = random.Float64() * (math.Pi / 2)

	s.template = vector2d.NewPolygon(
		40, -40, 120, -20, 120, 40, 80, 80, 20, 80, -40, 100, -80, 60, -60, 00, -20, -40,
	)

	s.template.Scale(scale, scale)
}

func (s *Asteroid) Update(dt float64, pressedButtons map[int]bool) {

	s.x += s.velocityX * dt
	s.y += s.velocityY * dt
	s.rotation += s.rotationD * dt

	// BOUNDS
	// We limit rotation to tau
	s.rotation = math.Mod(s.rotation, (2 * math.Pi))

	s.x, s.y = WrapXY(s.x, s.y, s.width, s.height)

	// update projection
	s.projection = s.template.Clone().Translate(s.x, s.y).Rotate(s.rotation)
}

func (s *Asteroid) Render(ctx *canvas.Context2D) {
	if s.projection != nil {
		RenderPolygon(ctx, s.projection)
	}
}
