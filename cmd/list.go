/*
Copyright © 2022 Siddharth Kannan <mail@siddharthkannan.in>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/icyflame/pagerduty-override-management/internal/pagerduty"
	"github.com/spf13/cobra"
)

func init() {
	var from, to, scheduleID string

	// listCmd represents the list command
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
			fmt.Println(pagerduty.ListOverrides(fromTime, toTime, scheduleID))
		},
	}

	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&from, "from", "", "", "Oldest time for searching overrides")
	listCmd.Flags().StringVarP(&to, "to", "", "", "Newest time for searching overrides")
	listCmd.Flags().StringVarP(&scheduleID, "schedule-id", "", "", "Schedule ID in which to look for overrides")
	listCmd.MarkFlagsRequiredTogether("from", "to", "schedule-id")
}
