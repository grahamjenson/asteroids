package game

import (
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/grahamjenson/asteroids/js/keys"
)

type Game struct {
	Asteroids []*Asteroid
	Ship      *Ship

	State     string
	nextState string

	Score int

	Height, Width int
	MenuText      string
	Dead          bool
	Win           bool
	Seed          int64
}

func (g *Game) Init(width, height int) {
	g.Height = height
	g.Width = width

	g.MenuText = ""

	g.InitGame()

	g.State = "menu"
	g.Score = 0
	g.Seed = 0
}

func (g *Game) InitGame() {
	// 10 is hard
	// 6 medium
	// 4 is easy
	rand.Seed(g.Seed)
	defer rand.Seed(time.Now().UnixNano())

	g.Asteroids = []*Asteroid{&Asteroid{}, &Asteroid{}, &Asteroid{}, &Asteroid{}, &Asteroid{}}
	for _, a := range g.Asteroids {
		a.Init(g.Width, g.Height, 1)
	}

	g.Ship = &Ship{}
	g.Ship.Init(g.Width, g.Height)
	g.Dead = false
	g.Win = false
}

func (g *Game) Update(now float64, pressedKeys map[int]bool) {
	g.StateTransition()

	switch g.State {
	case "menu":
		g.UpdateAlways(now, pressedKeys)
		g.UpdateMenu(now, pressedKeys)
	case "game":
		g.UpdateAlways(now, pressedKeys)
		g.UpdateGame(now, pressedKeys)
	}

}

func (g *Game) StateTransition() {
	if g.nextState == "" {
		return
	}

	if g.State == "menu" && g.nextState == "game" {
		g.InitGame()
		g.Score = 0
		// Start State
	} else if g.State == "game" && g.nextState == "menu" {

	}

	g.State = g.nextState
	g.nextState = ""
}

func (g *Game) UpdateAlways(now float64, pressedKeys map[int]bool) {
	for _, a := range g.Asteroids {
		a.Update(now, pressedKeys)
	}
}

func (g *Game) UpdateMenu(now float64, pressedKeys map[int]bool) {
	if pressedKeys[keys.KEY_ENTER] || pressedKeys[keys.KEY_RETURN] {
		g.nextState = "game"
	}
}

func (g *Game) UpdateGame(now float64, pressedKeys map[int]bool) {
	g.Ship.Update(now, pressedKeys)

	g.Collisions()
}

func (g *Game) EndGame() {
	g.Dead = true
	g.nextState = "menu"
}

func (g *Game) Collisions() {
	s := g.Ship

	// Win condition
	if len(g.Asteroids) == 0 {
		g.Win = true
		g.nextState = "menu"
	}

	// does the ship hit an asteroid
	for _, a := range g.Asteroids {
		hc := s.Projection.HitCheck(a.Projection)
		if hc.Collision {
			g.EndGame()
		}
	}

	// does a bullet hit an asteroid

	if s.IsBullet() {
		var hitAsteroid *Asteroid
		for _, a := range g.Asteroids {
			bhc := a.Projection.HitCheck(s.BulletProjection)
			if bhc.Collision {
				g.Score += 1
				// Stop bullet
				hitAsteroid = a
				break
			}
		}

		// hit
		if hitAsteroid != nil {
			s.bulletDistance = -1 // stop bullet
			newAsteroids := []*Asteroid{}
			for _, a := range g.Asteroids {
				if a != hitAsteroid {
					// Remove the hit asteroid
					newAsteroids = append(newAsteroids, a)
				}
			}

			scale := hitAsteroid.scale * 0.7

			if scale > 0.3 {
				a1 := &Asteroid{}
				a2 := &Asteroid{}
				// Limit the number
				a1.Init(g.Width, g.Height, scale)
				a2.Init(g.Width, g.Height, scale)

				rot1 := -math.Pi / 8
				rot2 := +math.Pi / 8

				a1.x = hitAsteroid.x - (40 * hitAsteroid.scale)
				a1.y = hitAsteroid.y + (40 * hitAsteroid.scale)
				a1.velocityX = 1.2 * (hitAsteroid.velocityX*math.Cos(rot1) - hitAsteroid.velocityY*math.Sin(rot1))
				a1.velocityY = 1.2 * (hitAsteroid.velocityX*math.Sin(rot1) + hitAsteroid.velocityY*math.Cos(rot1))

				a2.x = hitAsteroid.x - (40 * hitAsteroid.scale)
				a2.y = hitAsteroid.y + (40 * hitAsteroid.scale)
				a2.velocityX = 1.2 * (hitAsteroid.velocityX*math.Cos(rot2) - hitAsteroid.velocityY*math.Sin(rot2))
				a2.velocityY = 1.2 * (hitAsteroid.velocityX*math.Sin(rot2) + hitAsteroid.velocityY*math.Cos(rot2))

				// add the new asteroids
				newAsteroids = append(newAsteroids, a1)
				newAsteroids = append(newAsteroids, a2)
				sort.Slice(newAsteroids, func(i, j int) bool {
					return newAsteroids[i].scale > newAsteroids[j].scale
				})
			}
			g.Asteroids = newAsteroids
		}
	}

}

////
// Utils
////

// useful https://jlongster.com/Making-Sprite-based-Games-with-Canvas
func WrapXY(x, y float64, width, height int) (float64, float64) {
	nx := math.Mod(x, float64(width))
	ny := math.Mod(y, float64(height))

	if nx < 0 {
		nx = float64(width)
	}

	if ny < 0 {
		ny = float64(height)
	}
	return nx, ny
}

func Pythag(x, y float64) float64 {
	return math.Sqrt((x * x) + (y * y))
}
