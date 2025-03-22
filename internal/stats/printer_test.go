package stats

import (
	"strconv"
	"testing"

	"github.com/ivan-jabrony/gitfame/internal/fileparser"
)

var flagtests = []struct {
	name         string
	in           *[]fileparser.PersonData
	order        string
	outputFormat string
}{
	{"linear", &[]fileparser.PersonData{
		{
			Name:    "Ivan Zabrodin",
			Lines:   7,
			Files:   map[string]struct{}{"a.txt": {}, "b.txt": {}, "c.txt": {}},
			Commits: map[string]struct{}{"aaa": {}, "bbb": {}},
		},
		{
			Name:    "Arseniy Borozdov",
			Lines:   6,
			Files:   map[string]struct{}{"a.txt": {}, "b.txt": {}, "c.txt": {}, "d.txt": {}},
			Commits: map[string]struct{}{"ccc": {}, "ddd": {}}}},
		"Files",
		"tabular"},
}

func TestStatsPrinter(t *testing.T) {
	for i, tt := range flagtests {
		name := strconv.Itoa(i) + " " + tt.name

		t.Run(name, func(t *testing.T) {
			PrintStatistics(tt.in, tt.order, tt.outputFormat)
		})
	}
}
