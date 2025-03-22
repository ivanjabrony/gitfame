package stats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ivan-jabrony/gitfame/internal/fileparser"
)

type personForJSON struct {
	Name    string `json:"name"`
	Lines   int    `json:"lines"`
	Commits int    `json:"commits"`
	Files   int    `json:"files"`
}

func printTabular(stats *[]fileparser.PersonData) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprint(w, "Name\tLines\tCommits\tFiles\n")
	for _, person := range *stats {
		fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", person.Name, person.Lines, len(person.Commits), len(person.Files))
	}
	w.Flush()
}
func printCSV(stats *[]fileparser.PersonData) {
	w := csv.NewWriter(os.Stdout)
	err := w.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while writing csv: %s", err)
		os.Exit(1)
	}
	for _, person := range *stats {
		err := w.Write([]string{
			person.Name,
			strconv.Itoa(person.Lines),
			strconv.Itoa(len(person.Commits)),
			strconv.Itoa(len(person.Files)),
		})
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while writing csv: %s", err)
			os.Exit(1)
		}
	}
	w.Flush()
}

func printJSON(stats *[]fileparser.PersonData) {
	w := json.NewEncoder(os.Stdout)
	statsForJSON := make([]personForJSON, 0, len(*stats))
	for _, person := range *stats {
		outPerson := personForJSON{
			person.Name,
			person.Lines,
			len(person.Commits),
			len(person.Files),
		}
		statsForJSON = append(statsForJSON, outPerson)
	}
	err := w.Encode(statsForJSON)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while writing JSON: %s", err)
		os.Exit(1)
	}
}

func printJSONLines(stats *[]fileparser.PersonData) {
	w := json.NewEncoder(os.Stdout)
	for _, person := range *stats {
		outPerson := personForJSON{
			person.Name,
			person.Lines,
			len(person.Commits),
			len(person.Files),
		}
		err := w.Encode(&outPerson)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while writing JSON-lines: %s", err)
			os.Exit(1)
		}
	}
}

func PrintStatistics(stats *[]fileparser.PersonData, order, outputFormat string) {

	sort.Slice(*stats, func(i, j int) bool {
		switch order {
		case "lines":
			if (*stats)[i].Lines > (*stats)[j].Lines {
				return true
			} else if (*stats)[i].Lines < (*stats)[j].Lines {
				return false
			} else if len((*stats)[i].Commits) > len((*stats)[j].Commits) {
				return true
			} else if len((*stats)[i].Commits) < len((*stats)[j].Commits) {
				return false
			} else if len((*stats)[i].Files) > len((*stats)[j].Files) {
				return true
			} else if len((*stats)[i].Files) < len((*stats)[j].Files) {
				return false
			}
		case "commits":
			if len((*stats)[i].Commits) > len((*stats)[j].Commits) {
				return true
			} else if len((*stats)[i].Commits) < len((*stats)[j].Commits) {
				return false
			} else if (*stats)[i].Lines > (*stats)[j].Lines {
				return true
			} else if (*stats)[i].Lines < (*stats)[j].Lines {
				return false
			} else if len((*stats)[i].Files) > len((*stats)[j].Files) {
				return true
			} else if len((*stats)[i].Files) < len((*stats)[j].Files) {
				return false
			}
		case "files":
			if len((*stats)[i].Files) > len((*stats)[j].Files) {
				return true
			} else if len((*stats)[i].Files) < len((*stats)[j].Files) {
				return false
			} else if (*stats)[i].Lines > (*stats)[j].Lines {
				return true
			} else if (*stats)[i].Lines < (*stats)[j].Lines {
				return false
			} else if len((*stats)[i].Commits) > len((*stats)[j].Commits) {
				return true
			} else if len((*stats)[i].Commits) < len((*stats)[j].Commits) {
				return false
			}
		}

		return strings.ToLower((*stats)[i].Name) <= strings.ToLower((*stats)[j].Name)
	})

	switch outputFormat {
	case "csv":
		printCSV(stats)
	case "json":
		printJSON(stats)
	case "json-lines":
		printJSONLines(stats)
	default:
		printTabular(stats)
	}
}
