package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var allpaths = []string{}
var allpathscleaned = []string{}
var allallall = []string{}
var allpackagefucker = [][]string{}

// Function to find all paths from start to end
func findAllPath(links []string, start, end string, visited map[string]bool, path []string) {
	// Mark the current node as visited and add it to the path
	visited[start] = true
	path = append(path, start)
	// If the start node is the end node, print the path
	if start == end {
		allpaths = append(allpaths, strings.Join(path, "-"))
	} else {
		// Recur for all the vertices adjacent to this vertex
		for _, link := range links {
			first := strings.Split(link, "-")[0]
			second := strings.Split(link, "-")[1]
			if first == start && !visited[second] {
				findAllPath(links, second, end, visited, path)
			}
			if second == start && !visited[first] {
				findAllPath(links, first, end, visited, path)
			}
		}
	}
	// Remove current vertex from path and mark it as unvisited
	path = path[:len(path)-1]
	visited[start] = false
}

func main() {
	links := []string{"0-1", "2-4", "1-4", "0-2", "4-5", "3-0", "4-3", "0-5"}
	start := "0"
	end := "5"
	// Create a map to keep track of visited nodes
	visited := make(map[string]bool)
	// Call the function for each starting link
	findAllPath(links, start, end, visited, []string{})
	fmt.Println("hello", links)
	fmt.Println("hello", start)
	fmt.Println("hello", end)
	fmt.Println("hello", visited)
	fmt.Println("asdas", allpaths)
	allpaths = []string{"start-t-E-a-m-n-e-end", "start-t-E-a-m-n-h-A-c-k-end", "start-h-A-c-k-end", "start-0-o-e-n-end", "start-t-E-a-m-end", "start-0-o-n-m-end", "start-0-o-n-h-A-c-k-end"}
	// allpaths= []string{"0-2-3-1", "0-1", "0-2-4-3-6-1", "0-99-1", "0-2-1", "0-99-89-1", "0-3-1"}
	// allpaths= []string{ "start-0-o-e-n-end" ,"start-t-E-a-m-end" ,"start-0-o-n-m-end","start-h-A-c-k-end"}
	// fmt.Println(reflect.DeepEqual(m1, m2))
	allpaths = DoingTheWholeJob()
	old := []string{}
	DoingTheWholeJob()
	for reflect.DeepEqual(old, allallall) {
		old = DoingTheWholeJob()
	}
	for i := 0; i < len(allpaths); i++ {
		for j := 0; j < len(allpaths); j++ {
			temp := allpaths[j]
			if j+1 == len(allpaths) {
				break
			}
			allpaths[j] = allpaths[j+1]
			allpaths[j+1] = temp
			allpaths = DoingTheWholeJob()
			allpackagefucker = append(allpackagefucker, allpaths)
		}

	}
	fmt.Println("allpackagefucker", allpackagefucker)

	// fmt.Println("allpaths", allpaths)
	// allpaths = DoingTheWholeJob()
	fmt.Println("allpaths", allpaths)
	SimulateAnts(allpaths, 10)
	fmt.Println("alalalala:", allallall)
}

func DoingTheWholeJob() []string {
	var OUTPUT []string
	allallall = []string{}
	for i := 0; i < len(allpaths); i++ {
		allpathscleaned = []string{}
		findSmallestNonDisjointPaths(allpaths, i)
		maybe := findTheSmallestNumberOfRooms(allpathscleaned)
		fmt.Println("asa", maybe)
		allallall = append(allallall, maybe)
		allpathscleaned = []string{}
		fmt.Println("fina;", allallall)
		OUTPUT = RemoveDuplicateIFSeen(allallall)
		fmt.Println("OUTPUT", OUTPUT)
	}
	return OUTPUT
}

func findSmallestNonDisjointPaths(paths []string, index int) {
	for i := 0; i < len(paths); i++ {
		fmt.Println(i)
		line := paths[index]
		line2 := paths[i]
		// if i ==0 {
		// 		continue
		// 	}
		fmt.Println("\n\n\n\n\n\nline 1", line)
		fmt.Println("line 2", line2)
		rooms := strings.Split(line, "-")
		rooms2 := strings.Split(line2, "-")
		for i := 1; i < len(rooms)-1; i++ {
			for j := 1; j < len(rooms2)-1; j++ {
				if len(rooms) == 0 || len(rooms2) == 0 {
					continue
				}
				fmt.Println("length of rooms", len(rooms))
				fmt.Println("the rooms", rooms)
				fmt.Println("length of rooms2", len(rooms2))
				fmt.Println("the rooms2", rooms2)
				fmt.Printf("%s %s\n", rooms[i], rooms2[j])
				if rooms[i] == rooms2[j] {
					fmt.Println("length of rooms", len(rooms))
					fmt.Println("the rooms", rooms)
					fmt.Println("length of rooms2", len(rooms2))
					fmt.Println("the rooms2", rooms2)
					if len(rooms) >= len(rooms2) {
						rooms = nil
					}
				}
			}
		}
		fmt.Println("asas")
		if rooms == nil {
			fmt.Println("smallest non disjoint paths", rooms2)
			allpathscleaned = append(allpathscleaned, strings.Join(rooms2, "-"))
		} else {
			fmt.Println("smallest non disjoint paths", rooms)
			allpathscleaned = append(allpathscleaned, strings.Join(rooms, "-"))
		}
	}
}

