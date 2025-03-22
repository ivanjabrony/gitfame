package validation

import (
	"slices"
	"strconv"
	"testing"
)

var flagtests = []struct {
	name  string
	in    []string
	out   []string
	flags []string
}{
	{"extensions", []string{}, []string{}, []string{}},
	{"extensions", []string{"a.txt", "b.mp3", "c.dll"}, []string{"a.txt"}, []string{".txt"}},
	{"extensions", []string{"a.txt", "b.txt", "c.dll"}, []string{"a.txt", "b.txt"}, []string{".txt", ".bbb"}},

	{"languages", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"a1.cpp", "a2.cpp", "b.py"}, []string{"C++", "Python"}},
	{"languages", []string{"a.cpp", "b.py", "c.go"}, []string{"a.cpp", "b.py"}, []string{"c++", "python"}},
	{"languages", []string{"a.cpp", "b.py", "c.go"}, []string{"c.go"}, []string{"Go"}},

	{"exclude", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"b.py", "c.go"}, []string{"*.cpp"}},
	{"exclude", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"b.py"}, []string{"*.go", "*.cpp"}},
	{"exclude", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"a1.cpp", "a2.cpp", "c.go"}, []string{"*.py"}},

	{"restrict", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"a1.cpp", "a2.cpp"}, []string{"*.cpp"}},
	{"restrict", []string{"a1.cpp", "a2.cpp", "b.py", "c.go"}, []string{"b.py", "c.go"}, []string{"*.go", "*.py"}},
}

func TestExtensionFilter(t *testing.T) {
	for i, tt := range flagtests {
		name := strconv.Itoa(i) + " " + tt.name

		t.Run(name, func(t *testing.T) {
			result := []string{}
			var err error
			switch tt.name {
			case "extensions":
				result = ExtensionCheck(tt.in, tt.flags)
			case "languages":
				result, err = LanguagesCheck(tt.in, tt.flags)
				if err != nil {
					t.Errorf("Unexpected error occured: %v\n", err.Error())
				}
			case "exclude":
				result, err = ExcludeCheck(tt.in, tt.flags)
				if err != nil {
					t.Errorf("Unexpected error occured: %v\n", err.Error())
				}
			case "restrict":
				result, err = RestrictToCheck(tt.in, tt.flags)
				if err != nil {
					t.Errorf("Unexpected error occured: %v\n", err.Error())
				}
			}

			if !slices.Equal(result, tt.out) {
				t.Errorf("got %v, want %v", result, tt.out)
			}
		})
	}
}
