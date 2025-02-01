package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

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

	// Set the time range to the last 3 months.
	to := time.Now().UTC()
	from := to.AddDate(0, -3, 0)

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

	renderTermuiBarChart(weeklyContributions)
	fmt.Printf("Total weeks: %d\n", len(weeklyContributions))
}

func renderTermuiBarChart(weeklyContributions []WeekData) {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Create a new bar chart
	bc := widgets.NewBarChart()
	bc.Title = "GitHub Weekly Contributions (Last 3 Months)"
	bc.BarWidth = 4
	bc.BarGap = 1
	bc.BarColors = []ui.Color{ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack, ui.ColorClear, ui.ModifierBold)}

	// Prepare data for the bar chart
	labels := make([]string, len(weeklyContributions))
	data := make([]float64, len(weeklyContributions))
	lastMonth := ""

	for i, w := range weeklyContributions {
		date, err := time.Parse("2006-01-02", w.StartDate)
		if err != nil {
			log.Fatalf("Failed to parse date: %v", err)
		}

		currentMonth := date.Format("Jan")
		if currentMonth != lastMonth {
			labels[i] = currentMonth
			lastMonth = currentMonth
		} else {
			labels[i] = ""
		}

		data[i] = float64(w.Count)
		log.Printf("%s: %d", w.StartDate, float64(w.Count))
	}

	bc.Data = data
	bc.Labels = labels

	_, termHeight := ui.TerminalDimensions()
	bc.SetRect(2, 2, 5*19, termHeight/2)

	ui.Render(bc)

	// Wait for a key press to exit.
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "<Resize>":
			payload := e.Payload.(ui.Resize)
			bc.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(bc)
		}
	}
}
