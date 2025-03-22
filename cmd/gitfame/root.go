package main

import (
	"fmt"
	"os"

	"github.com/ivan-jabrony/gitfame/internal/fileparser"
	"github.com/ivan-jabrony/gitfame/internal/stats"
	"github.com/ivan-jabrony/gitfame/internal/validation"
	"github.com/ivan-jabrony/gitfame/internal/validation/flags"
	"github.com/spf13/cobra"
)

var (
	repository       string
	revision         string
	order            string
	useCommitter     bool
	format           string
	maxThreads       int
	extensions       []string
	languages        []string
	excludePatterns  []string
	restrictPatterns []string
)
var rootCmd = &cobra.Command{
	Use:   "gitfame",
	Short: "a CLI that can print statistics about git repository authors and their participation in a project",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		flags.Validate(&repository, &revision, &order, &useCommitter, &format, &extensions, &languages, &excludePatterns, &restrictPatterns)

		filePaths, err := fileparser.ParseLsTree(repository, revision)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while parsing ls-tree: %s", err)
			os.Exit(1)
		}

		filePaths = validation.ExtensionCheck(filePaths, extensions)

		if filePaths, err = validation.LanguagesCheck(filePaths, languages); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if filePaths, err = validation.ExcludeCheck(filePaths, excludePatterns); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if filePaths, err = validation.RestrictToCheck(filePaths, restrictPatterns); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		personData := stats.RunParallel(filePaths, maxThreads, revision, repository, useCommitter)
		stats.PrintStatistics(personData, order, format)
	},
}

func init() {
	curPath, _ := os.Executable()
	rootCmd.Flags().StringVar(&repository, "repository", curPath, "Path to git repository. Current path of executable on default")
	rootCmd.Flags().StringVar(&revision, "revision", "HEAD", "Pointer to the commit. 'HEAD' on default")
	rootCmd.Flags().StringVar(&order, "order-by", "lines", "Key order of printed result. May be 'lines'(default), 'commits' or 'files'")
	rootCmd.Flags().BoolVar(&useCommitter, "use-committer", false, "Changes counting policy from author statistics to committer statistics")
	rootCmd.Flags().StringVar(&format, "format", "tabular", "Defines output format. Either 'tabular'(default), 'csv', 'json' or 'json-file'")
	rootCmd.Flags().IntVar(&maxThreads, "maxthr", 8, "Defines maximum amount of threads at the same time. 8 by default'")
	rootCmd.Flags().StringArrayVar(&extensions, "extensions", nil, "Defines a list of file extensions that should be counted")
	rootCmd.Flags().StringArrayVar(&languages, "languages", nil, "Defines a list of languages that should be counted")
	rootCmd.Flags().StringArrayVar(&excludePatterns, "exclude", nil, "List of Glob patterns, all files that fall under one of the patterns will not be counted")
	rootCmd.Flags().StringArrayVar(&restrictPatterns, "restrict-to", nil, "List of Glob patterns, any file that don't fall under atleast one of the patterns won't be counted")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured on flag parsing stage: %s", err)
		os.Exit(1)
	}
}
