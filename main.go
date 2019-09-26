package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type question struct {
	text   string
	answer string
}

var (
	total   = 0
	correct = 0
)

func main() {
	filename := flag.String("f", "questions.csv", "provide a path to a question file (CSV format)")
	limit := flag.Duration("t", 30*time.Second, "a time limit in second for the whole set of question")
	flag.Parse()

	var enter string
	fmt.Println("Please press ENTER when you are ready...")
	fmt.Scanln(&enter)

	questions := loadQuestions(*filename)
	total = len(questions)
	done := make(chan int)

	go func() {
		for _, q := range questions {
			fmt.Printf("%s\n", q.text)
			fmt.Print("\r> ")
			var answer string
			fmt.Scanln(&answer)
			answer = strings.TrimSpace(answer)
			if answer == q.answer {
				correct++
			}
		}
		done <- 1
	}()

	timeout := time.After(*limit)

	select {
	case <-done:
		fmt.Printf("\n\nDone!\n\n")
	case <-timeout:
		fmt.Printf("\n\nTimeout!\n\n")
	}

	fmt.Printf("Total: %d\n", len(questions))
	fmt.Printf("Correct: %d\n", correct)
	fmt.Printf("Incorrect: %d\n\n", total-correct)
}

func loadQuestions(filename string) []question {
	questions := []question{}
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	lines, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range lines {
		questions = append(questions, question{line[0], strings.TrimSpace(line[1])})
	}

	return questions
}
