package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var countApprovedByUserCmd = &cobra.Command{
    Use:   "count-approved-by-user",
    Short: "Count the merge requests where a specified user has been an approver",
    Run: func(cmd *cobra.Command, args []string) {
        // Read token
        token := os.Getenv("GITLAB_CLI_TOKEN")
        if token == "" {
            fmt.Println("GITLAB_CLI_TOKEN environment variable is not set.")
            return
        }

        client, err := gitlab.NewClient(token)
        if err != nil {
            log.Fatalf("Failed to create client: %v", err)
        }

        // Prompt for username and dates
        var username, startDate, endDate string
        fmt.Print("Enter GitLab username: ")
        fmt.Scanln(&username)
        fmt.Print("Enter start date (YYYY-MM-DD): ")
        fmt.Scanln(&startDate)
        fmt.Print("Enter end date (YYYY-MM-DD): ")
        fmt.Scanln(&endDate)


        userOpt:= gitlab.SearchOptions{ListOptions: gitlab.ListOptions{PerPage: 1}};

        users, _, err := client.Search.Users(username, &userOpt);
        if err != nil {
            log.Fatalf("Failed to get user: %v", err)
        }
    

        // Convert dates to time.Time
        startParsed, err := time.Parse("2006-01-02", startDate)
        start := time.Date(startParsed.Year(), startParsed.Month(), startParsed.Day(), 0, 0, 0, 0, time.UTC)
        if err != nil {
            fmt.Println("Error parsing start date:", err)
            return
        }
        endParsed, err := time.Parse("2006-01-02", endDate)
        end := time.Date(endParsed.Year(), endParsed.Month(), endParsed.Day(), 23, 59, 59, 999, time.UTC)
        if err != nil {
            fmt.Println("Error parsing end date:", err)
            return
        }

        approvers:= []int{users[0].ID}
        

        opt:= gitlab.ListMergeRequestsOptions{
            ListOptions: gitlab.ListOptions{PerPage: 1},
            ApprovedByIDs: gitlab.ApproverIDs(approvers),
            CreatedAfter: gitlab.Time(start), 
            CreatedBefore: gitlab.Time(end),
            Scope: gitlab.String("all"),
        }

        _, resp2, err := client.MergeRequests.ListMergeRequests(&opt)
        if err != nil {
            log.Fatalf("Failed to list merge requests: %v", err)
        }

        xTotal:= resp2.Header.Get("X-Total");

        // Output the number of merge requests
        fmt.Printf("Number of merge requests approved by %s: %s\n", username, xTotal)
    },
}

func init() {
    rootCmd.AddCommand(countApprovedByUserCmd)
}
