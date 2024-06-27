package lemin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFile(filepath string) *AntFarm {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("ERROR: failed to open file.")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	Farm := &AntFarm{
		rooms: make(map[string]*Room),
	}

	var phases string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if Farm.AntsNum == 0 {
			Farm.AntsNum, err = strconv.Atoi(line)
			if err != nil {
				fmt.Println("ERROR: invalid data format for number of ants")
				os.Exit(0)
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				phases = "start"
			} else if line == "##end" {
				phases = "end"
			}
			continue
		}

		if strings.Contains(line, " ") {
			coordenaties := strings.Split(line, " ")
			if len(coordenaties) != 3 {
				fmt.Println("ERROR: Data format for room coordinates.")
				os.Exit(0)
			}
			xcoord, err := strconv.Atoi(coordenaties[1])
			if err != nil {
				fmt.Println("ERROR: Data format for coordinates.")
				os.Exit(0)
			}
			ycoord, err := strconv.Atoi(coordenaties[2])
			if err != nil {
				fmt.Println("ERROR: Data format for coordinates.")
				os.Exit(0)
			}
			room := &Room{
				name:  coordenaties[0],
				x:     xcoord,
				y:     ycoord,
				links: []*Room{},
			}
			Farm.rooms[room.name] = room
			if phases == "start" {
				Farm.startRoom = room
				phases = ""
			} else if phases == "end" {
				Farm.end = room
				phases = ""
			}
			continue
		}

		if strings.Contains(line, "-") {
			coordenaties := strings.Split(line, "-")
			if len(coordenaties) != 2 {
				fmt.Println("ERROR: Data format for links.")
				os.Exit(0)
			}
			room1, ok1 := Farm.rooms[coordenaties[0]]
			room2, ok2 := Farm.rooms[coordenaties[1]]
			if !ok1 || !ok2 {
				fmt.Println("ERROR: Room not found.")
				os.Exit(0)
			}
			room1.links = append(room1.links, room2)
			room2.links = append(room2.links, room1)
		}
	}
	err = scanner.Err()
	if err != nil {
		fmt.Println("ERROR: Scanner.")
		os.Exit(0)
	}
	if Farm.startRoom == nil || Farm.end == nil {
		fmt.Println("ERROR: Missing start or end.")
		os.Exit(0)
	}
	return Farm
}
