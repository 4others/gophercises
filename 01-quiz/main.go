package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type queAns struct {
	Question string
	Answer   int
}

var (
	csvFilePath      = "problems.csv"
	timeLimitDefault = 30
	done             chan bool
	questionsAnswers []queAns
	answered         int
)

func main() {
	timeLimitFlag := flag.Int("timeout", timeLimitDefault, "limit-time modifies time in seconds available for solving provided tasks")
	pathFlag := flag.String("csv", csvFilePath, "csv-path specifies path to csv file")
	flag.Parse()
	done := make(chan bool)

	csvFile, err := os.Open(*pathFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	for _, line := range csvLines {
		conv, err := strconv.Atoi(line[1])
		if err != nil {
			fmt.Println(err)
		}
		qA := queAns{
			Question: line[0],
			Answer:   conv,
		}
		questionsAnswers = append(questionsAnswers, qA)
	}

	questionsNumber := len(csvLines)

	fmt.Println("Press any key to start: ")
	var in string
	fmt.Scanln(&in)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeLimitFlag)*time.Second)
	defer cancel()

	go func() {
		for _, v := range questionsAnswers {
			var input int
			fmt.Printf("What is %v?\n", v.Question)
			fmt.Println("Type your answer: ")
			fmt.Scanln(&input)

			if v.Answer == input {
				answered++
			}
		}
		done <- true
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Timeout exceeded.")
		fmt.Printf("Available points: %v, collected points: %v\n", questionsNumber, answered)
		return
	case <-done:
		fmt.Printf("Available points: %v, collected points: %v\n", questionsNumber, answered)
		return
	}
}
