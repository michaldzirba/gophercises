package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	out "fmt"
	"io"
	"log"
	"os"
	"time"
	// "strings"
)

const (
	default_problems   = "problems.csv"
	default_timeout    = 30
	correct, incorrect = 1, 0
)

var (
	timeout_  int
	problems_ string
)

type Quiz struct {
	question, answer string
}

func (q *Quiz) test(given string) (score int) {
	if q.answer == given {
		return correct
	} else {
		return incorrect
	}
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
func NewQuiz(record []string) (quiz *Quiz) {
	if len(record) > 1 {
		quiz := Quiz{record[0], record[1]}

		return &quiz
	}
	panic("incorrect input data")
}

func (q *Quiz) String() string {
	return "quiz: " + q.question + " = " + q.answer
}

func test(q *Quiz) (score int) {
	out.Println(q.question)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() && q.answer == scanner.Text() {
		return correct
	}

	return incorrect
}

func init() {
	timeoutPtr := flag.Int("time", default_timeout, "how long will the program give for answers")
	problemsPtr := flag.String("problems", default_problems, "where to read the problems from")
	flag.Parse()
	timeout_ = *timeoutPtr
	problems_ = *problemsPtr
}

func main() {
	totalscore, maxscore := 0, 0

	timer := time.NewTimer(time.Second * time.Duration(timeout_)).C
	f, err := os.Open(problems_)
	check(err)
	reader := csv.NewReader(bufio.NewReader(f))
	defer f.Close() // close the file, at the end of method

	scoreCh, endCh := make(chan int), make(chan int)

	go func() {
		for {
			record, err := reader.Read()
			if io.EOF == err {
				endCh <- 0
				break
			}
			check(err)                       // other errors cause panic
			scoreCh <- test(NewQuiz(record)) // this parses, asks and returns points if correct (0 if not correct)
		}
	}()

loop:
	for {
		select {
		case <-timer:
			break loop
		case <-endCh:
			break loop
		case score := <-scoreCh:
			maxscore += correct
			totalscore += score
		}
	}

	out.Println("end, score", totalscore, " / ", maxscore)
}
