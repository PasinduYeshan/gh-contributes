package main

import (
	"log"
	"time"
)

// startOfWeek returns the Sunday of the week for a given time.
func startOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	return t.AddDate(0, 0, -weekday)
}

// processWeeklyContributions converts the raw ContributionCalendar data into a slice of WeekData.
func processWeeklyContributions(cc ContributionCalendar) []WeekData {

	var weeklyContributions []WeekData
	for _, w := range cc.Weeks {
		weekCount := 0
		var startDate time.Time
		for _, d := range w.ContributionDays {
			weekCount += d.ContributionCount
			date, err := time.Parse("2006-01-02", d.Date)
			if err != nil {
				log.Fatalf("Error parsing date: %v", err)
			}
			if startDate.IsZero() || date.Before(startDate) {
				startDate = date
			}
		}
		if !startDate.IsZero() {
			startOfWk := startOfWeek(startDate)
			weeklyContributions = append(weeklyContributions, WeekData{
				StartDate: startOfWk.Format("2006-01-02"),
				Count:     weekCount,
			})
		}
	}
	return weeklyContributions
}
