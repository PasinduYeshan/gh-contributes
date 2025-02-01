package main

import (
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func renderWeeklyGraphWithDailyStat(
	weeklyContributions []WeekData,
	weeks []Week,
) {

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

// generateLast5DaysText aggregates and returns a string for the last 5 days of contributions.
func generateLast5DaysText(weeks []Week) string {

	// Flatten all days from all weeks.
	var allDays []ContributionDay
	for _, w := range weeks {
		allDays = append(allDays, w.ContributionDays...)
	}

	// Sort days in descending order.
	for i := 0; i < len(allDays); i++ {
		for j := i + 1; j < len(allDays); j++ {
			dateI, _ := time.Parse("2006-01-02", allDays[i].Date)
			dateJ, _ := time.Parse("2006-01-02", allDays[j].Date)
			if dateI.Before(dateJ) {
				allDays[i], allDays[j] = allDays[j], allDays[i]
			}
		}
	}

	// Pick the last 5 days.
	if len(allDays) > 5 {
		allDays = allDays[:5]
	}

	text := "\n"
	for _, day := range allDays {
		text += fmt.Sprintf("  %s: %d\n", day.Date, day.ContributionCount)
	}
	text += "\nPress 'q' to quit"
	return text
}
