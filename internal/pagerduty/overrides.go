package pagerduty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type User struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Summary string `json:"summary"`
}

type Override struct {
	From time.Time `json:"start"`
	To   time.Time `json:"end"`
	User User      `json:"user"`
}

type ListOverridesResponse struct {
	Total     int        `json:"total"`
	Overrides []Override `json:"overrides"`
}

var authToken string = os.Getenv("AUTHORIZATION_TOKEN")

// ListOverrides ...
func ListOverrides(from, to time.Time, scheduleID string) []Override {
	endpoint := fmt.Sprintf("https://api.pagerduty.com/schedules/%s/overrides", scheduleID)

	urlParams := fmt.Sprintf("since=%s&until=%s", from.Local().Format(time.RFC3339), to.Local().Format(time.RFC3339))
	endpointWithParams := fmt.Sprintf("%s?%s", endpoint, urlParams)

	// urlParams := "since=2022-08-17T10:00:00+09:00&until=2022-08-31T10:00:00+09:00"
	// endpointWithParams := endpoint + "?" + urlParams

	fmt.Println(endpointWithParams)

	req, err := http.NewRequest("GET", endpointWithParams, nil)
	if err != nil {
		fmt.Printf("could not create GET request: %v", err)
		return nil
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Add("Authorization", fmt.Sprintf("Token token=%s", authToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("could not perform Get request: %v", err)
		return nil
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("could not read response body: %v", err)
		return nil
	}

	fmt.Println("full response: ", string(body))

	var response ListOverridesResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Printf("error while unmarshaling pagerduty response: %v", err)
		return nil
	}

	return response.Overrides
}
