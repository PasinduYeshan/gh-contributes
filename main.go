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

	to := time.Now().UTC()
	from := to.AddDate(0, -2, 0)
	from = startOfWeek(from)

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
			// Calculate the start of the week (Sunday)
			startOfWeek := startOfWeek(startDate)
			weeklyContributions = append(weeklyContributions, WeekData{
				StartDate: startOfWeek.Format("2006-01-02"),
				Count:     weekCount,
			})
		}
	}

	renderTermui(weeklyContributions, cc.Weeks)
}

func startOfWeek(date time.Time) time.Time {
	return date.AddDate(0, 0, -int(date.Weekday()))
}

func renderTermui(weeklyContributions []WeekData, weeks []struct {
	ContributionDays []struct {
		Date              string `json:"date"`
		ContributionCount int    `json:"contributionCount"`
	} `json:"contributionDays"`
}) {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer ui.Close()

	// Bar chart.
	bc := widgets.NewBarChart()
	bc.Title = "GitHub Weekly Contributions (Last 2 Months)"
	bc.BarWidth = 6
	bc.BarGap = 2
	bc.BarColors = []ui.Color{ui.ColorGreen}
	bc.LabelStyles = []ui.Style{ui.NewStyle(ui.ColorWhite, ui.ColorClear, ui.ModifierBold)}
	bc.NumStyles = []ui.Style{ui.NewStyle(ui.ColorBlack, ui.ColorClear, ui.ModifierBold)}

	labels := make([]string, len(weeklyContributions))
	data := make([]float64, len(weeklyContributions))

	for i, w := range weeklyContributions {
		date, err := time.Parse("2006-01-02", w.StartDate)
		if err != nil {
			labels[i] = fmt.Sprintf("Week %d", i+1)
		} else {
			labels[i] = date.Format("01/02")
		}
		data[i] = float64(w.Count)
	}

	bc.Data = data
	bc.Labels = labels

	// Stat text.
	p := widgets.NewParagraph()
	p.Title = "Last 5 Days Contributions"
	p.Text = generateLast5DaysText(weeks)
	p.SetRect(0, 0, 50, 8)

	// Grid.
	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(0.7, bc),
		ui.NewRow(0.3, p),
	)

	ui.Render(grid)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			// On 'q' or ctrl-c, exit
			return
		case "<Resize>":
			// On resize event, recalc the grid size
			payload := e.Payload.(ui.Resize)
			grid.SetRect(0, 0, payload.Width, payload.Height)
			ui.Clear()
			ui.Render(grid)
		}
	}
}

func generateLast5DaysText(weeks []struct {
	ContributionDays []struct {
		Date              string `json:"date"`
		ContributionCount int    `json:"contributionCount"`
	} `json:"contributionDays"`
}) string {
	var allDays []struct {
		Date              string `json:"date"`
		ContributionCount int    `json:"contributionCount"`
	}
	for _, w := range weeks {
		allDays = append(allDays, w.ContributionDays...)
	}

	for i := 0; i < len(allDays); i++ {
		for j := i + 1; j < len(allDays); j++ {
			dateI, _ := time.Parse("2006-01-02", allDays[i].Date)
			dateJ, _ := time.Parse("2006-01-02", allDays[j].Date)
			if dateI.Before(dateJ) {
				allDays[i], allDays[j] = allDays[j], allDays[i]
			}
		}
	}

	last5Days := allDays[:5]

	text := "\n"
	for _, day := range last5Days {
		text += fmt.Sprintf("  %s: %d\n", day.Date, day.ContributionCount)
	}
	text += "\nPress 'q' to quit"
	return text
}
