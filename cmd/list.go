/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
	"gopry/utils"
	"os"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all system packages.",
	Long: `Will list out the system packages currently installed on the OS. Mostly used for 
inventory and fact finding about system dependencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		if sysPkgs, err := utils.GetSystemPkgs(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		} else {
			// create the command line table
			table := tabby.New()
			table.AddHeader("ID", "Name", "Version", "Description")

			for _, pkg := range sysPkgs {
				table.AddLine(pkg.ID, pkg.Name, pkg.Version, pkg.Description)
			}

			table.Print()
		}
	},
}

func init() {
	pkgCmd.AddCommand(listCmd)

	// Local flags
	pkgCmd.Flags().BoolP("all", "a", false, "Does not audit the packages but rather just dumps the list")
}
