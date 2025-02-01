package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

// WeekData represents the weekly contribution data.
type WeekData struct {
	StartDate string
	Count     int
}

type contributionsResponse struct {
	Viewer struct {
		ContributionsCollection struct {
			ContributionCalendar struct {
				Weeks []struct {
					ContributionDays []struct {
						Date              string `json:"date"`
						ContributionCount int    `json:"contributionCount"`
					} `json:"contributionDays"`
				} `json:"weeks"`
			} `json:"contributionCalendar"`
		} `json:"contributionsCollection"`
	} `json:"viewer"`
}

func main() {
	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatalf("Error creating GraphQL client: %v", err)
	}

	// Set the time range to the last 3 months
	to := time.Now().UTC()
	from := to.AddDate(0, -3, 0) // 3 months ago

	query := `
    query($from: DateTime!, $to: DateTime!) {
      viewer {
        contributionsCollection(from: $from, to: $to) {
          contributionCalendar {
            weeks {
              contributionDays {
                date
                contributionCount
              }
            }
          }
        }
      }
    }
    `

	variables := map[string]interface{}{
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	var resp contributionsResponse
	err = client.Do(query, variables, &resp)
	if err != nil {
		log.Fatalf("GraphQL query failed: %v", err)
	}

	cc := resp.Viewer.ContributionsCollection.ContributionCalendar

	// Aggregate contributions by week
	var weeklyContributions []WeekData
	for _, w := range cc.Weeks {
		weekCount := 0
		for _, d := range w.ContributionDays {
			weekCount += d.ContributionCount
		}
		if len(w.ContributionDays) > 0 {
			weeklyContributions = append(weeklyContributions, WeekData{
				StartDate: w.ContributionDays[0].Date,
				Count:     weekCount,
			})
		}
	}

	renderTerminalGraph(weeklyContributions)
}

func renderTerminalGraph(weeklyContributions []WeekData) {
	fmt.Println("GitHub Contributions (Last 3 Months)")
	fmt.Println("------------------------------------")

	maxCount := 0
	for _, w := range weeklyContributions {
		if w.Count > maxCount {
			maxCount = w.Count
		}
	}

	for _, w := range weeklyContributions {
		bar := ""
		for i := 0; i < w.Count; i++ {
			bar += "█"
		}
		fmt.Printf("%s: %3d %s\n", w.StartDate, w.Count, bar)
	}
}
