package stats

import (
	"bufio"
	"fmt"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"xray-stats-telegram/queryDate"
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

func (p *StatsParser) Query(queryDate queryDate.QueryDate) []Stats {
	result := make([]Stats, 0, 5)

	cmd := exec.Command(p.statsQueryBin, "--plain", queryDate.String())
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

	sort.Slice(result, func(i, j int) bool {
		return result[i].DownBytes < result[j].DownBytes
	})

	return result
}

func (p *StatsParser) QueryUser(user string, queryDate queryDate.QueryDate) *Stats {
	cmd := exec.Command(p.statsQueryBin, "--plain", queryDate.String(), user)
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
		UserEmail: username,
		DownBytes: down,
		UpBytes:   up,
	}

	return stats
}
