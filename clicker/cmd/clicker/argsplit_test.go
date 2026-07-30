package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func newSplitTestCmd() *cobra.Command {
	root := &cobra.Command{Use: "root"}
	root.PersistentFlags().Bool("json", false, "")
	root.PersistentFlags().BoolP("verbose", "v", false, "")
	cmd := &cobra.Command{Use: "sub", Run: func(*cobra.Command, []string) {}}
	cmd.Flags().String("timeout", "", "")
	cmd.Flags().Float64("accuracy", 0, "")
	root.AddCommand(cmd)
	return cmd
}

func TestSplitFlagsFromArgs(t *testing.T) {
	cmd := newSplitTestCmd()
	cases := []struct {
		in        []string
		wantPos   []string
		wantFlags []string
	}{
		{[]string{"200", "--json"}, []string{"200"}, []string{"--json"}},
		{[]string{"--json", "200"}, []string{"200"}, []string{"--json"}},
		{[]string{"-5"}, []string{"-5"}, nil},
		{[]string{"37.7", "-122.4"}, []string{"37.7", "-122.4"}, nil},
		{[]string{"#x", "-2", "--json"}, []string{"#x", "-2"}, []string{"--json"}},
		{[]string{"#x", "v", "--timeout", "5s"}, []string{"#x", "v"}, []string{"--timeout", "5s"}},
		{[]string{"#x", "--timeout=5s", "v"}, []string{"#x", "v"}, []string{"--timeout=5s"}},
		{[]string{"37.7", "-122.4", "--accuracy", "10"}, []string{"37.7", "-122.4"}, []string{"--accuracy", "10"}},
		{[]string{"-v", "x"}, []string{"x"}, []string{"-v"}},
		{[]string{"--", "-v"}, []string{"-v"}, nil},
		{[]string{"#x", "--", "--not-a-flag"}, []string{"#x", "--not-a-flag"}, nil},
		{[]string{"--nope", "x"}, []string{"x"}, []string{"--nope"}},
		{[]string{"-"}, []string{"-"}, nil},
		{nil, nil, nil},
	}
	for _, c := range cases {
		pos, flags := splitFlagsFromArgs(cmd, c.in)
		if !reflect.DeepEqual(pos, c.wantPos) || !reflect.DeepEqual(flags, c.wantFlags) {
			t.Errorf("splitFlagsFromArgs(%q): got pos=%q flags=%q, want pos=%q flags=%q",
				c.in, pos, flags, c.wantPos, c.wantFlags)
		}
	}
}
