/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopry/utils"
)

var limit int

// auditCmd represents the audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if sysPkgs, err := utils.GetSystemPkgs(); err != nil {
			return err
		} else {
			// create the command line table
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()

			tbl := table.New("ID", "Name", "Version", "Base Score", "Severity", "Impact Score", "Exploitability Score", "CVE")
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

			for _, pkg := range sysPkgs {
				// get vuln info for that package NIST
				if risk, riskErr := utils.GetPkgRisk(pkg.Version, 0, limit); riskErr != nil {
					return riskErr
				} else {
					if risk.TotalResults > 0 {
						for idx, cpe := range risk.Result.CVEItems {
							if idx == 0 {
								tbl.AddRow(pkg.ID, pkg.Name, pkg.Version,
									cpe.Impact.BaseMetricV3.CvssV3.BaseScore,
									cpe.Impact.BaseMetricV3.CvssV3.BaseSeverity,
									cpe.Impact.BaseMetricV3.ImpactScore,
									cpe.Impact.BaseMetricV2.ExploitabilityScore)
							} else {
								tbl.AddRow(" - ", " - ", " - ",
									cpe.Impact.BaseMetricV3.CvssV3.BaseScore,
									cpe.Impact.BaseMetricV3.CvssV3.BaseSeverity,
									cpe.Impact.BaseMetricV3.ImpactScore,
									cpe.Impact.BaseMetricV2.ExploitabilityScore)
							}

						}
					}
				}
			}

			tbl.Print()

		}
		return nil
	},
}

func init() {
	pkgCmd.AddCommand(auditCmd)

	auditCmd.Flags().BoolP("findings", "f", false, "Only shows the packages with findings.")
	auditCmd.Flags().IntVarP(&limit, "limit", "l", 20, "Limits the number of findings per package. Default is 20 and -1 will get all (could take awhile).")
}
