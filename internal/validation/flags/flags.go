package flags

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Validate(
	repository *string,
	revision *string,
	order *string,
	useCommitter *bool,
	format *string,
	extensions *[]string,
	languages *[]string,
	excludePatterns *[]string,
	restrictPatterns *[]string) {

	cmd := exec.Command("git", "status", "2>/dev/null;", "echo", "$?")
	cmd.Dir = *repository
	output, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Flags validation failed: %s", err)
		os.Exit(1)
	} else if string(output) == "0" {
		_, _ = fmt.Fprintf(os.Stderr, "Flags validation failed: %v is not a git repository", *repository)
		os.Exit(1)
	}

	revCheck := exec.Command("git", "rev-parse", "--verify", "-q", *revision)
	revCheck.Dir = *repository
	revCheckOutput, err := revCheck.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Flags validation failed: %v is not a valid revision", revCheckOutput)
		os.Exit(1)
	}

	if *order != "lines" && *order != "commits" && *order != "files" {
		_, _ = fmt.Fprintf(os.Stderr, "Flags validation failed: %v is not a valid order format", order)
		os.Exit(1)
	}

	if *format != "tabular" && *format != "csv" && *format != "json" && *format != "json-lines" {
		_, _ = fmt.Fprintf(os.Stderr, "Flags validation failed: %v is not a valid output format", format)
		os.Exit(1)
	}

	if *extensions != nil {
		if strings.Contains((*extensions)[0], ",") {
			*extensions = strings.Split((*extensions)[0], ",")
		}
	}

	if *languages != nil {
		if strings.Contains((*languages)[0], ",") {
			*languages = strings.Split((*languages)[0], ",")
		}
	}

	if *excludePatterns != nil {
		if strings.Contains((*excludePatterns)[0], ",") {
			*excludePatterns = strings.Split((*excludePatterns)[0], ",")
		}
	}

	if *restrictPatterns != nil {
		if strings.Contains((*restrictPatterns)[0], ",") {
			*restrictPatterns = strings.Split((*restrictPatterns)[0], ",")
		}
	}
}