func findTheSmallestNumberOfRooms(allpaths []string) string {
	smallest := allpaths[0]
	for i := 0; i < len(allpaths); i++ {
		if strings.Count(allpaths[i], "-") < strings.Count(smallest, "-") {
			smallest = allpaths[i]
		}
	}
	return smallest
}

func RemoveDuplicateIFSeen(array []string) []string {
	seen := make(map[string]bool)
	j := 0
	for i := 0; i < len(array); i++ {
		if seen[array[i]] {
			continue
		}
		seen[array[i]] = true
		array[j] = array[i]
		j++
	}
	return array[:j]
}

func SimulateAnts(paths []string, numAnts int) {
	mapForLengthOfEachPAth := make(map[string]int)
	MapPathForAnt := make(map[string]string)
	for i := 0; i < len(paths); i++ {
		mapForLengthOfEachPAth[paths[i]] = strings.Count(paths[i], "-")
	}
	fmt.Println("mapForLengthOfEachPAth", mapForLengthOfEachPAth)
	// for i := 0; i < len(paths); i++ {

	// }
	for i := 0; i < numAnts; i++ {
		smallest := FindSmallestIntInMap(mapForLengthOfEachPAth)
		// fmt.Printf("ant %v has path %v\n", i+1, smallest)
		MapPathForAnt["L"+strconv.Itoa(i+1)] = smallest
		mapForLengthOfEachPAth[smallest]++
	}
	fmt.Println("MapPathForAnt", MapPathForAnt)
	printFormattedMap(MapPathForAnt)
	fmt.Println(getFormattedMap(MapPathForAnt))
	arr := getFormattedMap(MapPathForAnt)
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			fmt.Println(arr[i][j])
		}
	}
	max := len(allallall)
	min := 0
	size := len(arr)
	for i := min; i < size; i += 3 {
		appendToFile("hh.txt", arr, max, i)
		min = i
		max += 3
		if max > size {
			max = size
		}
	}
}

func appendToFile(filename string, formattedMap [][]string, max, min int) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	maxLen := 0
	for _, row := range formattedMap {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}

	for i := 0; i < maxLen; i++ {
		for j := min; j < max; j++ {
			if i < len(formattedMap[j]) {
				_, err := writer.WriteString(formattedMap[j][i] + " ")
				if err != nil {
					return err
				}
			}
		}

		_, err := writer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func ReadLine(filePath string, lineNumber int) (string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    currentLine := 1

    for scanner.Scan() {
        if currentLine == lineNumber {
            return scanner.Text(), nil
        }
        currentLine++
    }

    if err := scanner.Err(); err != nil {
        return "", err
    }

    return "", fmt.Errorf("line %d does not exist", lineNumber)
}

func FindSmallestIntInMap(mapForLengthOfEachPAth map[string]int) string {
	var smallest string
	var smallestInt int
	for key, value := range mapForLengthOfEachPAth {
		if smallestInt == 0 {
			smallestInt = value
			smallest = key
			continue
		}
		if value < smallestInt {
			smallestInt = value
			smallest = key
		}
	}
	return smallest
}

func printFormattedMap(MapPathForAnt map[string]string) {
	for i := 1; i <= len(MapPathForAnt); i++ {
		key := "L" + strconv.Itoa(i)
		path := MapPathForAnt[key]
		elements := strings.Split(path, "-")
		for _, element := range elements {
			fmt.Printf("%s-%s ", key, element)
		}

		fmt.Println()
	}
}

func getFormattedMap(MapPathForAnt map[string]string) [][]string {
	// Create a 2D slice to hold the results
	formattedMap := make([][]string, len(MapPathForAnt))

	for i := 1; i <= len(MapPathForAnt); i++ {
		key := "L" + strconv.Itoa(i)
		path := MapPathForAnt[key]
		elements := strings.Split(path, "-")

		// Create a slice to hold the formatted strings for this key
		formattedRow := make([]string, len(elements))
		for j, element := range elements {
			formattedRow[j] = fmt.Sprintf("%s-%s", key, element)
		}

		// Store the formatted row in the 2D slice
		formattedMap[i-1] = formattedRow
	}

	return formattedMap
}
