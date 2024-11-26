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

func (parser *StatsParser) Query(user string, byDate time.Time) *Stats {
	dateOnly := byDate.Format(time.DateOnly)

	cmd := exec.Command(parser.statsQueryBin, "--plain", dateOnly)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) == 3 {
			username := fields[0]
			down, _ := strconv.Atoi(fields[1])
			up, _ := strconv.Atoi(fields[2])
			stats := &Stats{
				User: username,
				Down: down,
				Up:   up,
			}

			if username == user {
				return stats
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading command output:", err)
		return nil
	}

	return nil
}
