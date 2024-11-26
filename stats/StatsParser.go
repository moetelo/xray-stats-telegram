package stats

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type StatsParser struct {
	statsQueryBin string

	dateToUserToStats      map[string]map[string]*Stats
	lastQueryExecutionTime time.Time

	mu sync.Mutex
}

func New(statsQueryBin string) *StatsParser {
	if statsQueryBin == "" {
		statsQueryBin = "/usr/local/bin/stats-query"
	}

	return &StatsParser{
		statsQueryBin:          statsQueryBin,
		dateToUserToStats:      make(map[string]map[string]*Stats),
		lastQueryExecutionTime: time.Time{},
	}
}

func (parser *StatsParser) Query(user string, byDate *time.Time) *Stats {
	dateOnly := byDate.Format(time.DateOnly)

	if parser.shouldInvalidateCache(dateOnly) {
		parser.invalidateCache(dateOnly)
	}

	cachedUserStats := parser.getOrInitCachedStats(dateOnly, user)
	if cachedUserStats != nil {
		return cachedUserStats
	}

	cmd := exec.Command(parser.statsQueryBin, "--plain", dateOnly)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return nil
	}

	parser.mu.Lock()
	parser.lastQueryExecutionTime = time.Now()
	parser.mu.Unlock()

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

			parser.mu.Lock()
			parser.dateToUserToStats[dateOnly][username] = stats
			parser.mu.Unlock()

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

func (parser *StatsParser) getOrInitCachedStats(dateOnly string, user string) *Stats {
	parser.mu.Lock()
	defer parser.mu.Unlock()
	if userCache, exists := parser.dateToUserToStats[dateOnly]; exists {
		if stats, exists := userCache[user]; exists {
			return stats
		}
	} else {
		parser.dateToUserToStats[dateOnly] = make(map[string]*Stats)
	}
	return nil
}

func (parser *StatsParser) shouldInvalidateCache(dateOnly string) bool {
	today := time.Now().Format(time.DateOnly)
	if dateOnly != today {
		return false
	}

	return time.Since(parser.lastQueryExecutionTime) > 5*time.Minute
}

func (parser *StatsParser) invalidateCache(dateOnly string) {
	parser.mu.Lock()
	defer parser.mu.Unlock()
	delete(parser.dateToUserToStats, dateOnly)
}
