package game

import (
	"math"

	"github.com/grahamjenson/asteroids/js/keys"
	"github.com/grahamjenson/asteroids/vector2d"
)

// https://codepen.io/anthonydugois/full/mewdyZ

///
// SHIP
///
const whiskerDistance float64 = 300.0

type Ship struct {
	template   *vector2d.Polygon
	Projection *vector2d.Polygon

	whiskersTemplate []*vector2d.Polygon
	Whiskers         []*vector2d.Polygon

	bulletTemplate   *vector2d.Polygon
	BulletProjection *vector2d.Polygon

	x, y                           float64
	velocityX, velocityY, Rotation float64

	width, height int

	bulletX, bulletY, bulletVelocityX, bulletVelocityY, bulletDistance float64

	Human bool
}

func (s *Ship) Init(width, height int) {
	s.width = width
	s.height = height

	s.x = float64(width / 2)
	s.y = float64(height / 2)
	s.velocityX = 0
	s.velocityY = 0
	s.Rotation = 0

	s.template = vector2d.NewPolygon(
		0, 2,
		-1.5, -2,
		-1, -1,
		1, -1,
		1.5, -2,
	)

	s.template.Scale(10, 10)
	s.Projection = s.template

	s.bulletTemplate = vector2d.NewPolygon(
		0, 2,
		2, 0,
		0, -2,
	)

	s.bulletDistance = -1 // No bullet

	// whikers
	s.whiskersTemplate = make([]*vector2d.Polygon, 8)
	s.Whiskers = make([]*vector2d.Polygon, 8)

	s.whiskersTemplate[0] = vector2d.NewPolygon(0, 0, 0, whiskerDistance)
	s.whiskersTemplate[1] = vector2d.NewPolygon(0, 0, 0, -whiskerDistance)
	s.whiskersTemplate[2] = vector2d.NewPolygon(0, 0, whiskerDistance, 0)
	s.whiskersTemplate[3] = vector2d.NewPolygon(0, 0, -whiskerDistance, 0)
	s.whiskersTemplate[4] = vector2d.NewPolygon(0, 0, whiskerDistance, whiskerDistance)
	s.whiskersTemplate[5] = vector2d.NewPolygon(0, 0, -whiskerDistance, -whiskerDistance)
	s.whiskersTemplate[6] = vector2d.NewPolygon(0, 0, whiskerDistance, -whiskerDistance)
	s.whiskersTemplate[7] = vector2d.NewPolygon(0, 0, -whiskerDistance, whiskerDistance)

	for i := 0; i < 8; i++ {
		s.Whiskers[i] = s.whiskersTemplate[i].Clone()
	}
}

