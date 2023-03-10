
package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "timer for the quiz")
	shuffle := flag.Bool("shuffle", false, "flag to know if to shuffle questions")

	flag.Parse()
	

	file, err := os.Open(*csvFilename)
	if err != nil {
		log.Fatalf("could not open the csv file %s\n",*csvFilename)
	}
	
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal("Failed to parse the provided csv file")
	}

	problems := parseLines(lines)
	// takes a flag of boolean to shuffle the slice
	// to shuffle the slice the command line would shuffle=true
	if *shuffle {
		problems = shuffleSlice(problems)
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	// problemloop: // this is a loop label
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.question)

		answerChannel := make(chan string)

		// this will like run in the background so fmt.Scanf doesn't block
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer) // this method blocks the app flow and doesn't allow the timer to stop the call even when it has expired
			answerChannel <- answer
		}()

		// for whicherver comes first this select statement will know what to do
		select {
		case <-timer.C:
			fmt.Printf("\nThe time limit has been reached.\nYou scored %d out of %d.", correct, len(problems))
			return // we coulssd use break but we want to break out of the loop totally so we use return
			// break statement will just end this iteration and move on to the next iteration in the loop
			// break problemloop // if we don't want to use a return statememnt we can break out of the loop using the
			// break statement on the label of the loop
		case answer := <-answerChannel:
			answer = strings.ToLower(answer)
			if answer == p.answer {
				correct++
			}
		}
	}

	// if you use the break statement this line will be executed
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(strings.ToLower(line[1])), // to ensure data is clean from the csv
		}
	}
	return ret
}

type problem struct {
	question string
	answer   string
}



func shuffleSlice(problems []problem) []problem {
	result := make([]problem, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(problems)) {
		result = append(result, problems[i])
	}
	return result
}