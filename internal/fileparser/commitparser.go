package fileparser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func parseGitLog(filepath, repositoryPath, revision string) Commit {
	gitLogCommand := exec.Command("git", "log", "--name-only", "--pretty=full", revision, "--", filepath)
	gitLogCommand.Dir = repositoryPath

	gitLogOutput, err := gitLogCommand.Output()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while parsing file(%v): %s", filepath, err)
		os.Exit(1)
	}

	splitOutput := strings.Split(string(gitLogOutput), "\n")
	if len(splitOutput) > 1 {

		authorLine := strings.Split(splitOutput[1], " ")
		committerLine := strings.Split(splitOutput[2], " ")

		lastCommitSha := strings.Split(splitOutput[0], " ")[1]
		lastCommitAuthor := strings.Join(authorLine[1:len(authorLine)-1], " ")
		lastCommitCommitter := strings.Join(committerLine[1:len(committerLine)-1], " ")

		return Commit{
			SHA:      lastCommitSha,
			Author:   lastCommitAuthor,
			Commiter: lastCommitCommitter,
			Lines:    0,
			File:     filepath}
	}
	return Commit{}
}

func parseCommitsInFile(commitData, file string) []Commit {
	commitBySHA := make(map[string]*Commit)
	commitStrings := make([]string, 0, len(commitBySHA))
	SHAByAuthor := make(map[string]string)
	splittedBlame := strings.Split(commitData, "\n")
	for _, line := range splittedBlame {
		if !strings.HasPrefix(line, "	") {
			commitStrings = append(commitStrings, line)
		}
	}

	for i, line := range commitStrings {
		splittedLine := strings.Split(line, " ")

		//new commit with new author and commiter
		if splittedLine[0] == "author" {
			curentSHA := strings.Split(commitStrings[i-1], " ")[0]
			currentCommiter := strings.Join(strings.Split(commitStrings[i+4], " ")[1:], " ")
			currentAuthor := strings.Join(strings.Split(commitStrings[i], " ")[1:], " ")
			SHAByAuthor[curentSHA] = currentAuthor

			commitBySHA[curentSHA] =
				&Commit{
					SHA:      curentSHA,
					Author:   currentAuthor,
					Commiter: currentCommiter,
					File:     file,
					Lines:    1}
		} else if curCommit, ok := commitBySHA[splittedLine[0]]; ok {
			curCommit.Lines++
		}

	}
	commits := make([]Commit, 0, len(commitBySHA))
	for _, v := range commitBySHA {
		commits = append(commits, *v)
	}
	return commits
}

func ParseFile(filepath, repositoryPath, revision string) []Commit {
	commits := []Commit{}
	cmd := exec.Command("git", "blame", "--porcelain", revision, "--", filepath)
	cmd.Dir = repositoryPath

	outputRaw, err := cmd.Output()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Internal error occured while parsing file(%v): %s", filepath, err)
		os.Exit(1)
	}
	if len(outputRaw) == 0 {
		commits = append(commits, parseGitLog(filepath, repositoryPath, revision))
	} else {
		commits = parseCommitsInFile(string(outputRaw), filepath)
	}

	log.Printf("%v - finished parsing file.\n", filepath)
	return commits
}
