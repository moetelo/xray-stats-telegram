package stats

import "fmt"

type Stats struct {
	UserEmail string

	DownBytes int
	UpBytes   int
}

func (s Stats) ToString() string {
	return fmt.Sprintf("↑ %.1f MB\n↓ %.1f MB", bytesToMB(s.UpBytes), bytesToMB(s.DownBytes))
}

func (s Stats) ToOneLineString() string {
	return fmt.Sprintf("↓ %.1f MB ↑ %.1f MB", bytesToMB(s.DownBytes), bytesToMB(s.UpBytes))
}

func bytesToMB(bytes int) float64 {
	return float64(bytes) / 1024 / 1024
}
