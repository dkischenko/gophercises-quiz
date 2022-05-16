package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		data [][]string
	)
	var (
		file    string
		limit   int
		total   int
		correct int
	)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Please, restart your program, you got an error: %s", err)
		}
	}()

	startQuiz()
	file, limit = getArguments()
	data = readFile(file)
	total, correct = run(data, limit)
	fmt.Printf("Total asked questions is %d\r\n", total)
	fmt.Printf("Total correct answers is %d\r\n", correct)
	os.Exit(0)
}

func run(data [][]string, limit int) (total int, correct int) {
	answerCh := make(chan string)
	defer close(answerCh)
	timer := time.NewTimer(0)
	if limit > 0 {
		timer.Reset(time.Second * time.Duration(limit))
		defer timer.Stop()
	}
	for _, line := range data {
		question, expected := line[0], strings.TrimSpace(line[1])
		fmt.Printf("Please, what %s, sir? ", question)
		total++
		go getAnswer(answerCh)
		select {
		case answer := <-answerCh:
			if answer == expected {
				fmt.Println("Correct answer!")
				correct++
			}
		case <-timer.C:
			fmt.Println("\r\nTimer stopped")
			return
		}
	}

	return
}

func getArguments() (string, int) {
	var (
		file  string
		limit int
	)

	flag.StringVar(&file, "file", "problems.csv", "filename for loading data")
	flag.IntVar(&limit, "limit", 30, "limit in seconds for answers")
	flag.Parse()
	return file, limit
}

func startQuiz() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("To start quiz press 'Enter'")
		char, _, err := reader.ReadRune()
		errNotNil(err)
		if char != 13 {
			fmt.Println("You pressed wrong key! Please, try again")
		} else {
			break
		}
	}
}

func readFile(file string) (data [][]string) {
	readFile, err := os.Open(file)
	errNotNil(err)
	defer readFile.Close()

	lines, err := csv.NewReader(readFile).ReadAll()
	errNotNil(err)
	return lines
}

func errNotNil(err error) {
	if err != nil {
		panic(err)
	}
}

func getAnswer(answerChan chan string) {
	r := bufio.NewReader(os.Stdin)
	answer, err := r.ReadString('\n')
	errNotNil(err)
	answerChan <- strings.TrimSpace(answer)
}
