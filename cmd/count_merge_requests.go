package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)

var countMergeRequestsCmd = &cobra.Command{
    Use:   "count-merge-requests",
    Short: "Count the merge requests created by a user between two dates",
    Run: func(cmd *cobra.Command, args []string) {
        // Read token
        token := os.Getenv("GITLAB_CLI_TOKEN")
        if token == "" {
            fmt.Println("GITLAB_CLI_TOKEN environment variable is not set.")
            return
        }

        git, err := gitlab.NewClient(token)
        if err != nil {
            log.Fatalf("Failed to create client: %v", err)
        }

        // Prompt for username and dates
        var username string
        fmt.Print("Enter GitLab username: ")
        fmt.Scanln(&username)
        // fmt.Print("Enter start date (YYYY-MM-DD): ")
        // fmt.Scanln(&startDate)
        // fmt.Print("Enter end date (YYYY-MM-DD): ")
        // fmt.Scanln(&endDate)

        userOpt:= gitlab.SearchOptions{ListOptions: gitlab.ListOptions{PerPage: 1}};

        users, _, err := git.Search.Users(username, &userOpt);
        if err != nil {
            log.Fatalf("Failed to get user: %v", err)
        }

        // Convert dates to time.Time
        start := time.Now().Add(-time.Hour * 720);
        if err != nil {
            fmt.Println("Error parsing start date:", err)
            return
        }
        end := time.Now()
        if err != nil {
            fmt.Println("Error parsing end date:", err)
            return
        }

        opt:= gitlab.ListMergeRequestsOptions{
            ListOptions: gitlab.ListOptions{PerPage: 1},
            AuthorID: gitlab.Int(users[0].ID),
            CreatedAfter: gitlab.Time(start), 
            CreatedBefore: gitlab.Time(end),
            Scope: gitlab.String("all"),
        }

        _, resp2, err := git.MergeRequests.ListMergeRequests(&opt)
        if err != nil {
            log.Fatalf("Failed to list merge requests: %v", err)
        }

        xTotal:= resp2.Header.Get("X-Total");

        // Output the number of merge requests
        fmt.Printf("Number of merge requests created by %s between %s and %s: %s\n", username, start, end, xTotal)
    },
}

func init() {
    rootCmd.AddCommand(countMergeRequestsCmd)
}
