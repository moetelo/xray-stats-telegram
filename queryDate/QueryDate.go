package queryDate

import (
	"fmt"
	"strconv"
	"time"
)

type QueryDate struct {
	year  int
	month time.Month
	day   *int
}

var fullIsoLength = len("2024-01-01")
var yearMonthIsoLength = len("2024-01")

func FromTime(t time.Time) QueryDate {
	day := t.Day()
	return QueryDate{
		year:  t.Year(),
		month: t.Month(),
		day:   &day,
	}
}

func Now() QueryDate {
	return FromTime(time.Now())
}

func (d QueryDate) String() string {
	if d.day != nil {
		return fmt.Sprintf("%d-%02d-%02d", d.year, d.month, *d.day)
	}

	return fmt.Sprintf("%d-%02d", d.year, d.month)
}

func (d QueryDate) time() time.Time {
	if d.day != nil {
		return time.Date(d.year, d.month, *d.day, 0, 0, 0, 0, time.UTC)
	}

	return time.Date(d.year, d.month, 1, 0, 0, 0, 0, time.UTC)
}

func (date QueryDate) Prev() QueryDate {
	t := date.time()
	if date.day != nil {
		return FromTime(t.AddDate(0, 0, -1))
	}

	newTime := FromTime(t.AddDate(0, -1, 0))
	newTime.day = nil
	return newTime
}

func (date QueryDate) Next() QueryDate {
	t := date.time()
	if date.day != nil {
		return FromTime(t.AddDate(0, 0, +1))
	}

	newTime := FromTime(t.AddDate(0, +1, 0))
	newTime.day = nil
	return newTime
}

func Parse(possiblyDate string) (QueryDate, error) {
	if len(possiblyDate) != fullIsoLength && len(possiblyDate) != yearMonthIsoLength {
		return QueryDate{}, fmt.Errorf("invalid date format")
	}

	year, err := strconv.Atoi(possiblyDate[:len("2024")])
	if err != nil {
		return QueryDate{}, err
	}

	month, err := strconv.Atoi(possiblyDate[len("2024-"):7])
	if err != nil {
		return QueryDate{}, err
	}

	if len(possiblyDate) == yearMonthIsoLength {
		return QueryDate{
			year:  year,
			month: time.Month(month),
		}, nil
	}

	day, err := strconv.Atoi(possiblyDate[len("2024-01-"):])
	if err != nil {
		return QueryDate{}, err
	}

	return QueryDate{
		year:  year,
		month: time.Month(month),
		day:   &day,
	}, nil
}
