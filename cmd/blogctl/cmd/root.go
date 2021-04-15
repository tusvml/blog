package cmd

import (
	"github.com/spf13/cobra"
)

var version bool

var rootCmd = &cobra.Command{
	Use:   "notion2hugo",
	Short: "notion2hugo is a CLI for converting Notion pages to Hugo blogs",
	Long:  "notion2hugo is a CLI for converting Notion pages to Hugo blogs.",
	RunE:  runRoot,
}

func runRoot(cmd *cobra.Command, args []string) error {
	if version {
		return runVersion(cmd, args)
	}

	return nil
}

func init() {
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "show version")
}

func Execute() {
	rootCmd.Execute() //nolint:errcheck
}
