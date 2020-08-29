package render

import (
	"fmt"
	"math"

	"github.com/grahamjenson/asteroids/game"
	"github.com/grahamjenson/asteroids/js/canvas"
	"github.com/grahamjenson/asteroids/vector2d"
)

func RenderGame(ctx *canvas.Context2D, g *game.Game) {
	switch g.State {
	case "menu":
		RenderMenu(ctx, g)
		RenderAlways(ctx, g)
	case "game":
		RenderShip(ctx, g.Ship, g.Asteroids)
		RenderAlways(ctx, g)
	}

}

func RenderAlways(ctx *canvas.Context2D, g *game.Game) {

	for _, a := range g.Asteroids {

		if a.Projection != nil {
			RenderPolygon(ctx, a.Projection)
		}

	}

	ctx.Save()
	ctx.SetTextAlign("center")
	ctx.SetFont("bold 20px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
	ctx.FillText(fmt.Sprintf("Score: %v", g.Score), g.Width/2, 40, g.Width)
	ctx.Restore()

}

func RenderMenu(ctx *canvas.Context2D, g *game.Game) {
	ctx.Save()
	ctx.SetTextAlign("center")
	ctx.SetFont("bold 40px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
	ctx.FillText(g.MenuText, g.Width/2, g.Height/2, g.Width)
	ctx.Restore()
}

func RenderShip(ctx *canvas.Context2D, s *game.Ship, asteroids []*game.Asteroid) {
	// Leave custom calculations for collision
	RenderPolygon(ctx, s.Projection)

	for i := 0; i < 8; i++ {
		whisker := s.Whiskers[i]
		x, y := whisker.Centroid()
		ctx.Save()
		ctx.SetLineDash([]interface{}{5, 15})
		ctx.SetStrokeStyle("rgba(0, 0, 0, .5)")
		RenderPolygon(ctx, whisker)

		distance := s.WhiskerDistance(whisker, asteroids)

		ctx.SetTextAlign("center")
		ctx.SetFont("bold 12px Courier New")
		ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
		ctx.FillText(
			fmt.Sprintf("%0.f", distance),
			int(x),
			int(y),
			1000,
		)

		ctx.Restore()
	}

	if s.IsBullet() {
		RenderPolygon(ctx, s.BulletProjection)
	}

}

func RenderScorePlayer(ctx *canvas.Context2D, g *game.Game, humanScore, botScore []int, human bool, botNumber int, timeLeft float64) {
	ctx.Save()

	ctx.SetTextAlign("center")
	ctx.SetFont("bold 20px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 1)")
	ctx.FillText("AI Scores", (g.Width/2)-g.Width/4, 100, g.Width/4)

	for i, s := range botScore {
		ctx.FillText(fmt.Sprintf("%v", s), (g.Width/2)-g.Width/4, 130+(i*30), g.Width/4)
	}

	ctx.FillText("Human Scores", (g.Width/2)+g.Width/4, 100, g.Width/4)
	for i, s := range humanScore {
		ctx.FillText(fmt.Sprintf("%v", s), (g.Width/2)+g.Width/4, 130+(i*30), g.Width/4)
	}

	tl := fmt.Sprintf("%0.0f", timeLeft)
	ctx.SetFont("bold 80px Courier New")
	ctx.SetStrokeStyle("rgba(125, 0, 0, 1)")
	if human {
		ctx.FillText("HUMAN "+tl+"s", g.Width/2, g.Height-40, g.Width)
	} else {
		ctx.FillText(fmt.Sprintf("SPECIES %v "+tl+"s", botNumber), g.Width/2, g.Height-40, g.Width)
	}

	ctx.Restore()
}

// Utils

func RenderPolygon(ctx *canvas.Context2D, s *vector2d.Polygon) {
	ctx.Save()

	ctx.BeginPath()
	first := true
	var firstPoint vector2d.Vector

	for _, v := range s.Matrix {
		if first {
			ctx.MoveTo(v[0], v[1])
			first = false
			firstPoint = v
		} else {
			ctx.LineTo(v[0], v[1])
		}
	}

	ctx.LineTo(firstPoint[0], firstPoint[1])
	ctx.Stroke()
	ctx.Restore()
}

func Pythag(x, y float64) float64 {
	return math.Sqrt((x * x) + (y * y))
}
