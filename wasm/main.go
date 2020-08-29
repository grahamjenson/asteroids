package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/grahamjenson/asteroids/game"
	"github.com/grahamjenson/asteroids/js/canvas"
	"github.com/grahamjenson/asteroids/neat/bot"
	"github.com/grahamjenson/asteroids/neat/players"
	"github.com/grahamjenson/asteroids/render"
	"github.com/yaricom/goNEAT/neat/genetics"
	"github.com/yaricom/goNEAT/neat/network"
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

	humanButtons := map[int]bool{}
	human := false

	humanscore := []int{}
	botscore := []int{}

	bots := []*network.Network{
		getGenome(players.PLAYER_1),
		getGenome(players.PLAYER_2),
		getGenome(players.PLAYER_3),
		getGenome(players.PLAYER_4),
		getGenome(players.PLAYER_5),
	}
	botN := 0

	TIME_TO_PLAY := 20.0
	timeLeft := TIME_TO_PLAY

	window.Call(
		"addEventListener",
		"keyup",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			e.Call("preventDefault")
			humanButtons[e.Get("keyCode").Int()] = false
			return nil
		}))

	window.Call(
		"addEventListener",
		"keydown",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			e := args[0]
			e.Call("preventDefault")
			humanButtons[e.Get("keyCode").Int()] = true
			return nil
		}))

	game := &game.Game{}
	game.Init(width, height)

	diagnostics := &Diagnostics{}
	diagnostics.Init(width, height)

	var gameLoop js.Func
	prevNow := 0.0
	gameLoop = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		go func() {
			now := args[0].Float()
			dt := (now - prevNow) / 1000
			prevNow = now

			// kill if timeout
			if timeLeft < 0 {
				game.EndGame()
				timeLeft = TIME_TO_PLAY
			}

			if game.State == "menu" {
				game.Update(dt, humanButtons)
				timeLeft = TIME_TO_PLAY
			} else if human {
				game.Update(dt, humanButtons)
				timeLeft -= dt

			} else {
				botButtons, err := bot.GetOutputs(bots[botN], game)
				if err != nil {
					fmt.Println("bot outputs error", err)
					return
				}
				game.Update(dt, botButtons)
				timeLeft -= dt

			}

			diagnostics.Update(dt)
			if game.Dead || game.Win {
				for _, b := range bots {
					b.Flush()
				}

				if human {
					humanscore = append(humanscore, game.Score)
				} else {
					botscore = append(botscore, game.Score)
				}
				human = !human
				if human {
					game.Seed = int64(len(humanscore)) + 10
				} else {
					botN = (botN + 1) % 5
					game.Seed = int64(len(botscore)) + 10
				}

				game.Dead = false
				game.Win = false
			}

			// Clear Screen
			ctx.ClearRect(0, 0, width, height) // clear canvas

			// Render
			render.RenderGame(ctx, game)
			render.RenderScorePlayer(ctx, game, humanscore, botscore, human, botN, timeLeft)
			diagnostics.Render(ctx)

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

func getGenome(genomeStr string) *network.Network {
	genome, _ := genetics.ReadGenome(strings.NewReader(genomeStr), 1)

	net, _ := genome.Genesis(1)

	return net
}
