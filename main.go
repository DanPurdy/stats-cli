package main

import (
	"fmt"
	"stats-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}