func (s *Ship) IsBullet() bool {
	return s.bulletDistance != -1
}
func (s *Ship) Update(dt float64, pressedButtons map[int]bool) {

	controlCoef := 1.0
	if !s.Human {
		controlCoef = 2.0
	}
	if pressedButtons[keys.KEY_LEFT] {
		s.Rotation += math.Pi * dt * 2 * controlCoef
	}
	if pressedButtons[keys.KEY_RIGHT] {
		s.Rotation += -math.Pi * dt * 2 * controlCoef
	}
	if pressedButtons[keys.KEY_UP] {
		s.velocityX += dt * 60 * math.Sin(s.Rotation) * controlCoef
		s.velocityY += dt * 60 * math.Cos(s.Rotation) * controlCoef
	}

	// Move will add to the velocity based on the direction the

	// we add friction on the velocity
	s.velocityX *= dt * (0.90 * 60)
	s.velocityY *= dt * (0.90 * 60)

	s.x += s.velocityX
	s.y += s.velocityY

	// BOUNDS
	// We limit Rotation to tau
	s.Rotation = math.Mod(s.Rotation, (2 * math.Pi))

	s.x, s.y = WrapXY(s.x, s.y, s.width, s.height)

	// Bullet

	if pressedButtons[keys.KEY_SPACE] {
		// Init fireing the bullet
		if !s.IsBullet() {
			s.bulletDistance = 0
			s.bulletVelocityX = s.velocityX + (dt * 1000 * math.Sin(s.Rotation))
			s.bulletVelocityY = s.velocityY + (dt * 1000 * math.Cos(s.Rotation))

			// Make it align with the front of the ship
			s.bulletX = s.x + (20 * math.Sin(s.Rotation))
			s.bulletY = s.y + (20 * math.Cos(s.Rotation))
		}
	}

	if s.IsBullet() {
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

	// update Projections
	s.Projection = s.template.Clone().Translate(s.x, s.y).Rotate(s.Rotation)

	for i := 0; i < 8; i++ {
		s.Whiskers[i] = s.whiskersTemplate[i].Clone().Translate(s.x, s.y).RotateAbout(s.Rotation, s.x, s.y)
	}

	if s.IsBullet() {
		s.BulletProjection = s.bulletTemplate.Clone().Translate(s.bulletX, s.bulletY)
	}
}

// return the distances to the closet asteroid from a whisker
func (s *Ship) WhiskerDistance(whisker *vector2d.Polygon, asteroids []*Asteroid) float64 {
	minDistance := whiskerDistance
	projections := []*vector2d.Polygon{}
	for _, a := range asteroids {
		projections = append(projections, a.Projection)
	}

	// If the endpoint of the whisker pokes into another screen
	// 1 2 3
	// 4 S 5
	// 6 7 8
	// We copy the asteroids projections into that screen to check collisions there as well
	wx := whisker.Matrix[1][0]
	wy := whisker.Matrix[1][1]

	w := float64(s.width)
	h := float64(s.height)
	for _, a := range asteroids {
		switch {
		case wx < 0 && wy < 0:
			projections = append(projections, a.Projection.Clone().Translate(-w, -h)) // 1
		case wx > 0 && wx < w && wy < 0:
			projections = append(projections, a.Projection.Clone().Translate(0, -h)) // 2
		case wx > w && wy < 0:
			projections = append(projections, a.Projection.Clone().Translate(w, -h)) // 3
		case wx < 0 && wy > 0 && wy < h:
			projections = append(projections, a.Projection.Clone().Translate(-w, 0)) // 4
		case wx > w && wy > 0 && wy < h:
			projections = append(projections, a.Projection.Clone().Translate(w, 0)) // 5
		case wx < 0 && wy > h:
			projections = append(projections, a.Projection.Clone().Translate(-w, h)) // 6
		case wx > 0 && wx < w && wy > h:
			projections = append(projections, a.Projection.Clone().Translate(0, h)) // 7
		case wx > w && wy > h:
			projections = append(projections, a.Projection.Clone().Translate(w, h)) // 8
		}
	}

	for _, projection := range projections {
		bhc := whisker.HitCheck(projection)
		if bhc.Collision {
			// TODO find closest intersection point, e.g. https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

			// find closest point to ship (this is not perfect)
			for _, p := range projection.Matrix {
				d := Pythag(p[0]-s.x, p[1]-s.y)
				minDistance = math.Min(minDistance, d)
			}
		}
	}
	return minDistance
}

// return the relative XY coordinates of an asteroid
// so an asteroid directly infront of the ship is on the Y axis
func (s *Ship) AsteroidXYDistance(a *Asteroid) (float64, float64, float64) {
	// Find centers of the projections
	sx, sy := s.Projection.Centroid()
	ax, ay := a.Projection.Centroid()

	x := sx - ax
	y := sy - ay
	distance := Pythag(x, y)

	rx := (x * math.Cos(s.Rotation)) - (y * math.Sin(s.Rotation))
	ry := (y * math.Cos(s.Rotation)) + (x * math.Sin(s.Rotation))

	return rx, ry, distance
}
