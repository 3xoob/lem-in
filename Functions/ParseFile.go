package lemin

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

	Farm, err := parseAntFarm(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	return Farm
}

func parseAntFarm(r io.Reader) (*AntFarm, error) {
	scanner := bufio.NewScanner(r)
	Farm := &AntFarm{
		rooms: make(map[string]*Room),
	}

	var phase string
	antsSet := false
	linkSeen := make(map[string]bool)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if !antsSet {
			n, err := strconv.Atoi(line)
			if err != nil || n <= 0 {
				return nil, errors.New("ERROR: invalid data format, invalid number of Ants")
			}
			Farm.AntsNum = n
			antsSet = true
			continue
		}
		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				phase = "start"
			} else if line == "##end" {
				phase = "end"
			}
			continue
		}

		if strings.Contains(line, " ") {
			fields := strings.Split(line, " ")
			if len(fields) != 3 {
				return nil, errors.New("ERROR: invalid data format, room coordinates.")
			}
			name := fields[0]
			if _, exists := Farm.rooms[name]; exists {
				return nil, errors.New("ERROR: invalid data format, duplicate room name.")
			}
			x, err := strconv.Atoi(fields[1])
			if err != nil {
				return nil, errors.New("ERROR: invalid data format, room coordinates.")
			}
			y, err := strconv.Atoi(fields[2])
			if err != nil {
				return nil, errors.New("ERROR: invalid data format, room coordinates.")
			}
			room := &Room{name: name, x: x, y: y, links: []*Room{}}
			Farm.rooms[name] = room
			switch phase {
			case "start":
				if Farm.start != nil {
					return nil, errors.New("ERROR: invalid data format, multiple start rooms.")
				}
				Farm.start = room
				phase = ""
			case "end":
				if Farm.end != nil {
					return nil, errors.New("ERROR: invalid data format, multiple end rooms.")
				}
				Farm.end = room
				phase = ""
			}
			continue
		}

		if strings.Contains(line, "-") {
			fields := strings.Split(line, "-")
			if len(fields) != 2 {
				return nil, errors.New("ERROR: invalid data format, links.")
			}
			a, b := fields[0], fields[1]
			if a == b {
				return nil, errors.New("ERROR: invalid data format, room cannot link to itself.")
			}
			room1, ok1 := Farm.rooms[a]
			room2, ok2 := Farm.rooms[b]
			if !ok1 || !ok2 {
				return nil, errors.New("ERROR: invalid data format, room not found.")
			}
			key := linkKey(a, b)
			if linkSeen[key] {
				return nil, errors.New("ERROR: invalid data format, duplicate link.")
			}
			linkSeen[key] = true
			room1.links = append(room1.links, room2)
			room2.links = append(room2.links, room1)
			continue
		}

		return nil, errors.New("ERROR: invalid data format, unrecognized line.")
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("ERROR: invalid data format, scanner error.")
	}
	if Farm.start == nil || Farm.end == nil {
		return nil, errors.New("ERROR: invalid data format, missing start or end room.")
	}
	return Farm, nil
}

func linkKey(a, b string) string {
	if a < b {
		return a + "-" + b
	}
	return b + "-" + a
}
