package train

import (
	"fmt"
	"math"
	"os"

	"github.com/grahamjenson/asteroids/game"
	"github.com/grahamjenson/asteroids/neat/bot"
	"github.com/yaricom/goNEAT/experiments"
	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/genetics"
	"github.com/yaricom/goNEAT/neat/network"
)

type AsteroidGenerationEvaluator struct {
	OutputPath string

	PlayTimeInSeconds  float64
	FrameRatePerSecond float64
}

func (ex AsteroidGenerationEvaluator) GenerationEvaluate(pop *genetics.Population, epoch *experiments.Generation, context *neat.NeatContext) (err error) {

	for _, org := range pop.Organisms {
		res := ex.orgEvaluate(org, epoch.Id)
		if res && (epoch.Best == nil || org.Fitness > epoch.Best.Fitness) {
			epoch.Solved = true
			epoch.WinnerNodes = len(org.Genotype.Nodes)
			epoch.WinnerGenes = org.Genotype.Extrons()
			epoch.WinnerEvals = context.PopSize*epoch.Id + org.Genotype.Id
			epoch.Best = org
		}
	}

	epoch.FillPopulationStatistics(pop)

	if epoch.Best != nil {
		best_org_path := fmt.Sprintf("%s/best_%04d", experiments.OutDirForTrial(ex.OutputPath, epoch.TrialId), epoch.Id)
		file, err := os.Create(best_org_path)
		if err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
		} else {
			org := epoch.Best
			fmt.Fprintf(file, "/* Organism #%d Fitness: %.3f Error: %.3f */\n",
				org.Genotype.Id, org.Fitness, org.Error)
			org.Genotype.Write(file)
		}
	}

	if epoch.Solved || epoch.Id%context.PrintEvery == 0 {

		pop_path := fmt.Sprintf("%s/gen_%04d", experiments.OutDirForTrial(ex.OutputPath, epoch.TrialId), epoch.Id)
		file, err := os.Create(pop_path)
		if err != nil {
			neat.ErrorLog(fmt.Sprintf("Failed to dump population, reason: %s\n", err))
		} else {
			pop.WriteBySpecies(file)
		}
	}

	return err
}

func (ex *AsteroidGenerationEvaluator) orgEvaluate(organism *genetics.Organism, gen int) bool {

	fit := ex.runGame(organism.Phenotype, gen)

	if fit >= 1.0 {
		organism.IsWinner = true
	}

	organism.Error = 1.0 - fit
	organism.Fitness = fit

	return organism.IsWinner
}

func (ex *AsteroidGenerationEvaluator) runGame(net *network.Network, gen int) float64 {

	game := &game.Game{}

	game.Init(1280, 720)
	game.State = "game"
	frames := int(ex.PlayTimeInSeconds * ex.FrameRatePerSecond)

	frame := 0
	died := false

	for frame = 0; frame < frames; frame++ {
		outputs, err := bot.GetOutputs(net, game)
		if err != nil {
			neat.ErrorLog(fmt.Sprintf("Error GetOutputs %s", err))
			return 0
		}

		game.Update(1.0/ex.FrameRatePerSecond, outputs)

		if game.State != "game" {
			died = game.Dead
			break
		}

		if game.Score == 0 && frame > 60 {
			return 0
		}
	}

	net.Flush()

	secondsSurvived := float64(frame) / ex.FrameRatePerSecond

	maxScore := 75.0 // 5 + 10 + 20 + 40
	maxSeconds := ex.PlayTimeInSeconds

	secondsPercent := secondsSurvived / maxSeconds

	fit := fitness(float64(game.Score), secondsSurvived, secondsPercent, maxScore, died)
	norm := fit

	fmt.Printf("%2.0f, %4.2f, %0.2f\n", float64(game.Score), secondsSurvived, norm)

	return norm
}

func fitness(score, seconds, secondsPercent, maxScore float64, died bool) float64 {
	if score != maxScore {
		return 0.0
	}
	return 1 - secondsPercent
}

func exp(n, e float64) float64 {
	return (math.Pow(e, n) - 1.0) / (e - 1)
}
