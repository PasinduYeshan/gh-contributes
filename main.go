package main

import (
	"log"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

func main() {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatalf("Error creating GraphQL client: %v", err)
	}

	to := time.Now().UTC()
	from := to.AddDate(0, -2, 0)
	from = startOfWeek(from)

	queryParams := map[string]interface{}{
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	// Query the GitHub API for the weekly contributions.
	var resp WeeksContributionsResponse
	err = client.Do(getWeeksDataQuery, queryParams, &resp)
	if err != nil {
		log.Fatalf("GraphQL query failed: %v", err)
	}

	cc := resp.Viewer.ContributionsCollection.ContributionCalendar
	var weeklyContributions []WeekData = processWeeklyContributions(cc)

	renderWeeklyGraphWithDailyStat(weeklyContributions, cc.Weeks)
}
