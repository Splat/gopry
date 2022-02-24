/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopry",
	Short: "Linux system inspection",
	Long: `Gopry currently allows for the listing and auditing of system packages 
against known CVEs in the wild. Gopry also has the ability to patch specific packages
when specified. 

By default will run an audit of the system packages the same as: "gpry pkg audit"`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root command called")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
