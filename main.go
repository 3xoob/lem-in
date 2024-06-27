package main

import (
	"bufio"
	"fmt"
	lemin "lemin/Functions"
	"os"
)

func main() {
	arg := os.Args
	if len(arg) != 2 {
		fmt.Println("ERROR: invalid number of arguments")
		return
	}
	filepath := arg[1]

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("ERROR: failed to open file.")
		return
	}
	defer file.Close()

	Scan := bufio.NewScanner(file)
	var FileContent []string
	for Scan.Scan() {
		FileContent = append(FileContent, Scan.Text())
	}

	for _, line := range FileContent {
		fmt.Println(line)
	}
	fmt.Println()

	FarmOfAnt := lemin.ParseFile(filepath)

	PathsNeeded := lemin.BFS(FarmOfAnt)
	if len(PathsNeeded) == 0 {
		fmt.Println("ERROR: no valid paths found")
		return
	}
	lemin.AntsMovement(PathsNeeded, FarmOfAnt)
}
