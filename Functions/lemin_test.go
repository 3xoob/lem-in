package lemin

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func mustParse(t *testing.T, input string) *AntFarm {
	t.Helper()
	farm, err := parseAntFarm(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return farm
}

func TestParseAntFarm_Valid(t *testing.T) {
	input := "2\n##start\nstart 0 0\nmid 1 1\n##end\nend 2 2\nstart-mid\nmid-end\n"
	farm := mustParse(t, input)

	if farm.AntsNum != 2 {
		t.Errorf("AntsNum = %d, want 2", farm.AntsNum)
	}
	if farm.start == nil || farm.start.name != "start" {
		t.Errorf("start room = %v, want %q", farm.start, "start")
	}
	if farm.end == nil || farm.end.name != "end" {
		t.Errorf("end room = %v, want %q", farm.end, "end")
	}
	if len(farm.rooms) != 3 {
		t.Errorf("len(rooms) = %d, want 3", len(farm.rooms))
	}
	if len(farm.start.links) != 1 || farm.start.links[0].name != "mid" {
		t.Errorf("start.links = %v, want [mid]", farm.start.links)
	}
}

func TestParseAntFarm_Errors(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "zero ants",
			input: "0\n##start\ns 0 0\n##end\ne 1 1\ns-e\n",
			want:  "invalid number of Ants",
		},
		{
			name:  "negative ants",
			input: "-5\n##start\ns 0 0\n##end\ne 1 1\ns-e\n",
			want:  "invalid number of Ants",
		},
		{
			name:  "non-numeric ants",
			input: "abc\n##start\ns 0 0\n##end\ne 1 1\ns-e\n",
			want:  "invalid number of Ants",
		},
		{
			name:  "duplicate room name",
			input: "1\nA 0 0\nA 5 5\n##start\ns 0 0\n##end\ne 1 1\ns-e\n",
			want:  "duplicate room name",
		},
		{
			name:  "multiple start rooms",
			input: "1\n##start\ns1 0 0\n##start\ns2 1 1\n##end\ne 2 2\ns1-e\ns2-e\n",
			want:  "multiple start rooms",
		},
		{
			name:  "multiple end rooms",
			input: "1\n##start\ns 0 0\n##end\ne1 1 1\n##end\ne2 2 2\ns-e1\ns-e2\n",
			want:  "multiple end rooms",
		},
		{
			name:  "self link",
			input: "1\n##start\ns 0 0\n##end\ne 1 1\ns-e\ns-s\n",
			want:  "link to itself",
		},
		{
			name:  "duplicate link",
			input: "1\n##start\ns 0 0\n##end\ne 1 1\ns-e\ns-e\n",
			want:  "duplicate link",
		},
		{
			name:  "link to unknown room",
			input: "1\n##start\ns 0 0\n##end\ne 1 1\ns-ghost\n",
			want:  "room not found",
		},
		{
			name:  "missing start",
			input: "1\n##end\ne 1 1\n",
			want:  "missing start or end",
		},
		{
			name:  "missing end",
			input: "1\n##start\ns 0 0\n",
			want:  "missing start or end",
		},
		{
			name:  "malformed room line",
			input: "1\n##start\ns 0 0 0\n##end\ne 1 1\ns-e\n",
			want:  "room coordinates",
		},
		{
			name:  "malformed link line",
			input: "1\n##start\ns 0 0\n##end\ne 1 1\ns-e-x\n",
			want:  "links",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := parseAntFarm(strings.NewReader(tc.input))
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.want)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Errorf("error = %q, want it to contain %q", err.Error(), tc.want)
			}
		})
	}
}

