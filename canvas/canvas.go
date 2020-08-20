package canvas

import "syscall/js"

type Context2D struct {
	JS *js.Value
}

func NewContext2D(canvas js.Value) *Context2D {
	ctx := canvas.Call("getContext", "2d")
	return &Context2D{JS: &ctx}
}

////
// Sets
////
func (ctx *Context2D) SetGlobalCompositeOperation(co string) {
	ctx.JS.Set("globalCompositeOperation", co)
}

func (ctx *Context2D) SetFont(co string) {
	ctx.JS.Set("font", co)
}

func (ctx *Context2D) SetFillStyle(co string) {
	ctx.JS.Set("fillStyle", co)
}

func (ctx *Context2D) SetStrokeStyle(co string) {
	ctx.JS.Set("strokeStyle", co)
}

func (ctx *Context2D) SetTextAlign(co string) {
	ctx.JS.Set("textAlign", co)
}

////
// Calls
////

func (ctx *Context2D) Save() {
	ctx.JS.Call("save")
}

func (ctx *Context2D) Restore() {
	ctx.JS.Call("restore")
}

func (ctx *Context2D) BeginPath() {
	ctx.JS.Call("beginPath")
}

func (ctx *Context2D) ClosePath() {
	ctx.JS.Call("closePath")
}

func (ctx *Context2D) Stroke() {
	ctx.JS.Call("stroke")
}

func (ctx *Context2D) StrokeText(text string, x, y, maxWidth int) {
	ctx.JS.Call("strokeText", text, x, y, maxWidth)
}

func (ctx *Context2D) FillText(text string, x, y, maxWidth int) {
	ctx.JS.Call("fillText", text, x, y, maxWidth)
}

func (ctx *Context2D) ClearRect(x, y, width, height int) {
	ctx.JS.Call("clearRect", x, y, width, height)
}

func (ctx *Context2D) FillRect(x, y, width, height int) {
	ctx.JS.Call("fillRect", x, y, width, height)
}

// Transformations
func (ctx *Context2D) Rotate(angle float64) {
	ctx.JS.Call("rotate", angle)
}

func (ctx *Context2D) Translate(x, y float64) {
	ctx.JS.Call("translate", x, y)
}

func (ctx *Context2D) DrawImage(e js.Value, cords ...float64) {
	args := []interface{}{e}
	for _, c := range cords {
		args = append(args, c)
	}

	ctx.JS.Call("drawImage", args...)
}

/// Stroke

func (ctx *Context2D) MoveTo(x, y float64) {
	ctx.JS.Call("moveTo", x, y)
}

func (ctx *Context2D) LineTo(x, y float64) {
	ctx.JS.Call("lineTo", x, y)
}

func (ctx *Context2D) Arc(x, y int, radius, startAngle, endAngle float64, anticlockwise bool) {
	ctx.JS.Call("arc", x, y, radius, startAngle, endAngle, anticlockwise)
}
