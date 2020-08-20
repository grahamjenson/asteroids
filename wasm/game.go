package main

import (
	"fmt"
	"sort"

	"github.com/grahamjenson/asteroids/canvas"
)

type Game struct {
	asteroids   []*Asteroid
	ship        *Ship
	diagnostics *Diagnostics

	state     string
	nextState string

	score int

	height, width int
	menuText      string
}

func (g *Game) Init(width, height int) {
	g.height = height
	g.width = width

	g.menuText = "*** PRESS ENTER TO START ***"

	g.InitGame()
	g.diagnostics = &Diagnostics{}

	g.state = "menu"
	g.score = 0

	// INIT

	g.diagnostics.Init(width, height)
}

func (g *Game) InitGame() {
	g.asteroids = []*Asteroid{&Asteroid{}, &Asteroid{}, &Asteroid{}, &Asteroid{}, &Asteroid{}}
	for _, a := range g.asteroids {
		a.Init(g.width, g.height, 1)
	}

	g.ship = &Ship{}
	g.ship.Init(g.width, g.height)
}

func (g *Game) Render(ctx *canvas.Context2D) {
	switch g.state {
	case "menu":
		g.RenderMenu(ctx)
		g.RenderAlways(ctx)
	case "game":
		g.RenderGame(ctx)
		g.RenderAlways(ctx)
	}

}

func (g *Game) Update(now float64, pressedKeys map[int]bool) {
	g.stateTransition()

	switch g.state {
	case "menu":
		g.UpdateAlways(now, pressedKeys)
		g.UpdateMenu(now, pressedKeys)
	case "game":
		g.UpdateAlways(now, pressedKeys)
		g.UpdateGame(now, pressedKeys)
	}

}

func (g *Game) stateTransition() {
	if g.nextState == "" {
		return
	}

	if g.state == "menu" && g.nextState == "game" {
		g.InitGame()
		g.score = 0
		// Start State
	} else if g.state == "game" && g.nextState == "menu" {
		g.menuText = "*** GAME OVER! PRESS ENTER TO START ***"
	}

	g.state = g.nextState
	g.nextState = ""
}

func (g *Game) RenderAlways(ctx *canvas.Context2D) {
	for _, a := range g.asteroids {
		a.Render(ctx)
	}
	g.diagnostics.Render(ctx)

	ctx.Save()
	ctx.SetTextAlign("center")
	ctx.SetFont("bold 20px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
	ctx.FillText(fmt.Sprintf("Score: %v", g.score), g.width/2, 40, g.width)
	ctx.Restore()

}

func (g *Game) RenderMenu(ctx *canvas.Context2D) {
	ctx.Save()
	ctx.SetTextAlign("center")
	ctx.SetFont("bold 40px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
	ctx.FillText(g.menuText, g.width/2, g.height/2, g.width)
	ctx.Restore()
}

func (g *Game) RenderGame(ctx *canvas.Context2D) {
	g.ship.Render(ctx)
}

func (g *Game) UpdateAlways(now float64, pressedKeys map[int]bool) {
	for _, a := range g.asteroids {
		a.Update(now, pressedKeys)
	}

	g.diagnostics.Update(now, pressedKeys)
}

func (g *Game) UpdateMenu(now float64, pressedKeys map[int]bool) {
	if pressedKeys[KEY_ENTER] || pressedKeys[KEY_RETURN] {
		g.nextState = "game"
	}
}

func (g *Game) UpdateGame(now float64, pressedKeys map[int]bool) {
	g.ship.Update(now, pressedKeys)

	g.Collisions()
}

func (g *Game) Collisions() {
	s := g.ship

	// Win condition
	if len(g.asteroids) == 0 {
		g.nextState = "menu"
	}

	// does the ship hit an asteroid
	for _, a := range g.asteroids {
		hc := s.projection.HitCheck(a.projection)
		if hc.Collision {
			g.nextState = "menu"
		}
	}

	for i, a1 := range g.asteroids {
		for j := i + 1; j < len(g.asteroids); j++ {
			a2 := g.asteroids[j]

			hc := a1.projection.HitCheck(a2.projection)
			if hc.Collision {
				// The first asteroid has preference
				a2.x -= hc.OverlapV[0]
				a2.y -= hc.OverlapV[1]
			}
		}
	}

	// does a bullet hit an asteroid

	if s.bulletDistance != -1 {
		var hitAsteroid *Asteroid
		for _, a := range g.asteroids {
			bhc := a.projection.HitCheck(s.bulletProjection)
			if bhc.Collision {
				g.score += 10
				// Stop bullet
				hitAsteroid = a
			}
		}

		// hit
		if hitAsteroid != nil {
			s.bulletDistance = -1 // stop bullet
			newAsteroids := []*Asteroid{}
			for _, a := range g.asteroids {
				if a != hitAsteroid {
					// Remove the hit asteroid
					newAsteroids = append(newAsteroids, a)
				}
			}

			scale := hitAsteroid.scale * 0.7
			if scale > 0.05 {
				a1 := &Asteroid{}
				a2 := &Asteroid{}
				// Limit the number
				a1.Init(g.width, g.height, scale)
				a2.Init(g.width, g.height, scale)

				a1.x = hitAsteroid.x - (20 * hitAsteroid.scale)
				a1.y = hitAsteroid.y + (20 * hitAsteroid.scale)
				a1.velocityX = hitAsteroid.velocityX
				a1.velocityY = hitAsteroid.velocityY * 1.1

				a2.x = hitAsteroid.x - (20 * hitAsteroid.scale)
				a2.y = hitAsteroid.y + (20 * hitAsteroid.scale)
				a2.velocityX = hitAsteroid.velocityY * 1.1

				// add the new asteroids
				newAsteroids = append(newAsteroids, a1)
				newAsteroids = append(newAsteroids, a2)
				sort.Slice(newAsteroids, func(i, j int) bool {
					return newAsteroids[i].scale > newAsteroids[j].scale
				})
			}
			g.asteroids = newAsteroids
		}
	}
}
