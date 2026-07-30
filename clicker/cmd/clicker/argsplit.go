package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vibium/clicker/internal/log"
)

// splitFlagsFromArgs separates the raw args of a command that disables
// cobra flag parsing into positional args and flag tokens. These commands
// take values that can start with a dash (sleep -5, geolocation -122.4,
// fill "#x" "-2"), which pflag would reject as unknown shorthand flags.
// A dash token that parses as a number stays positional; everything after
// "--" is positional. A literal non-numeric dash value (fill "#x" "--on")
// needs the "--" separator.
func splitFlagsFromArgs(cmd *cobra.Command, raw []string) (pos []string, flags []string) {
	afterDashDash := false
	for i := 0; i < len(raw); i++ {
		tok := raw[i]
		switch {
		case afterDashDash:
			pos = append(pos, tok)
		case tok == "--":
			afterDashDash = true
		case !strings.HasPrefix(tok, "-") || tok == "-" || isNumericToken(tok):
			pos = append(pos, tok)
		default:
			flags = append(flags, tok)
			if flagTakesSeparateValue(cmd, tok) && i+1 < len(raw) {
				i++
				flags = append(flags, raw[i])
			}
		}
	}
	return pos, flags
}

func isNumericToken(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// flagTakesSeparateValue reports whether tok is a non-boolean flag written
// without =value, meaning the next token is its value.
func flagTakesSeparateValue(cmd *cobra.Command, tok string) bool {
	if strings.Contains(tok, "=") {
		return false
	}
	var f *pflag.Flag
	if name, ok := strings.CutPrefix(tok, "--"); ok {
		f = lookupCommandFlag(cmd, func(fs *pflag.FlagSet) *pflag.Flag { return fs.Lookup(name) })
	} else if shorthand := strings.TrimPrefix(tok, "-"); len(shorthand) == 1 {
		f = lookupCommandFlag(cmd, func(fs *pflag.FlagSet) *pflag.Flag { return fs.ShorthandLookup(shorthand) })
	}
	return f != nil && f.Value.Type() != "bool"
}

func lookupCommandFlag(cmd *cobra.Command, lookup func(*pflag.FlagSet) *pflag.Flag) *pflag.Flag {
	if f := lookup(cmd.Flags()); f != nil {
		return f
	}
	return lookup(cmd.InheritedFlags())
}

// parseSplitFlags parses the separated flag tokens (including inherited
// globals like --json) and handles --help. Returns false when the command
// should stop because help was shown. Unknown flags exit like cobra does.
func parseSplitFlags(cmd *cobra.Command, flags []string) bool {
	// cmd.ParseFlags is a no-op when DisableFlagParsing is set. InheritedFlags
	// merges the persistent flags (--json, --verbose, ...) into cmd.Flags() so
	// they can be parsed directly.
	cmd.InheritedFlags()
	if err := cmd.Flags().Parse(flags); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprint(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}
	if help, _ := cmd.Flags().GetBool("help"); help {
		cmd.Help()
		return false
	}
	// The root's PersistentPreRun enabled logging before this parse ran, so
	// apply --verbose here or it would be read too late to take effect.
	if verbose {
		log.Setup(log.LevelVerbose)
	}
	return true
}

// requirePositionalArgs enforces the positional count with cobra's error
// wording, since Args validators see the raw args when flag parsing is
// disabled.
func requirePositionalArgs(cmd *cobra.Command, pos []string, min, max int) {
	if len(pos) >= min && len(pos) <= max {
		return
	}
	if min == max {
		fmt.Fprintf(os.Stderr, "Error: accepts %d arg(s), received %d\n", min, len(pos))
	} else {
		fmt.Fprintf(os.Stderr, "Error: accepts between %d and %d arg(s), received %d\n", min, max, len(pos))
	}
	fmt.Fprint(os.Stderr, cmd.UsageString())
	os.Exit(1)
}
