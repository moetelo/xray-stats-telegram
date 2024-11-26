package stats

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type StatsParser struct {
	statsQueryBin string
}

func New(statsQueryBin string) *StatsParser {
	if statsQueryBin == "" {
		statsQueryBin = "/usr/local/bin/stats-query"
	}

	return &StatsParser{
		statsQueryBin: statsQueryBin,
	}
}

func (p *StatsParser) Query(byDate time.Time) []Stats {
	result := make([]Stats, 0, 5)

	dateOnly := byDate.Format(time.DateOnly)

	cmd := exec.Command(p.statsQueryBin, "--plain", dateOnly)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return result
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		stats := p.parseStatsLine(line)
		if stats == nil {
			continue
		}

		result = append(result, *stats)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return result
	}

	return result
}

func (p *StatsParser) QueryUser(user string, byDate time.Time) *Stats {
	dateOnly := byDate.Format(time.DateOnly)

	cmd := exec.Command(p.statsQueryBin, "--plain", dateOnly, user)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil
	}

	return p.parseStatsLine(string(output))
}

func (p *StatsParser) parseStatsLine(line string) *Stats {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		fmt.Println("Unexpected output:", line)
		return nil
	}

	username := fields[0]
	down, _ := strconv.Atoi(fields[1])
	up, _ := strconv.Atoi(fields[2])
	stats := &Stats{
		User: username,
		Down: down,
		Up:   up,
	}

	return stats
}