// A diamond with two independent branches of different lengths between
// start and end: start must never see both branches share a room.
func TestBFS_RoomDisjointPaths(t *testing.T) {
	input := "2\n##start\nstart 0 0\n##end\nend 5 0\na1 1 0\na2 2 0\nb1 1 1\nstart-a1\na1-a2\na2-end\nstart-b1\nb1-end\n"
	farm := mustParse(t, input)

	paths := BFS(farm)
	if len(paths) != 2 {
		t.Fatalf("len(paths) = %d, want 2", len(paths))
	}

	seen := map[string]int{}
	for _, p := range paths {
		for _, r := range p.rooms {
			if r == farm.start || r == farm.end {
				continue
			}
			seen[r.name]++
			if seen[r.name] > 1 {
				t.Errorf("room %q used by more than one path", r.name)
			}
		}
	}

	// Shortest path first.
	if paths[0].RoomsNum > paths[1].RoomsNum {
		t.Errorf("paths not sorted by length: %d before %d", paths[0].RoomsNum, paths[1].RoomsNum)
	}
}

var moveRe = regexp.MustCompile(`^L(\d+)-(\S+)$`)

// checkInvariants replays the recorded turns and verifies that no room
// (other than start/end) ever holds more than one ant at once, no link is
// crossed twice in the same turn, and every ant reaches the end room.
func checkInvariants(t *testing.T, farm *AntFarm, turns []string, antsNum int) {
	t.Helper()

	links := map[[2]string]bool{}
	for _, r := range farm.rooms {
		for _, n := range r.links {
			links[[2]string{r.name, n.name}] = true
		}
	}

	position := map[int]string{}
	for i := 1; i <= antsNum; i++ {
		position[i] = farm.start.name
	}

	for turnIdx, line := range turns {
		roomsEntered := map[string]bool{}
		tunnelsUsed := map[[2]string]bool{}
		for _, tok := range strings.Fields(line) {
			m := moveRe.FindStringSubmatch(tok)
			if m == nil {
				t.Fatalf("turn %d: malformed move token %q", turnIdx+1, tok)
			}
			ant, _ := strconv.Atoi(m[1])
			room := m[2]
			prev := position[ant]

			if prev != room && !links[[2]string{prev, room}] {
				t.Errorf("turn %d: ant %d moved %s->%s but no such link", turnIdx+1, ant, prev, room)
			}
			tKey := [2]string{prev, room}
			if tunnelsUsed[tKey] {
				t.Errorf("turn %d: tunnel %s-%s used twice", turnIdx+1, prev, room)
			}
			tunnelsUsed[tKey] = true

			if room != farm.end.name {
				if roomsEntered[room] {
					t.Errorf("turn %d: room %q occupied by more than one ant", turnIdx+1, room)
				}
				roomsEntered[room] = true
			}
			position[ant] = room
		}
	}

	for i := 1; i <= antsNum; i++ {
		if position[i] != farm.end.name {
			t.Errorf("ant %d never reached the end room (last seen at %q)", i, position[i])
		}
	}
}

func TestEndToEnd_Examples(t *testing.T) {
	cases := []struct {
		file     string
		maxTurns int
	}{
		{"example00.txt", 6},
		{"example01.txt", 8},
		{"example02.txt", 11},
		{"example03.txt", 6},
		{"example04.txt", 6},
		{"example05.txt", 8},
	}

	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			f, err := os.Open("../examples/" + tc.file)
			if err != nil {
				t.Fatalf("open %s: %v", tc.file, err)
			}
			defer f.Close()

			farm, err := parseAntFarm(f)
			if err != nil {
				t.Fatalf("parse %s: %v", tc.file, err)
			}

			paths := BFS(farm)
			if len(paths) == 0 {
				t.Fatalf("no paths found for %s", tc.file)
			}

			turns := AntsMovement(paths, farm)
			if len(turns) > tc.maxTurns {
				t.Errorf("%s took %d turns, want <= %d", tc.file, len(turns), tc.maxTurns)
			}

			checkInvariants(t, farm, turns, farm.AntsNum)
		})
	}
}

func TestBadExamples_Reject(t *testing.T) {
	for _, file := range []string{"badexample00.txt", "badexample01.txt"} {
		t.Run(file, func(t *testing.T) {
			f, err := os.Open("../examples/" + file)
			if err != nil {
				t.Fatalf("open %s: %v", file, err)
			}
			defer f.Close()

			if _, err := parseAntFarm(f); err == nil {
				t.Errorf("expected %s to be rejected, but it parsed successfully", file)
			}
		})
	}
}
