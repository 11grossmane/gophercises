package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	timePointer := flag.Int("t", 2, "set timer")
	flag.Parse()
	startTime := time.Now()
	fmt.Printf("Welcome to the quiz.  You have %d seconds to complete it \n", *timePointer)
	filePointer, _ := os.Open("./problems.csv")
	rd := csv.NewReader(filePointer)
	score := 0
	total := 0

	fmt.Println("press enter to start")
	fmt.Scanln()

	for {
		record, err := rd.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("What is " + record[0] + "?")
		var answer string
		fmt.Scanln(&answer)
		if time.Since(startTime).Seconds() > float64(*timePointer) {
			fmt.Println("You ran out of time :(")
			return
		}
		myanswer, _ := strconv.Atoi(answer)
		solution, _ := strconv.Atoi(record[1])
		if myanswer == solution {
			fmt.Println("correct")
			score++
		} else {
			fmt.Println("incorrect")
		}
		total++
	}
	fmt.Printf("You got %d/%d \n", score, total)

}
