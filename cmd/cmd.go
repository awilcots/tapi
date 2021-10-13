package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/awilcots/tapi/store"
	"github.com/spf13/cobra"
)

var (
	fromPipe bool
	fromFile bool
	fromArgs bool
	duration string

	show bool

	rootCmd = &cobra.Command{
		Use:   "tapi",
		Short: "Tapi is a simple api built to test testing",
		Long:  "Test API or tapi is a simple program build so that we can have something to make some tests for",
		Run:   rootRun,
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&fromPipe, "pipe", "p", false, "if the data to store is being piped here")
	rootCmd.Flags().BoolVarP(&fromFile, "file", "f", false, "file to store")
	rootCmd.Flags().BoolVarP(&fromArgs, "data", "d", false, "data to store")
	rootCmd.Flags().StringVarP(&duration, "ttl", "t", "10m", "amount of time you'd like the data to be stored")

	rootCmd.Flags().BoolVarP(&show, "show", "s", false, "view stored data")
}

func rootRun(cmd *cobra.Command, args []string) {
	kvs := store.NewKVSStore()
	if show {
		kvs.Show()
		return
	}

	if fromPipe && fromFile && fromArgs ||
		fromPipe && fromFile ||
		fromPipe && fromArgs ||
		fromFile && fromArgs {
		fmt.Fprintln(os.Stderr, "[!] Must only supply one of pipe, file, or data flags")
		cmd.Usage()
		os.Exit(1)
	}
	if fromFile && len(args) < 1 {
		fmt.Fprintln(os.Stderr, "[!] Must supply file name if file flag is used")
		cmd.Usage()
		os.Exit(1)
	}
	if fromArgs && len(args) < 1 {
		fmt.Fprintln(os.Stderr, "[!] No data to store found")
		cmd.Usage()
		os.Exit(1)
	}
	if !fromPipe && !fromFile && !fromArgs {
		cmd.Usage()
		os.Exit(0)
	}
	if !regexp.MustCompile(`^\d+[smdh]$`).MatchString(duration) {
		fmt.Fprintln(os.Stderr, "[!!] ttl must be in the format number(s) letter. e.g. 10s, 20d, 4m")
	}

	kvs.Strategy = store.StoreStrategy{
		Pipe: fromPipe,
		File: fromFile,
		Args: fromArgs,
	}

	kvs.Save(duration, args)
	fmt.Println("[*] Data Stored!")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
