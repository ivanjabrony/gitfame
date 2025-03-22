package fileparser

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ParseLsTree(repositoryPath, revision string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", revision)
	cmd.Dir = repositoryPath

	lsTreeResult, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while parsing file-tree: %s", err)
		os.Exit(1)
	}
	res := strings.Split(string(lsTreeResult), "\n")

	return res[0 : len(res)-1], nil
}
