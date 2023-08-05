package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var countMergeRequestsCmd = &cobra.Command{
    Use:   "count-merge-requests",
    Short: "Count the merge requests created by a user between two dates",
    Run: func(cmd *cobra.Command, args []string) {
        // Read token
        token := os.Getenv("GITLAB_TOKEN")
        if token == "" {
            fmt.Println("GITLAB_TOKEN environment variable is not set.")
            return
        }

        // Prompt for username and dates
        var username, startDate, endDate string
        fmt.Print("Enter GitLab username: ")
        fmt.Scanln(&username)
        fmt.Print("Enter start date (YYYY-MM-DD): ")
        fmt.Scanln(&startDate)
        fmt.Print("Enter end date (YYYY-MM-DD): ")
        fmt.Scanln(&endDate)

        // Convert dates to time.Time
        start, err := time.Parse("2006-01-02", startDate)
        if err != nil {
            fmt.Println("Error parsing start date:", err)
            return
        }
        end, err := time.Parse("2006-01-02", endDate)
        if err != nil {
            fmt.Println("Error parsing end date:", err)
            return
        }

        // Send request to GitLab API
        url := fmt.Sprintf("%s/merge_requests?scope=all&state=all&author_username=%s&created_after=%s&created_before=%s", "gitlabAPIURL", username, start.Format(time.RFC3339), end.Format(time.RFC3339))
        client := &http.Client{}
        req, err := http.NewRequest("GET", url, nil)
        if err != nil {
            fmt.Printf("Error creating request: %v\n", err)
            return
        }

        // Set the authorization header
        req.Header.Set("Authorization", "Bearer "+token)

        resp, err := client.Do(req)
        if err != nil {
            fmt.Printf("Error making request: %v\n", err)
            return
        }
        defer resp.Body.Close()

        // Get the X-Total header
        xTotal := resp.Header.Get("X-Total")

        // Output the number of merge requests
        fmt.Printf("Number of merge requests created by %s between %s and %s: %s\n", username, startDate, endDate, xTotal)
    },
}

func init() {
    rootCmd.AddCommand(countMergeRequestsCmd)
}
