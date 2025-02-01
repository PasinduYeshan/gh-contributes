package main

// WeeksContributionsResponse represents the structure of the GraphQL response.
type WeeksContributionsResponse struct {
	Viewer struct {
		ContributionsCollection struct {
			ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
		} `json:"contributionsCollection"`
	} `json:"viewer"`
}

// ContributionCalendar represents the calendar with weeks.
type ContributionCalendar struct {
	Weeks []Week `json:"weeks"`
}

// Week represents a week of contribution days.
type Week struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

// ContributionDay represents a single day in the calendar.
type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

// WeekData is a simplified structure for a week’s data.
type WeekData struct {
	StartDate string
	Count     int
}
