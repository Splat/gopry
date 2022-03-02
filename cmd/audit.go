/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopry/utils"
	"regexp"
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
			tbl := table.New("ID", "Name", "Version", "Risk", "Base", "Impact", "Exploitability", "CVE", "Description", "Published")
			// format some things
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			columnFmt := color.New(color.FgYellow).SprintfFunc()
			tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithWidthFunc(printStrWidth)
			// build the rows
			//tbl.AddRow(greened("ID"), greened("Name"), greened("Version"), greened("Risk"), greened("Base"), greened("Impact"), greened("Ease"), greened("CVE"), greened("Desc"), greened("Published"))
			for _, pkg := range sysPkgs {
				// get vuln info for that package NIST
				if risk, riskErr := utils.GetPkgRisk(pkg.Version, 0, limit); riskErr != nil {
					return riskErr
				} else {
					if risk.TotalResults > 0 {
						for _, cpe := range risk.Result.CVEItems {
							//if idx == 0 {
							tbl.AddRow(pkg.ID, pkg.Name, pkg.Version,
								getRowColor(cpe.Impact.BaseMetricV3.CvssV3.BaseSeverity)(
									cpe.Impact.BaseMetricV3.CvssV3.BaseSeverity),
								fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV3.CvssV3.BaseScore),
								fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV3.ImpactScore),
								fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV2.ExploitabilityScore),
								cpe.Cve.CVEDataMeta.ID,
								cpe.Cve.Description.DescriptionData[0].Value[:25]+`...`,
								cpe.PublishedDate).WithWidthFunc(printStrWidth)
							//} else {
							//	tbl.AddRow(" - ", " - ", " - ",
							//		colorF(cpe.Impact.BaseMetricV3.CvssV3.BaseSeverity),
							//		fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV3.CvssV3.BaseScore),
							//		fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV3.ImpactScore),
							//		colorF(fmt.Sprintf("%.1f", cpe.Impact.BaseMetricV2.ExploitabilityScore)),
							//		colorF(cpe.Cve.CVEDataMeta.ID),
							//		cpe.Cve.Description.DescriptionData[0].Value[:25]+`...`,
							//		cpe.PublishedDate).WithWidthFunc(runewidth.StringWidth)
							//}

						}
					}
				}
			}

			tbl.Print()
		}
		return nil
	},
}

// StringWidth return width (int) as you can see in the terminal.
// Strips ANSI escape codes to calculate the length of a string as it would appear in the terminal used
// in conjunction with the rodaine CLI table library.
// Parameters:
//		str (string) - the string who's width needs to be calculated as it would render on the screen
func printStrWidth(str string) (width int) {
	const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

	var ansiRegxp = regexp.MustCompile(ansi)

	cleanStr := ansiRegxp.ReplaceAllString(str, "")
	return len(cleanStr)
}

// getRowColor (private) - Given a severity returns the proper colorize function to be
// used in the table output. Default output is non-colorized text.
// Parameters:
//		cpeSeverity (str) - the severity { "LOW", "MEDIUM", "HIGH", "CRITICAL" }. Should be an enum-like struct refactor.
//
// TODO: This all should eventually be moved to a separate table printer function/module as the output of audit will eventually include other options like CSV and JSON
func getRowColor(cpeSeverity string) func(string) string {
	switch cpeSeverity {
	case "LOW":
		return func(value string) string { return color.BlueString(value) }
	case "MEDIUM":
		return func(value string) string { return color.HiYellowString(value) }
	case "HIGH":
		return func(value string) string { return color.YellowString(value) }
	case "CRITICAL":
		return func(value string) string { return color.RedString(value) }
	default:
		return func(value string) string {
			if value == "" {
				return color.GreenString("UNKNOWN")
			}
			return color.GreenString(value)
		}
	}
}

func greened(str string) string {
	return color.New(color.FgGreen, color.Underline).Sprintf(str)
}

func init() {
	pkgCmd.AddCommand(auditCmd)

	auditCmd.Flags().BoolP("findings", "f", false, "Only shows the packages with findings.")
	auditCmd.Flags().IntVarP(&limit, "limit", "l", 20, "Limits the number of findings per package. Default is 20 and -1 will get all (could take awhile).")
}
