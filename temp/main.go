package main

import (
	"fmt"
	"strings"
)

// Function to find all paths from start to end
func findAllPath(links []string, start, end string, visited map[string]bool, path []string) {
	// Mark the current node as visited and add it to the path
	visited[start] = true
	path = append(path, start)

	// If the start node is the end node, print the path
	if start == end {
		fmt.Println(strings.Join(path, "-"))
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
	links := []string{"0-2", "0-0", "2-1"}
	start := "0"
	end := "1"

	// Create a map to keep track of visited nodes
	visited := make(map[string]bool)

	// Call the function for each starting link
	findAllPath(links, start, end, visited, []string{})
	
}
