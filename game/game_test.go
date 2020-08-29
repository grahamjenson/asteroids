package game

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/grahamjenson/asteroids/js/keys"
)

func Test_TrainEmptyAsteroids(t *testing.T) {
	totaltime := 0.0
	mintime := 1000.0
	maxtime := -1000.0
	iters := 1000
	for i := 0; i < iters; i++ {
		t := PlayConstantGame(map[int]bool{
			keys.KEY_SPACE: true,
			keys.KEY_LEFT:  true,
		}, i)
		maxtime = math.Max(maxtime, t)
		mintime = math.Min(mintime, t)
		totaltime += t
	}
	meantime := totaltime / float64(iters)
	fmt.Println("DO NOTHING", meantime, maxtime, mintime)
}

func PlayConstantGame(inputs map[int]bool, seed int) float64 {
	rand.Seed(int64(seed))
	game := &Game{}
	game.Init(1280, 720)
	game.State = "game" // force no menu

	game.Update(0.1, inputs)

	playtime := 60.0
	framerate := 10.0
	frames := int(playtime * framerate)
	frame := 0
	for frame = 0; frame < frames; frame++ {
		game.Update(1.0/framerate, inputs)

		if game.State != "game" {

			// Also add the number of seconds that he survived
			break
		}
	}
	timplayed := float64(frame) / framerate
	return timplayed
}
