package main

import (
	"math"

	"github.com/grahamjenson/asteroids/canvas"
	"github.com/grahamjenson/asteroids/vector2d"
)

// https://codepen.io/anthonydugois/full/mewdyZ

///
// SHIP
///

type Ship struct {
	template   *vector2d.Polygon
	projection *vector2d.Polygon

	bulletTemplate   *vector2d.Polygon
	bulletProjection *vector2d.Polygon

	x, y                           float64
	velocityX, velocityY, rotation float64

	width, height int

	bulletX, bulletY, bulletVelocityX, bulletVelocityY, bulletDistance float64
}

func (s *Ship) Init(width, height int) {
	s.width = width
	s.height = height

	s.x = float64(width / 2)
	s.y = float64(height / 2)

	s.template = vector2d.NewPolygon(
		0, 2,
		-1.5, -2,
		-1, -1,
		1, -1,
		1.5, -2,
	)

	s.template.Scale(10, 10)

	s.bulletTemplate = vector2d.NewPolygon(
		0, 2,
		2, 0,
		0, -2,
	)

	s.bulletDistance = -1 // No bullet
}

func (s *Ship) Update(dt float64, pressedButtons map[int]bool) {

	if pressedButtons[KEY_LEFT] {
		s.rotation += math.Pi * dt * 2
	}
	if pressedButtons[KEY_RIGHT] {
		s.rotation += -math.Pi * dt * 2
	}
	if pressedButtons[KEY_UP] {
		s.velocityX += dt * 60 * math.Sin(s.rotation)
		s.velocityY += dt * 60 * math.Cos(s.rotation)

	}

	// Move will add to the velocity based on the direction the

	// we add friction on the velocity
	s.velocityX *= dt * (0.90 * 60)
	s.velocityY *= dt * (0.90 * 60)

	s.x += s.velocityX
	s.y += s.velocityY

	// BOUNDS
	// We limit rotation to tau
	s.rotation = math.Mod(s.rotation, (2 * math.Pi))

	s.x, s.y = WrapXY(s.x, s.y, s.width, s.height)

	// Bullet

	if pressedButtons[KEY_SPACE] {
		// Init fireing the bullet
		if s.bulletDistance == -1 {
			s.bulletDistance = 0
			s.bulletVelocityX = s.velocityX + (dt * 1000 * math.Sin(s.rotation))
			s.bulletVelocityY = s.velocityY + (dt * 1000 * math.Cos(s.rotation))

			// Make it align with the front of the ship
			s.bulletX = s.x + (20 * math.Sin(s.rotation))
			s.bulletY = s.y + (20 * math.Cos(s.rotation))
		}
	}

	if s.bulletDistance != -1 {
		prevX := s.bulletX
		prevY := s.bulletY
		s.bulletX += s.bulletVelocityX
		s.bulletY += s.bulletVelocityY

		// Must calculate distance before wrapping
		s.bulletDistance += Pythag(prevX-s.bulletX, prevY-s.bulletY)

		s.bulletX, s.bulletY = WrapXY(s.bulletX, s.bulletY, s.width, s.height)

		// Handle bullet

		if s.bulletDistance > 300 {
			s.bulletDistance = -1
		}
	}

	// update projections
	s.projection = s.template.Clone().Translate(s.x, s.y).Rotate(s.rotation)
	if s.bulletDistance != -1 {
		s.bulletProjection = s.bulletTemplate.Clone().Translate(s.bulletX, s.bulletY)
	}
}

func (s *Ship) Render(ctx *canvas.Context2D) {
	// Leave custom calculations for collision
	RenderPolygon(ctx, s.projection)

	if s.bulletDistance != -1 {
		RenderPolygon(ctx, s.bulletProjection)
	}

}
