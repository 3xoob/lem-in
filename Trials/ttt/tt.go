package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {

	input := os.Args[1]
	new_farm, err := PrepFarm(input)
	if err != nil {
		fmt.Errorf("new farm err %w", err)
		return
	}

	fmt.Println("=================================")
	fmt.Println("Start: ", new_farm.Start)
	fmt.Println("End: ", new_farm.End)
	fmt.Println("antAmunt: ", new_farm.AntAmount)
	fmt.Println("AdjacencyList: ", new_farm.AdjacencyList)
	fmt.Println("=================================")
	fmt.Println()
	fmt.Println("path:        ", new_farm.Paths)
	fmt.Println("AntAmaunt:        ", new_farm.AntAmount)
	// fmt.Println("queue result: ", Queue(new_farm.Paths, new_farm.AntAmount))

}

func DFS(adjacencyList map[string][]string, currentVertex string, listOfVisited []string) []string {
	listOfVisited = append(listOfVisited, currentVertex)
	childs := adjacencyList[currentVertex]
	if len(childs) == 0 {
		return listOfVisited
	}

	for _, child := range childs {
		if !isContain(listOfVisited, child) {
			listOfVisited = DFS(adjacencyList, child, listOfVisited)
		}
	}

	return listOfVisited
}

// checks whether the array contains an element
func isContain(arr []string, target string) bool {
	if len(arr) == 0 {
		return false
	}

	for _, elem := range arr {
		if elem == target {
			return true
		}
	}

	return false
}

// sorted farm
type UpdatedFarm struct {
	AntAmount     int
	Start         string
	End           string
	AdjacencyList map[string][]string
	Weights       map[[2]string]bool
	Queue         []int
	Paths         [][]string
}

// base farm
type Farm struct {
	AntAmount string
	Start     string
	End       string
	Links     []string
	Rooms     []string
}

// farm parser
func PrepFarm(filename string) (result UpdatedFarm, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return UpdatedFarm{}, fmt.Errorf("prepFarm: %w", err)
	}
	defer f.Close()

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	var fileLines []string
	var idxOfStart, idxOfEnd int
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "##start" {
			idxOfStart = len(fileLines)
		}
		if line == "##end" {
			idxOfEnd = len(fileLines)
		}
		fileLines = append(fileLines, line)
	}
	if len(fileLines) == 0 {
		return UpdatedFarm{}, fmt.Errorf("prepFarm: empty file")
	}

	result.AntAmount, err = strconv.Atoi(fileLines[0])
	if err != nil {
		return UpdatedFarm{}, fmt.Errorf("prepFarm: failed to parse ant amount: %w", err)
	}

	// var emptyV Vertex
	var links []string
	var rooms []string

	for i, line := range fileLines[1:] {
		switch i {
		case idxOfStart, idxOfEnd:
			if i == idxOfStart {
				result.Start = line
			} else {
				result.End = line
			}
		default:
			if strings.Contains(line, "-") {
				links = append(links, line)
			} else if !strings.Contains(line, "#") {
				rooms = append(rooms, line)
			}
		}
	}

	result.AdjacencyList = TransformToAdjacencyList(links)
	if len(result.AdjacencyList) == 0 {
		return UpdatedFarm{}, fmt.Errorf("prepFarm: failed to transform links field")
	}

	return result, nil
}

func UpdateFarm(raw_farm Farm) (result UpdatedFarm, err error) {
	result.AntAmount, err = strconv.Atoi(raw_farm.AntAmount)
	if err != nil {
		return result, fmt.Errorf("UpdateFarm: %w", err)
	}

	result.Start = GetName(raw_farm.Start)

	result.End = GetName(raw_farm.End)

	result.AdjacencyList = TransformToAdjacencyList(raw_farm.Links)
	if len(result.AdjacencyList) == 0 {
		return UpdatedFarm{}, fmt.Errorf("UpdateFarm: can't transform links field of Farm")
	}

	result.Weights = make(map[[2]string]bool)

	return result, nil
}

func GetName(info string) string {
	splittedData := strings.Split(info, " ")
	if len(splittedData) == 0 {
		return ""
	}

	return splittedData[0]
}

// function that transform list of edges to adjacency list
func TransformToAdjacencyList(listOfEdges []string) map[string][]string {
	var result map[string][]string = make(map[string][]string)

	if len(listOfEdges) == 0 {
		return result
	}

	for _, pairOfVertex := range listOfEdges {
		vertexes := strings.Split(pairOfVertex, "-")
		result[vertexes[0]] = append(result[vertexes[0]], vertexes[1])
		result[vertexes[1]] = append(result[vertexes[1]], vertexes[0])
	}

	return result
}

func Queue(paths [][]string, antsAmount int) []int {

	antsQueue := make([]int, len(paths))
	rooms := make([]int, len(paths))

	for i := range paths {
		rooms[i] = len(paths[i])
	}

	for antsAmount > 0 {
		fmt.Println("rooms: ", rooms)
		fmt.Println("paths: ", len(paths))
		fmt.Println("antsQueue: ", antsQueue)
		fmt.Println("antsAmount: ", antsAmount)

		indexOfInsert := checkLowestPath(rooms, antsQueue)
		antsQueue[indexOfInsert] += 1
		antsAmount -= 1
	}
	return antsQueue
}

func checkLowestPath(rooms []int, antsQueue []int) int {
	lowestValue := 10000
	lowestInd := 0
	for indOfPath := range rooms {
		// summ of rooms and ants in one path
		sum := rooms[indOfPath] + antsQueue[indOfPath]
		if sum < lowestValue {
			lowestValue = sum
			lowestInd = indOfPath
		}
	}

	return lowestInd
}
