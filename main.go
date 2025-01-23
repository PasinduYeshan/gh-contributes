package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
)

// 3. Response struct matching JSON fields.
type contributionsResponse struct {
	Viewer struct {
		ContributionsCollection struct {
			TotalCommitContributions            int `json:"totalCommitContributions"`
			TotalPullRequestContributions       int `json:"totalPullRequestContributions"`
			TotalPullRequestReviewContributions int `json:"totalPullRequestReviewContributions"`
			TotalIssueContributions             int `json:"totalIssueContributions"`
			TotalRepositoryContributions        int `json:"totalRepositoryContributions"`
			RestrictedContributionsCount        int `json:"restrictedContributionsCount"`
			ContributionCalendar                struct {
				TotalContributions int `json:"totalContributions"`
				Weeks              []struct {
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

	// 1 year ago.
	to := time.Now().UTC()
	from := to.AddDate(-1, 0, 0)

	query := `
    query($from: DateTime!, $to: DateTime!) {
      viewer {
        contributionsCollection(from: $from, to: $to) {
          totalCommitContributions
          totalPullRequestContributions
          totalPullRequestReviewContributions
          totalIssueContributions
          totalRepositoryContributions
          restrictedContributionsCount
          contributionCalendar {
            totalContributions
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

	// 4. Make the request with from/to as variables
	variables := map[string]interface{}{
		"from": from.Format(time.RFC3339),
		"to":   to.Format(time.RFC3339),
	}

	var resp contributionsResponse
	err = client.Do(query, variables, &resp)
	if err != nil {
		log.Fatalf("GraphQL query failed: %v", err)
	}

	cc := resp.Viewer.ContributionsCollection

	// 6. Gather ALL daily contributions from the entire year.
	type DayData struct {
		Date  string
		Count int
	}
	var allDays []DayData
	for _, w := range cc.ContributionCalendar.Weeks {
		for _, d := range w.ContributionDays {
			allDays = append(allDays, DayData{
				Date:  d.Date,
				Count: d.ContributionCount,
			})
		}
	}

	// Ensure we have at least 7 days total.
	totalDays := len(allDays)
	if totalDays == 0 {
		fmt.Println("No contributions in the last year!")
		return
	}

	// We'll slice out the last 7 (or fewer if total < 7)
	startIndex := totalDays - 5
	if startIndex < 0 {
		startIndex = 0
	}
	last5Days := allDays[startIndex:]

	fmt.Println("------------------------------------------")
	fmt.Println("👋 Your GitHub Contributions (Last 5 days):")
	for i := len(last5Days) - 1; i >= 0; i-- {
		d := last5Days[i]
		fmt.Printf("  %s: contributions: %d\n", d.Date, d.Count)
	}

	// 5. Print the "full year" stats.
	fmt.Println("------------------------------------------")
	fmt.Println("👋 Your GitHub Contributions (Last Year):")
	fmt.Printf(" • Total Commits:              %d\n", cc.TotalCommitContributions)
	fmt.Printf(" • Total Pull Requests:        %d\n", cc.TotalPullRequestContributions)
	fmt.Printf(" • Total Pull Request Reviews: %d\n", cc.TotalPullRequestReviewContributions)
	fmt.Printf(" • Total Issues:               %d\n", cc.TotalIssueContributions)
	fmt.Printf(" • Total Repositories:         %d\n", cc.TotalRepositoryContributions)
	fmt.Printf(" • Private Contributions:      %d\n", cc.RestrictedContributionsCount)
	fmt.Printf(" • Overall Contributions:      %d\n\n", cc.ContributionCalendar.TotalContributions)
	fmt.Println("------------------------------------------")
}
