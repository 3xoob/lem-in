package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var FileLines []string
var start, end int
var links []string
var rooms []string
var AllPaths = []string{}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: Please provide the input file as an argument.")
		fmt.Println("USAGE: go run . <FILENAME>")
		os.Exit(0)
	}

	input := os.Args[1]
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("ERROR: Invalid file.")
		os.Exit(0)
	}
	he, _ := os.ReadFile(input)
	if len(he) == 0 {
		fmt.Println("ERROR: File Empty.")
		os.Exit(0)
	}
	defer file.Close()

	Scan := bufio.NewScanner(file)
	for Scan.Scan() {
		line := Scan.Text()
		FileLines = append(FileLines, line)
		if line == "##start" {
			start = len(FileLines) - 1
		}
		if line == "##end" {
			end = len(FileLines) - 1
		}
	}

	AmountOfAnt := FileLines[0]

	for _, line := range FileLines[1:] {
		if !strings.Contains(line, "#") && !strings.Contains(line, "-") {
			rooms = append(rooms, line)
		} else if strings.Contains(line, "-") {
			links = append(links, line)
		}
	}

	for _, edge := range links {
		vertexes := strings.Split(edge, "-")
		if len(vertexes) != 2 {
			fmt.Println("ERROR: Links has more than 2.")
			os.Exit(0)
		}
	}

	fmt.Println()
	visited := make(map[string]bool)
	fmt.Println(links)
	fmt.Println(start)
	fmt.Println(end)
	fmt.Println(visited)
	fmt.Println(AmountOfAnt)

}
