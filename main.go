package main

import (
	"bufio"
	"fmt"
	lemin "lemin/Functions"
	"os"
)


func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("ERROR: invalid number of arguments")
		return
	}
	filepath := os.Args[1]

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("ERROR: failed to open file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}

	for _, line := range fileLines {
		fmt.Println(line)
	}
	fmt.Println()

	antFarm, err := lemin.ParseFile(filepath)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	PathsNeeded := lemin.BFS(antFarm)
	if len(PathsNeeded) == 0 {
		fmt.Println("ERROR: no valid paths found")
		return
	}
	lemin.AntsMovement(PathsNeeded, antFarm)
}
