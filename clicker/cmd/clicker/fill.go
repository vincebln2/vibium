package main

import (
	"time"

	"github.com/spf13/cobra"
)

func newFillCmd() *cobra.Command {
	var timeout time.Duration
	cmd := &cobra.Command{
		Use:   "fill [selector] [text]",
		Short: "Clear an input field and type new text",
		Example: `  vibium fill "input[name=email]" "user@example.com"
  # Clear the field and type new value
  vibium fill "#search" "vibium"
  # Replace search field contents

  vibium fill "#search" "vibium" --timeout 5s
  # Custom timeout (5s, or 5000 for milliseconds)`,
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			pos, flags := splitFlagsFromArgs(cmd, args)
			if !parseSplitFlags(cmd, flags) {
				return
			}
			requirePositionalArgs(cmd, pos, 2, 2)
			selector := pos[0]
			text := pos[1]
			result, err := daemonCall("browser_fill", map[string]interface{}{
				"selector": selector,
				"value":    text,
				"timeout":  float64(timeout.Milliseconds()),
			})
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	addTimeoutFlag(cmd, &timeout)
	return cmd
}
