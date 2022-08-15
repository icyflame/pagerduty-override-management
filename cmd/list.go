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
	var from, to, scheduleID string

	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "List overrides between the from and to timestamp provided to this command as arguments",
		Run: func(cmd *cobra.Command, args []string) {
			fromTime, err := time.Parse("2006-01-02", from)
			if err != nil {
				fmt.Printf("could not parse 'from' time: %s; error: %v", from, err)
				return
			}
			toTime, err := time.Parse("2006-01-02", to)
			if err != nil {
				fmt.Printf("could not parse 'to' time: %s; error: %v", to, err)
				return
			}

			overrides := pagerduty.ListOverrides(fromTime, toTime, scheduleID)
			for _, override := range overrides {
				fmt.Println(strings.Join([]string{
					override.Start,
					" -> ",
					override.End,
					override.ID,
					override.User.Summary,
				}, "\t"))
			}
		},
	}

	listCmd.Flags().StringVarP(&from, "from", "", "", "Start time for listing overrides")
	listCmd.Flags().StringVarP(&to, "to", "", "", "End time for listing overrides")
	listCmd.Flags().StringVarP(&scheduleID, "schedule-id", "", "", "Schedule to look for overrides in")
	listCmd.MarkFlagsRequiredTogether("from", "to", "schedule-id")

	rootCmd.AddCommand(listCmd)
}
