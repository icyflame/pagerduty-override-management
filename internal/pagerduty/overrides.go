package pagerduty

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

var authToken string = os.Getenv("AUTHORIZATION_TOKEN")

// ListOverrides ...
func ListOverrides(from, to time.Time, scheduleID string) []pagerduty.Override {
	cli := pagerduty.NewClient(authToken)
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
