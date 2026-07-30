package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func newGeolocationCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "geolocation [latitude] [longitude]",
		Short: "Override the browser geolocation",
		Example: `  vibium geolocation 40.7128 -74.006
  # Set location to New York City

  vibium geolocation 51.5074 -0.1278 --accuracy 10
  # Set location to London with 10m accuracy`,
		DisableFlagParsing: true,
		Args:               cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			pos, flags := splitFlagsFromArgs(cmd, args)
			if !parseSplitFlags(cmd, flags) {
				return
			}
			requirePositionalArgs(cmd, pos, 2, 2)
			lat, err := strconv.ParseFloat(pos[0], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid latitude: %s\n", pos[0])
				os.Exit(1)
			}
			lng, err := strconv.ParseFloat(pos[1], 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid longitude: %s\n", pos[1])
				os.Exit(1)
			}

			accuracy, _ := cmd.Flags().GetFloat64("accuracy")

			callArgs := map[string]interface{}{
				"latitude":  lat,
				"longitude": lng,
			}
			if accuracy > 0 {
				callArgs["accuracy"] = accuracy
			}
			result, err := daemonCall("browser_set_geolocation", callArgs)
			if err != nil {
				printError(err)
				return
			}
			printResult(result)
		},
	}
	cmd.Flags().Float64("accuracy", 0, "Accuracy in meters (default: 1)")
	return cmd
}
