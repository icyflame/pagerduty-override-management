package pagerduty

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

type ShiftTimings struct {
	From, To time.Time
}

type CreatedShifts struct {
	Shift *pagerduty.Override
	Error error
}

var authToken string = os.Getenv("AUTHORIZATION_TOKEN")
var cli *pagerduty.Client = pagerduty.NewClient(authToken)

// ListOverrides ...
func ListOverrides(from, to time.Time, scheduleID string) []pagerduty.Override {
	res, err := cli.ListOverridesWithContext(context.Background(), scheduleID, pagerduty.ListOverridesOptions{
		Since: from.Format(time.RFC3339),
		Until: to.Format(time.RFC3339),
	})
	if err != nil {
		fmt.Printf("could not get from pagerduty: %v", err)
		return nil
	}

	return res.Overrides
}

func CreateOverrides(shifts []ShiftTimings, userID, scheduleID string, dryRun bool) []CreatedShifts {
	output := make([]CreatedShifts, len(shifts))
	for i, s := range shifts {
		overrideToCreate := pagerduty.Override{
			Start: s.From.Format(time.RFC3339),
			End:   s.To.Format(time.RFC3339),
			User: pagerduty.APIObject{
				ID:   userID,
				Type: "user_reference",
			},
		}

		if dryRun {
			output[i].Shift = &overrideToCreate
			continue
		}

		override, err := cli.CreateOverride(scheduleID, overrideToCreate)
		if err != nil {
			output[i].Error = err
		} else {
			output[i].Shift = override
		}
	}
	return output
}

// DeleteOverride ...
func DeleteOverride(scheduleID, overrideID string) error {
	return cli.DeleteOverrideWithContext(context.Background(), scheduleID, overrideID)
}
