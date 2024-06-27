package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Room in the ant farm
type Room struct {
	Name  string
	X, Y  int
	Links []*Room
}

// An ant
type Ant struct {
	ID   int
	Room *Room
}

// Ant farm
type Farm struct {
	Rooms   map[string]*Room
	Start   *Room
	End     *Room
	NumAnts int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no input file provided")
		return
	}
	filename := os.Args[1]
	farm, err := ParseFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(farm.NumAnts)
	for _, room := range farm.Rooms {
		fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
	}
	for _, room := range farm.Rooms {
		for _, link := range room.Links {
			if room.Name < link.Name {
				fmt.Printf("%s-%s\n", room.Name, link.Name)
			}
		}
	}
	SimulateAnts(farm)
}

// Initializes a new Farm
func NewFarm() *Farm {
	return &Farm{
		Rooms: make(map[string]*Room),
	}
}

// add a room to the farm
func (f *Farm) AddRoom(name string, x, y int) error {
	if _, exists := f.Rooms[name]; exists {
		return fmt.Errorf("ERROR: invalid data format, duplicated rooms")
	}
	if strings.HasPrefix(name, "L") || strings.Contains(name, " ") {
		return fmt.Errorf("ERROR: invalid data format, invalid room name")
	}
	f.Rooms[name] = &Room{Name: name, X: x, Y: y, Links: []*Room{}}
	return nil
}

// add a link between two rooms
func (f *Farm) AddLink(name1, name2 string) error {
	room1, ok1 := f.Rooms[name1]
	room2, ok2 := f.Rooms[name2]
	if !ok1 || !ok2 {
		return fmt.Errorf("ERROR: invalid data format, links to unknown rooms")
	}
	for _, link := range room1.Links {
		if link == room2 {
			return fmt.Errorf("ERROR: invalid data format, duplicate links")
		}
	}
	room1.Links = append(room1.Links, room2)
	room2.Links = append(room2.Links, room1)
	return nil
}

// set the start room
func (f *Farm) SetStart(name string) error {
	room, ok := f.Rooms[name]
	if !ok {
		return fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	f.Start = room
	return nil
}

// set the end room
func (f *Farm) SetEnd(name string) error {
	room, ok := f.Rooms[name]
	if !ok {
		return fmt.Errorf("ERROR: invalid data format, no end room found")
	}
	f.End = room
	return nil
}

// Parsing File

func ParseFile(filename string) (*Farm, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	farm := NewFarm()
	scanner := bufio.NewScanner(file)
	var numAnts int
	var phase string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if phase == "" {
			if _, err := fmt.Sscanf(line, "%d", &numAnts); err == nil {
				if numAnts <= 0 {
					return nil, fmt.Errorf("ERROR: invalid data format, invalid number of Ants")
				}
				farm.NumAnts = numAnts
				phase = "rooms"
			} else {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid number of Ants")
			}
		} else if phase == "rooms" {
			if line == "##start" {
				phase = "start"
			} else if line == "##end" {
				phase = "end"
			} else if strings.Contains(line, "-") {
				phase = "links"
				name1, name2 := "", ""
				names := strings.Split(line, "-")
				if len(names) != 2 {
					return nil, fmt.Errorf("ERROR: invalid data format, too many rooms in link")
				}
				name1, name2 = names[0], names[1]
				if err := farm.AddLink(name1, name2); err != nil {
					return nil, err
				}
			} else {
				name := ""
				var x, y int
				if _, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y); err == nil {
					if err := farm.AddRoom(name, x, y); err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("ERROR: invalid data format, invalid room format")
				}
			}
		} else if phase == "start" {
			name := ""
			var x, y int
			if _, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y); err == nil {
				if err := farm.AddRoom(name, x, y); err != nil {
					return nil, err
				}
				if err := farm.SetStart(name); err != nil {
					return nil, err
				}
				phase = "rooms"
			} else {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid start room format")
			}
		} else if phase == "end" {
			name := ""
			var x, y int
			if _, err := fmt.Sscanf(line, "%s %d %d", &name, &x, &y); err == nil {
				if err := farm.AddRoom(name, x, y); err != nil {
					return nil, err
				}
				if err := farm.SetEnd(name); err != nil {
					return nil, err
				}
				phase = "rooms"
			} else {
				return nil, fmt.Errorf("ERROR: invalid data format, invalid end room format")
			}
		} else if phase == "links" {
			name1, name2 := "", ""
			names := strings.Split(line, "-")
			if len(names) != 2 {
				return nil, fmt.Errorf("ERROR: invalid data format, too many rooms in link")
			}
			name1, name2 = names[0], names[1]

			if err := farm.AddLink(name1, name2); err != nil {
				return nil, err
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if farm.Start == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no start room found")
	}
	if farm.End == nil {
		return nil, fmt.Errorf("ERROR: invalid data format, no end room found")
	}
	return farm, nil
}

// Path finding
func Pathfinding(start, end *Room) ([]*Room, error) {
	Q := [][]*Room{{start}}
	visited := map[*Room]bool{start: true}
	for len(Q) > 0 {
		path := Q[0]
		Q = Q[1:]
		last := path[len(path)-1]
		if last == end {
			return path, nil
		}
		for _, link := range last.Links {
			if !visited[link] {
				visited[link] = true
				newPath := append([]*Room{}, path...)
				newPath = append(newPath, link)
				Q = append(Q, newPath)
			}
		}
	}
	return nil, fmt.Errorf("ERROR: no path found from start to end")
}

// Ant Movement
func SimulateAnts(farm *Farm) {
	ants := make([]*Ant, farm.NumAnts)
	for i := range ants {
		ants[i] = &Ant{ID: i + 1, Room: farm.Start}
	}

	allReached := false

	for !allReached {
		moveStr := ""
		allReached = true

		for _, ant := range ants {
			if ant.Room != farm.End {
				allReached = false
				path, err := Pathfinding(ant.Room, farm.End)
				if err == nil && len(path) > 1 {
					ant.Room = path[1]
					moveStr += fmt.Sprintf("L%d-%s ", ant.ID, ant.Room.Name)
				}
			}
		}

		if moveStr != "" {
			fmt.Println(strings.TrimSpace(moveStr))
		}
	}
}
