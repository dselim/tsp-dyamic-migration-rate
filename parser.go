package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func parseFile(filepath string) []point {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var (
		line        string
		canRead     bool
		problemSize int
	)
	for canRead = scanner.Scan(); canRead; canRead = scanner.Scan() {
		line := scanner.Text()
		if line == "NODE_COORD_SECTION" {
			break
		}
		tokens := strings.Split(line, " : ")
		if len(tokens) == 2 && tokens[0] == "DIMENSION" {
			problemSize, err = strconv.Atoi(tokens[1])
			if err != nil {
				panic(err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	if !canRead {
		panic("no coordinates to read")
	}

	nodes := make([]point, 0, problemSize)
	for scanner.Scan() {
		line = scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) != 3 {
			break
		}

		x, err := strconv.ParseFloat(tokens[1], 64)
		if err != nil {
			panic(err)
		}
		y, err := strconv.ParseFloat(tokens[2], 64)
		if err != nil {
			panic(err)
		}

		nodes = append(nodes, point{x, y})
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return nodes
}
