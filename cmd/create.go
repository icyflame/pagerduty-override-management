/*
Copyright © 2022 Siddharth Kannan <mail@siddharthkannan.in>

*/
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/icyflame/pagerduty-override-management/internal/pagerduty"
	"github.com/spf13/cobra"
)

func init() {
	var from, to string
	var userID, scheduleID string
	var shiftLength time.Duration
	var shiftDays, gapDays int
	var dryRun bool

	// createCmd represents the create command
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create overrides in the given schedule",
		Long: `This command is used to create overrides in the specified schedule.

Overrides are created using the following rules:

1. Each override is of length "shift-length"
2. Each override will be created after the "from" parameter
3. All sets of overrides except the last set will be create before the "to" parameter
4. Each set of overrides will be "shift-days" number of days long
5. The gap between sets of overrides will be "gap-days" long
6. The first shift will start from the time provided at "from"
`,
		Run: func(cmd *cobra.Command, args []string) {
			fromTime, err := time.Parse(time.RFC3339, from)
			if err != nil {
				fmt.Printf("could not parse 'from' time: %s; error: %v", from, err)
				return
			}
			toTime, err := time.Parse(time.RFC3339, to)
			if err != nil {
				fmt.Printf("could not parse 'to' time: %s; error: %v", to, err)
				return
			}

			currentTime := fromTime
			var overridesToCreate []pagerduty.ShiftTimings
			for currentTime.Before(toTime) {
				for i := 0; i < shiftDays; i++ {
					shiftTimeFrom := currentTime
					shiftTimeTo := currentTime.Add(shiftLength)
					overridesToCreate = append(overridesToCreate, pagerduty.ShiftTimings{
						From: shiftTimeFrom,
						To:   shiftTimeTo,
					})

					currentTime = currentTime.Add(24 * time.Hour)
				}

				currentTime = currentTime.Add(time.Duration(gapDays) * 24 * time.Hour)
			}

			createdShifts := pagerduty.CreateOverrides(overridesToCreate, userID, scheduleID, dryRun)
			for i, cs := range createdShifts {
				output := []string{
					overridesToCreate[i].From.Format(time.RFC3339),
					" -> ",
					overridesToCreate[i].To.Format(time.RFC3339),
				}
				if cs.Error != nil {
					output = append(output, "ERROR", cs.Error.Error())
				} else {
					if dryRun {
						output = append(output, "DRY-RUN")
					} else {
						output = append(output, cs.Shift.ID, cs.Shift.User.Summary)
					}
				}
				fmt.Println(strings.Join(output, "\t"))
			}
		},
	}

	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&userID, "user-id", "", "", "User who will take these overrides")
	createCmd.Flags().StringVarP(&scheduleID, "schedule-id", "", "", "Schedule where the overrides should be created")
	createCmd.Flags().StringVarP(&from, "from", "", "", "From time to start creating overrides")
	createCmd.Flags().StringVarP(&to, "to", "", "", "Maximum time until which the overrides might be created")
	createCmd.Flags().DurationVarP(&shiftLength, "shift-length", "", time.Hour, "Length of each shift")
	createCmd.Flags().IntVarP(&shiftDays, "shift-days", "", 0, "Days of each shift")
	createCmd.Flags().IntVarP(&gapDays, "gap-days", "", 0, "Gap between the set of days for each shift")
	createCmd.MarkFlagsRequiredTogether("from", "to", "shift-length", "shift-days", "gap-days", "user-id", "schedule-id")
	createCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Run this script without creating anything on PagerDuty")
}
