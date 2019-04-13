package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

// Example usage:
// $ sudo ping -i 0.01 -t 5 192.168.1.1 | ts -s '[%.s]' | sed -l -n -e 's/\[\(.*\)\].*time=\(.*\) ms/\1 \2/p' | ./ping-heatmap > data.json

func main() {
	// TODO: Accept a duration flag

	timestamps := make([]string, 0)
	latenciesMs := make([]string, 0)

	c := make(chan os.Signal, 1)

	go func () {
		fmt.Fprintln(os.Stderr, "Up and running")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			tokens := strings.Split(scanner.Text(), " ")
			timestamps = append(timestamps, tokens[0])
			latenciesMs = append(latenciesMs, tokens[1])
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}

		fmt.Fprintln(os.Stderr, "Done with loop")
		c<- os.Interrupt
	}()

	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Fprintln(os.Stderr, "Assembling heatmap data...")

	if err := os.Stdin.Close(); err != nil {
		log.Fatalf("stdin: %v", err)
	}

	fmt.Println("{")

	rows := make([]string, 100)
	for i := 0; i < 100; i++ {
		rows[i] = strconv.Itoa(i)
	}

	rowsJson := strings.Join(rows, ",")
	fmt.Printf("\t\"rows\": [%s],\n", rowsJson)


	columns := make([]string, len(timestamps)/100 + 1)
	for i := 0; i < len(columns); i++ {
		columns[i] = strconv.Itoa(i)
	}
	columnsJson := strings.Join(columns, ",")
	fmt.Printf("\t\"columns\": [%s],\n", columnsJson)

	fmt.Print("\t\"values\": [")

	for i := 0; i < len(columns); i++ {
		var vals []string
		if (i+1) * 100 < len(latenciesMs) {
			vals = latenciesMs[i*100:(i+1)*100]
		} else {
			vals = latenciesMs[i*100:]
		}
		fmt.Printf("\t[%s]", strings.Join(vals, ","))

		if len(vals) == 100 {
			fmt.Println(",")
		}

		fmt.Print("\n")
	}

	fmt.Print("\t]")

	fmt.Println("}")
}
