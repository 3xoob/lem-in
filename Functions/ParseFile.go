package lemin

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseFile(filepath string) (*AntFarm, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("ERROR: failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	Farm := &AntFarm{
		rooms: make(map[string]*Room),
	}

	var pendingType string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if Farm.numberOfAnts == 0 {
			Farm.numberOfAnts, err = strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format for number of ants")
			}
			continue
		}
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				pendingType = "start"
			} else if line == "##end" {
				pendingType = "end"
			}
			continue
		}

		if strings.Contains(line, " ") {
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				return nil, fmt.Errorf("ERROR: invalid data format for room coordinates")
			}
			x, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format for coordinates")
			}
			y, err := strconv.Atoi(parts[2])
			if err != nil {
				return nil, fmt.Errorf("ERROR: invalid data format for coordinates")
			}
			room := &Room{
				name:  parts[0],
				x:     x,
				y:     y,
				links: []*Room{},
			}
			Farm.rooms[room.name] = room
			if pendingType == "start" {
				Farm.startRoom = room
				pendingType = ""
			} else if pendingType == "end" {
				Farm.endRoom = room
				pendingType = ""
			}
			continue
		}

		if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("ERROR: invalid data format for link")
			}
			room1, ok1 := Farm.rooms[parts[0]]
			room2, ok2 := Farm.rooms[parts[1]]
			if !ok1 || !ok2 {
				return nil, fmt.Errorf("ERROR: invalid link: room not found")
			}
			room1.links = append(room1.links, room2)
			room2.links = append(room2.links, room1)
		}
	}
	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("ERROR: scanner error")
	}
	if Farm.startRoom == nil || Farm.endRoom == nil {
		return nil, fmt.Errorf("ERROR: missing start room or end room")
	}
	return Farm, nil
}
