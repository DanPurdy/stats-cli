package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var countApprovedByUserCmd = &cobra.Command{
    Use:   "count-approved-by-user",
    Short: "Count the merge requests where a specified user has been an approver",
    Run: func(cmd *cobra.Command, args []string) {
        // Read token
        token := os.Getenv("GITLAB_TOKEN")
        if token == "" {
            fmt.Println("GITLAB_TOKEN environment variable is not set.")
            return
        }

        // Prompt for username
        var username string
        fmt.Print("Enter GitLab username of the approver: ")
        fmt.Scanln(&username)

        // Send request to GitLab API
        url := fmt.Sprintf("%s/merge_requests?scope=all&state=merged&approver_usernames=%s", gitlabAPIURL, username)
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
        fmt.Printf("Number of merge requests approved by %s: %s\n", username, xTotal)
    },
}

func init() {
    rootCmd.AddCommand(countApprovedByUserCmd)
}
