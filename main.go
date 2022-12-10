package main

import (
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/go-co-op/gocron"
	"github.com/mroth/weightedrand"
)

var every string
var promptFile string
var reloadFile bool

func init() {
	flag.StringVar(&every, "every", "42m", "Every specifies the duration and supports formats like 5s, 42m, 1h ")
	flag.StringVar(&promptFile, "f", "./prompts.txt", "Which file to load prompts from")
	flag.BoolVar(&reloadFile, "reload", false, "Whether to reload prompts file on each run...")
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // always seed random!
	flag.Parse()

	preloadedChoices := loadChoices(promptFile)
	s := gocron.NewScheduler(time.UTC)

	s.Every(every).Do(func() {
		var chooser *weightedrand.Chooser
		if reloadFile {
			choices := loadChoices(promptFile)
			chooser, _ = weightedrand.NewChooser(choices...)
		} else {
			chooser, _ = weightedrand.NewChooser(preloadedChoices...)
		}

		result := chooser.Pick()
		err := beeep.Notify("Promptato", result.(string), "question.png")
		if err != nil {
			log.Fatalf("failed to notify user, got %v", err)
		}
	})

	// starts the scheduler and blocks current execution path
	s.StartBlocking()
}

func loadChoices(sourceFile string) []weightedrand.Choice {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Fatalf("failed to load choices file, got %v", err)
		return nil
	}

	lines := strings.Split(string(data), "\n")
	choices := make([]weightedrand.Choice, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(strings.TrimPrefix(line, "- "))
		choices = append(choices, weightedrand.NewChoice(trimmed, 1))
	}

	return choices
}
