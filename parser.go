package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	Next chan LatencyDataPoint
}

func NewParser(r io.Reader) *Parser {
	p := &Parser{
		Next: make(chan LatencyDataPoint, 100),
	}

	go func () {
		defer close(p.Next)

		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			tokens := strings.Split(scanner.Text(), " ")
			offset, err := strconv.ParseFloat(tokens[0], 64)
			if err != nil {
				fmt.Printf("offset parse error for '%s': %v", tokens[0], err)
			}
			latency, err := strconv.ParseFloat(tokens[1], 64)
			if err != nil {
				fmt.Printf("latency parse error for '%s': %v", tokens[1], err)
			}

			p.Next<- LatencyDataPoint{offset, latency}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
		}
	}()

	return p
}