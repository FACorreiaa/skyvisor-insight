package trips

import (
	"sort"
	"time"

	"github.com/FACorreiaa/Aviation-tracker/app/apiclient"
)

// TimelineEntry is one flight segment placed on the chronological timeline.
type TimelineEntry struct {
	TripID   string
	TripName string
	Segment  apiclient.TripSegment
}

// TimelineDay groups the entries that depart on the same local calendar day.
type TimelineDay struct {
	Date    time.Time
	Entries []TimelineEntry
}

// BuildTimeline flattens every trip's segments into a chronological, day-grouped
// timeline. Segments without a departure time are collected separately so they
// are not silently dropped. Times are shown in the viewer's chosen location.
func BuildTimeline(trips []apiclient.Trip, loc *time.Location) (days []TimelineDay, undated []TimelineEntry) {
	if loc == nil {
		loc = time.UTC
	}
	var dated []TimelineEntry
	for _, trip := range trips {
		for _, segment := range trip.Segments {
			entry := TimelineEntry{TripID: trip.ID, TripName: trip.Name, Segment: segment}
			if segment.DepartsAt.IsZero() {
				undated = append(undated, entry)
				continue
			}
			dated = append(dated, entry)
		}
	}
	sort.SliceStable(dated, func(i, j int) bool {
		return dated[i].Segment.DepartsAt.Before(dated[j].Segment.DepartsAt)
	})

	byDay := map[string]int{}
	for _, entry := range dated {
		local := entry.Segment.DepartsAt.In(loc)
		key := local.Format("2006-01-02")
		index, ok := byDay[key]
		if !ok {
			index = len(days)
			byDay[key] = index
			days = append(days, TimelineDay{Date: time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)})
		}
		days[index].Entries = append(days[index].Entries, entry)
	}
	return days, undated
}

func timelineDayLabel(day TimelineDay) string {
	return day.Date.Format("Monday, 2 January 2006")
}

func timelineSegmentTime(entry TimelineEntry, loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	segment := entry.Segment
	if segment.DepartsAt.IsZero() {
		return "Time not set"
	}
	departs := segment.DepartsAt.In(loc).Format("15:04")
	if segment.ArrivesAt.IsZero() {
		return departs
	}
	return departs + " – " + segment.ArrivesAt.In(loc).Format("15:04")
}
