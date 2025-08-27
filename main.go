package main

import (
	"fmt"
	"os"

	"github.com/antrusd/gke-gcloud-auth-plugin/pkg/cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("failure during execution: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
