package main

import (
	"fmt"

	"github.com/grahamjenson/asteroids/js/canvas"
)

type Diagnostics struct {
	markCount     int
	tdiffSum      float64
	currentText   string
	width, height int
}

func (d *Diagnostics) Init(width, height int) {
	d.width = width
	d.height = height
}

func (d *Diagnostics) Update(dt float64) {
	d.tdiffSum += dt
	d.markCount++
}

func (d *Diagnostics) Render(ctx *canvas.Context2D) {
	if d.markCount >= 10 {
		// calculate an average of every 10 frames
		d.currentText = fmt.Sprintf("FPS: %.01f", float64(d.markCount)/d.tdiffSum)
		d.tdiffSum = 0
		d.markCount = 0
	}

	ctx.Save()
	ctx.SetFont("12px Courier New")
	ctx.SetStrokeStyle("rgba(0, 0, 0, 0.4)")
	ctx.StrokeText(d.currentText, 0, d.height-6, d.width)
	ctx.Restore()
}
