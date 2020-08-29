package bot

import (
	"fmt"
	"math"

	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/network"

	"github.com/grahamjenson/asteroids/game"
	"github.com/grahamjenson/asteroids/js/keys"
)

type cood struct {
	x, y, distance float64
}

// Contains methods used by both training and playing
func FindInputs(g *game.Game) []float64 {
	// start with BIAS
	inputs := []float64{1.0}

	for i := 0; i < 8; i++ {
		distance := g.Ship.WhiskerDistance(g.Ship.Whiskers[i], g.Asteroids)
		norm := (300 - distance) / 300
		inputs = append(inputs, exp(norm, 40))
	}

	return inputs
}

// takes value between 0 and 1 and bends to exp curve between 0 and 1
func exp(n, e float64) float64 {
	return (math.Pow(e, n) - 1.0) / (e - 1)
}

func GetOutputs(net *network.Network, g *game.Game) (map[int]bool, error) {
	inputs := FindInputs(g) // 8 inputs + bias

	outputs := map[int]bool{}

	net.LoadSensors(inputs)

	/*-- activate the network based on the input --*/
	if res, err := net.Activate(); !res || err != nil {
		//If it loops, exit returning only fitness of 1 step
		neat.DebugLog(fmt.Sprintf("Failed to activate Network, reason: %s", err))
		return outputs, err
	}

	for i, output := range net.ReadOutputs() {
		var key int

		switch i {
		case 0:
			key = keys.KEY_UP
		case 1:
			key = keys.KEY_SPACE
		case 2:
			key = keys.KEY_LEFT
		case 3:
			key = keys.KEY_RIGHT
		}

		// if output is activated we roll
		if output > 0.5 {
			outputs[key] = true
		}
	}

	return outputs, nil
}
