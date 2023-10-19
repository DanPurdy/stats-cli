package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
)



var main = &cobra.Command{
    Use:   "get-group-members",
    Short: "Get a list of gitlab groups and their members",
    Run: func(cmd *cobra.Command, args []string) {
		getGroupMembers();
		
	},
}

func getGroupMembers() () {
	token := os.Getenv("GITLAB_CLI_TOKEN")
	if token == "" {
		fmt.Println("GITLAB_CLI_TOKEN environment variable is not set.")
		return
	}

	client, err := gitlab.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	groups, _, err :=client.Groups.ListGroups(&gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{PerPage: 100},
	});

	groupNames := make([] Group, 0, len(groups))

	for i, value := range groups {
		groupNames = append(groupNames, Group{Name: value.Name, Members: []Member{}});

		members, _, err  := client.Groups.ListGroupMembers(value.ID, &gitlab.ListGroupMembersOptions{ 
			
			ListOptions: gitlab.ListOptions{PerPage: 100}, })
		if err != nil {
			fmt.Print("Failed to get Group Members for group", value.Name, "\n\n");
		} else {
			fmt.Print("GROUP: ", value.Name, ":\n")

			for _, member := range members {
				groupNames[i].Members = append(groupNames[i].Members, Member {Username: member.Username, Name: member.Name})
				fmt.Print(member.Username, ",\n")
			}
			
			fmt.Print("\n\n");
		}

		

	}

	if err != nil {
		fmt.Print("Failed to get Group");
	}

	return;
		
}

func init() {
    rootCmd.AddCommand(main)
}

type Group struct{Name string; Members []Member}
type Member struct{Username string; Name string}