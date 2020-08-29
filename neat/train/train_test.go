package train

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/yaricom/goNEAT/experiments"
	"github.com/yaricom/goNEAT/neat"
	"github.com/yaricom/goNEAT/neat/genetics"
)

func Test_TrainAsteroids(t *testing.T) {

	rand.Seed(time.Now().Unix())

	out_dir_path := "/tmp/neat/out/asteroids_test"

	context := neat.LoadContext(strings.NewReader(CONFIG_FILE))

	fmt.Println(GENOME_INIT)
	start_genome, err := genetics.ReadGenome(strings.NewReader(GENOME_INIT), 1)
	if err != nil {
		t.Error("Failed to read start genome")
		return
	}

	neat.LogLevel = neat.LogLevelDebug

	if _, err := os.Stat(out_dir_path); err == nil {

		os.RemoveAll(out_dir_path)
	}

	err = os.MkdirAll(out_dir_path, os.ModePerm)
	if err != nil {
		t.Errorf("Failed to create output directory, reason: %s", err)
		return
	}

	context.NumRuns = 1
	experiment := experiments.Experiment{
		Id:     0,
		Trials: make(experiments.Trials, context.NumRuns),
	}
	err = experiment.Execute(context, start_genome, AsteroidGenerationEvaluator{
		OutputPath:         out_dir_path,
		PlayTimeInSeconds:  200,
		FrameRatePerSecond: 60,
	})
	if err != nil {
		t.Error("Failed to perform Asteroids experiment:", err)
		return
	}
}
