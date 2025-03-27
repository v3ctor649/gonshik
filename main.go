package main

import (
	"fmt"

	"github.com/spf13/cobra"
	bruter "github.com/v3ctor649/gonshik/scripts"
)

var (
	usernameDict string
	passwordDict string
	url          string
	rate         int
	domain       string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "mycli",
		Short: "A simple CLI tool",
		Run: func(cmd *cobra.Command, args []string) {
			bruter.Brute(&usernameDict, &passwordDict, &url, &domain, &rate)
		},
	}

	rootCmd.Flags().StringVarP(&usernameDict, "username", "u", "", "Username (required)")
	rootCmd.Flags().StringVarP(&passwordDict, "password", "p", "", "Password (required)")
	rootCmd.Flags().StringVarP(&url, "url", "", "", "URL (required)")
	rootCmd.Flags().IntVar(&rate, "rate", 1, "Number of concurrent requests")
	rootCmd.Flags().StringVarP(&domain, "domain", "d", "", "Number of concurrent requests")

	// Mark required flags
	rootCmd.MarkFlagRequired("username")
	rootCmd.MarkFlagRequired("password")
	rootCmd.MarkFlagRequired("url")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
