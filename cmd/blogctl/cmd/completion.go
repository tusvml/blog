package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var shell string

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Output shell completion (bash/fish/powershell/zsh)",
	Long:  "Output shell completion (bash/fish/powershell/zsh).",
	RunE:  runCompletion,
}

func runCompletion(cmd *cobra.Command, args []string) error {
	switch shell {
	case "bash":
		return fmt.Errorf("failed to create completion: %w", rootCmd.GenBashCompletion(os.Stdout))
	case "fish":
		return fmt.Errorf("failed to create completion: %w", rootCmd.GenFishCompletion(os.Stdout, true))
	case "powershell":
		return fmt.Errorf("failed to create completion: %w", rootCmd.GenPowerShellCompletion(os.Stdout))
	case "zsh":
		return fmt.Errorf("failed to create completion: %w", rootCmd.GenZshCompletion(os.Stdout))
	}

	return fmt.Errorf("invalid shell: %s", shell)
}

func init() {
	completionCmd.Flags().StringVarP(&shell, "shell", "s", "bash", "shell type (bash/fish/powershell/zsh)")

	rootCmd.AddCommand(completionCmd)
}
