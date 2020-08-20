package main

import (
	"math"
	"syscall/js"

	"github.com/grahamjenson/asteroids/canvas"
	"github.com/grahamjenson/asteroids/vector2d"
)

func initScreen(window js.Value, w, h int) *canvas.Context2D {
	document := window.Get("document")
	canvasE := document.Call("createElement", "canvas")

	canvasE.Set("id", "canvas")
	canvasE.Set("width", w)
	canvasE.Set("height", h)
	document.Get("body").Call("appendChild", canvasE)
	document.Get("body").Set("style", "margin: 0px; overflow: hidden")

	return canvas.NewContext2D(canvasE)
}

func main() {
	width := 1280
	height := 720

	window := js.Global()
	ctx := initScreen(window, width, height)
	ctx.SetGlobalCompositeOperation("destination-over")

	pressedButtons := map[int]bool{}

	window.Call(
		"addEventListener",
		"keyup",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			e.Call("preventDefault")
			pressedButtons[e.Get("keyCode").Int()] = false
			return nil
		}))

	window.Call(
		"addEventListener",
		"keydown",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			e.Call("preventDefault")
			pressedButtons[e.Get("keyCode").Int()] = true
			return nil
		}))

	game := Game{}
	game.Init(width, height)

	var gameLoop js.Func
	prevNow := 0.0
	gameLoop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			now := args[0].Float()
			dt := (now - prevNow) / 1000
			prevNow = now

			game.Update(dt, pressedButtons)

			// Clear Screen
			ctx.ClearRect(0, 0, width, height) // clear canvas
			// Render

			game.Render(ctx)

			// Interactions phase
			window.Call("requestAnimationFrame", gameLoop)
		}()
		return nil
	})

	// start the game render
	window.Call("requestAnimationFrame", gameLoop)

	// Never return
	done := make(chan struct{}, 0)
	<-done
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
