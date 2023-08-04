package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

func main() {
	timerDuration := flag.Int("timer", 30, "the time limit for the quiz in seconds")
	csvFilename := flag.String("csv", "problem.csv", "a csv file in the format of 'question,answer'")

	flag.Parse()

	timer := time.NewTimer(time.Duration(*timerDuration) * time.Second)

	f, err := os.Open(*csvFilename)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)

	var score int

	for {
		line, err := csvReader.Read()
		var question string
		question = line[0]

		fmt.Printf(question)

		ansCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanln(&answer)
			ansCh <- answer
		}()

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		select {
		case <-timer.C:
			fmt.Printf("Time's up! for %d\n", *timerDuration)
			fmt.Println("Your score is", score)
			return
		case answer := <-ansCh:
			if answer == line[1] {
				score++
			}
		}
	}
}
