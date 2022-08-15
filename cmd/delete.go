/*
Copyright © 2022 Siddharth Kannan <mail@siddharthkannan.in>

*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/icyflame/pagerduty-override-management/internal/pagerduty"
	"github.com/spf13/cobra"
)

func init() {
	var filePath, scheduleID string

	// deleteCmd represents the delete command
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			f, err := os.Open(filePath)
			if err != nil {
				fmt.Println("could not read the file with override IDs", err)
				return
			}

			csvReader := csv.NewReader(f)
			for {
				record, err := csvReader.Read()
				if err != nil {
					break
				}

				if len(record) == 0 {
					continue
				}

				overrideID := record[0]
				output := []string{
					overrideID,
				}
				err = pagerduty.DeleteOverride(scheduleID, overrideID)
				if err != nil {
					output = append(output, err.Error())
				} else {
					output = append(output, "DELETED")
				}

				fmt.Println(strings.Join(output, "\t"))
			}
		},
	}

	deleteCmd.Flags().StringVarP(&scheduleID, "schedule-id", "", "", "Schedule to look for overrides in")
	deleteCmd.Flags().StringVarP(&filePath, "file-path", "", "", "File path with the list of override IDs to delete")
	deleteCmd.MarkFlagsRequiredTogether("schedule-id", "file-path")

	rootCmd.AddCommand(deleteCmd)
}
