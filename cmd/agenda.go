package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/garaemon/org-agenda-cli/pkg/agenda"
	"github.com/garaemon/org-agenda-cli/pkg/config"
	"github.com/garaemon/org-agenda-cli/pkg/item"
	"github.com/garaemon/org-agenda-cli/pkg/parser"
	"github.com/garaemon/org-agenda-cli/pkg/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	agendaRange string
	agendaDate  string
	agendaTag   string
	agendaTui   bool
)

// agendaCmd represents the agenda command
var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "Displays the agenda view",
	Long:  `Displays the agenda view. Aggregates tasks with schedules and deadlines within a specified period.`,
	Run: func(cmd *cobra.Command, args []string) {
		var start time.Time
		var err error

		if agendaDate != "" {
			start, err = time.Parse("2006-01-02", agendaDate)
			if err != nil {
				fmt.Printf("Invalid date format: %v. Use YYYY-MM-DD.\n", agendaDate)
				return
			}
		} else {
			start = time.Now()
		}

		// Truncate time part for date-only comparison
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end := start
		if agendaRange == "week" {
			end = start.AddDate(0, 0, 7)
		}

		paths := viper.GetStringSlice("org_files")
		if len(paths) == 0 {
			if _, err := os.Stat("sample.org"); err == nil {
				paths = []string{"sample.org"}
			} else {
				fmt.Println("No org files configured.")
				return
			}
		}

		orgFiles := config.ResolveOrgFiles(paths)
		if !agendaTui {
			fmt.Printf("Agenda for %s to %s:\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
		}

		var allItems []*item.Item
		for _, file := range orgFiles {
			content, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			items := parser.ParseString(string(content), file)
			if agendaTui {
				allItems = append(allItems, items...)
			} else {
				filtered := agenda.FilterItemsByRange(items, start, end)
				allItems = append(allItems, filtered...)
			}
		}

		if agendaTui {
			err := tui.Run(allItems, start, agendaRange, "")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return
		}

		for _, item := range allItems {
			dateStr := ""
			// Helper to check if a date is strictly within [start, end]
			inRange := func(t *time.Time) bool {
				if t == nil {
					return false
				}
				// Normalize to UTC midnight for date comparison
				d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
				s := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
				e := time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)

				return (d.After(s) || d.Equal(s)) && (d.Before(e) || d.Equal(e))
			}

			if inRange(item.Scheduled) {
				dateStr = fmt.Sprintf("Sched: %s", item.Scheduled.Format("2006-01-02"))
			} else if inRange(item.Deadline) {
				dateStr = fmt.Sprintf("Dead:  %s", item.Deadline.Format("2006-01-02"))
			}
			fmt.Printf("%s: [%s] %s (%s:%d)\n", dateStr, item.Status, item.Title, item.FilePath, item.LineNumber)
		}
	},
}

func init() {
	rootCmd.AddCommand(agendaCmd)

	agendaCmd.Flags().StringVar(&agendaRange, "range", "day", "Specify the display range (day|week)")
	agendaCmd.Flags().StringVar(&agendaDate, "date", "", "Specify the reference date (YYYY-MM-DD, default: today)")
	agendaCmd.Flags().StringVar(&agendaTag, "tag", "", "Filter items by a specific tag")
	agendaCmd.Flags().BoolVar(&agendaTui, "tui", false, "Enable interactive TUI mode")
}
