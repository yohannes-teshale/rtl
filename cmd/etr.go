package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"rtl/etr"
)

var etrCommand = &cobra.Command{
	Use:   "etr",
	Short: "Enzyme to React testing library migration Tool.",
	Long:  `Long( Enzyme to React testing library migration Tool. )`,
	Run:   entry,
}

func init() {
	rootCmd.AddCommand(etrCommand)
	etrCommand.Flags().StringP("testDir", "t", "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing/src", "Target Directory containing test files.")
	etrCommand.Flags().StringP("projDir", "p", "/Users/yohannes/playground/hackathon/react-typescript-jest-enzyme-testing", "Target project Directory.")
	etrCommand.Flags().StringP("outDir", "o", "/generatedRTLfiles", "Directory for generated files.")
}

func entry(cmd *cobra.Command, args []string) {
	testDir, _ := cmd.Flags().GetString("testDir")
	projDir, _ := cmd.Flags().GetString("projDir")
	outDir, _ := cmd.Flags().GetString("outDir")
	outDir = projDir + outDir
	fmt.Println(testDir, "\n", projDir, "\n", outDir)

	etr.Main(testDir, projDir, outDir)

}
