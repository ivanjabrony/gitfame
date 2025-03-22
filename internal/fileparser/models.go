package fileparser

type PersonData struct {
	Name    string
	Lines   int
	Commits map[string]struct{}
	Files   map[string]struct{}
}

type Commit struct {
	SHA      string
	Author   string
	Commiter string
	File     string
	Lines    int
}
