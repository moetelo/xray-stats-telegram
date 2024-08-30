package stats

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type StatsParser struct {
	TrafficDataDirectory string
}

func New(trafficDataDirectory string) *StatsParser {
	return &StatsParser{
		TrafficDataDirectory: trafficDataDirectory,
	}
}

type Stats struct {
	Down int
	Up   int
}

func (s Stats) ToString() string {
	return fmt.Sprintf("↑ %d (mb)\n↓ %d (mb)", s.Up/1024/1024, s.Down/1024/1024)
}

func (s Stats) ToOneLineString() string {
	return fmt.Sprintf("↓ %d (mb) ↑ %d (mb)", s.Down/1024/1024, s.Up/1024/1024)
}

func (parser *StatsParser) sumLines(file string) int {
	sum := 0

	lines := strings.Split(file, "\n")
	for _, line := range lines {
		val, _ := strconv.Atoi(string(line))
		sum += val
	}

	return sum
}

func (parser *StatsParser) GetToday(user string) *Stats {
	today := time.Now().Format(time.DateOnly)

	path := fmt.Sprintf("%s/%s/down/%s", parser.TrafficDataDirectory, user, today)

	var stats Stats

	file, err := os.ReadFile(path)
	if err == nil {
		stats.Down = parser.sumLines(string(file))
	}

	fmt.Println(path, err)

	path = fmt.Sprintf("%s/%s/up/%s", parser.TrafficDataDirectory, user, today)
	file, err = os.ReadFile(path)
	if err == nil {
		stats.Up = parser.sumLines(string(file))
	}

	fmt.Println(path, err)

	return &stats
}
