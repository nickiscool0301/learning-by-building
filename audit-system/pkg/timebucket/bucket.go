package timebucket

import "time"

type Granularity string

const (
	Daily  Granularity = "daily"
	Hourly Granularity = "hourly"
)

// take timestamp, return its bucket key.
func Compute(t time.Time, g Granularity) string {
	switch g {
	case Daily:
		return t.UTC().Format("2006-01-02")
	case Hourly:
		return t.UTC().Format("2006-01-02-15")
	}
	panic("invalid granularity")
}

func Range(start, end time.Time, g Granularity) []string {
	switch g {
	case Daily:
		return rangeDaily(start, end)
	case Hourly:
		return rangeHourly(start, end)
	}
	panic("invalid granularity")
}

func rangeDaily(start, end time.Time) []string {
	var buckets []string
	for start.Before(end) {
		buckets = append(buckets, start.UTC().Format("2006-01-02"))
		start = start.AddDate(0, 0, 1)
	}
	return buckets
}

func rangeHourly(start, end time.Time) []string {
	var buckets []string
	for start.Before(end) {
		buckets = append(buckets, start.UTC().Format("2006-01-02-15"))
		start = start.Add(time.Hour)
	}
	return buckets
}
